package engine

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "net/http"
    "strings"
    "time"

    "github.com/pkg/errors"
    "gopkg.in/tomb.v2"

    "baetyl-simulator/common"
    "baetyl-simulator/config"
    "baetyl-simulator/constants"
    "baetyl-simulator/core/utils"
    "baetyl-simulator/middleware/log"
    "baetyl-simulator/service"
    specV1 "baetyl-simulator/spec/v1"
)

type Engine interface {
    Start(ctx context.Context)
    Close()
}

type engine struct {
    nodeService service.NodeService
    syncService service.SyncService
    node        *specV1.Node
    tomb        tomb.Tomb
    cfg         *config.Config
    log         *log.Logger
    tmpl        service.TemplateService
}

func NewEngine(ctx context.Context, name string, config *config.Config) (Engine, error) {
    c := context.WithValue(ctx, "nodeName", name)
    syncService, err := service.NewSyncService(c, &config.Cloud.Sync)
    if err != nil {
        return nil, err
    }
    nodeService, err := service.NewNodeService(&config.Cloud.Admin)
    if err != nil {
        return nil, err
    }

    node, err := nodeService.Get("", name)
    if err != nil {
        return nil, err
    }
    log.L().Debug("lookup node", log.Any("node", node))

    if err = common.Cache.Set(name, node, -1); err != nil {
        return nil, err
    }

    templatePath := config.Template.Path
    if templatePath == "" {
        templatePath = constants.DEFAULT_TEMPLATE_DIR
    }

    t, err := service.NewTemplateService(templatePath, map[string]interface{}{})
    if err != nil {
        return nil, err
    }

    return &engine{
        node: node,
        nodeService: nodeService,
        syncService: syncService,
        cfg: config,
        log:          log.With(log.Any("core", "engine"), log.Any("node", name)),
        tmpl: t,
    }, nil
}

func (e *engine) Start(ctx context.Context) {
    e.tomb.Go(e.reporting)
    e.tomb.Go(e.receiving)

    for {
       select {
       case <- ctx.Done():
           e.log.Info("node stop", log.Any("error", ctx.Err()))
           e.Close()
           return
       default:
           e.log.Debug("engine running")
           time.Sleep(30 * time.Second)
       }
    }
}

func (e *engine) Close() {
    e.tomb.Kill(nil)
    if err := e.tomb.Wait(); err != nil {
        log.L().Error("engine close error", log.Any("error", err))
    }
}

func (e *engine) reporting() error {
    defer utils.HandlePanic()

    e.log.Info("start reporting...")

    params := e.getParams()
    params["NodeName"] = e.node.Name
    params["ReportTime"] = time.Now().Local()
    reportData, err := e.tmpl.ParseTemplate(constants.REPORT_TEMPLATE, params)
    if err != nil {
        return err
    }

    t := time.NewTicker(e.cfg.Engine.Report.Interval)
    defer t.Stop()
    for {
        select {
        case <-t.C:
            res, _ := e.reportAndDesire(specV1.MessageReport, reportData)
            if err != nil {
                e.log.Error("report error", log.Any("error", errors.WithStack(err)))
            } else {
                var desire specV1.Desire
                err = res.Content.Unmarshal(&desire)
                if err != nil {
                    e.log.Error("report error", log.Any("error", errors.WithStack(err)))
                } else {
                    if err = e.MergeDesire(desire); err != nil {
                        e.log.Error("report error", log.Any("error", errors.WithStack(err)))
                    }
                }
            }
        case <-e.tomb.Dying():
            e.log.Info("stop syncApps")
            return nil
        }
    }
}

func (e *engine) receiving() error {
    defer utils.HandlePanic()
    e.log.Info("start receiving...")

    t := time.NewTicker(e.cfg.Engine.Desire.Interval)
    defer t.Stop()

    for {
        select {
        case <-t.C:
            delta, err := e.GetShadowDelta()
            if err != nil {
                e.log.Error("get node shadow delta error", log.Any("error", err))
            }

            if delta != nil {
                if err = e.syncApps(true, delta); err != nil {
                    e.log.Error("sync sys apps error", log.Any("error", err))
                }

                if err = e.syncApps(false, delta); err != nil {
                    e.log.Error("sync apps error", log.Any("error", err))
                }
            }
        case <-e.tomb.Dying():
            log.L().Info("stop syncApps")
            return nil
        }
    }
}

func (e *engine) syncApps(isSys bool, delta specV1.Delta) error {
    defer utils.HandlePanic()
    //TODO remove get  desire data from template
    //params := e.getParams()
    //params["NodeName"] = e.node.Name
    //dd, err := e.tmpl.ParseTemplate(constants.DESIRE_TEMPLATE, params)
    //if err != nil {
    //    return err
    //}

    dapps := specV1.Desire(delta).AppInfos(isSys)
    if dapps == nil {
       return nil
    }

    appInfo := make(map[string]string)
    for _, info := range dapps {
       appInfo[info.Name] = info.Version
    }
    dq := specV1.DesireRequest{Infos: e.genResourceInfos(specV1.KindApplication, appInfo)}
    e.log.Info("start desire", log.Any("delta", dq))
    desire, err := json.Marshal(dq)
    if err != nil {
       return err
    }

    res, err := e.reportAndDesire(specV1.MessageDesire, desire)
    if err != nil {
        return errors.Wrap(err, "call reportAndDesire error")
    } else {
        var report specV1.Report
        err = res.Content.Unmarshal(&report)
        if err != nil {
            return errors.Wrap(err, "unmarshal syncApps error")
        } else {
            if err = e.MergeReport(report); err != nil {
                return errors.Wrap(err, "merge report error")
            }
        }
    }

    return nil
}

func (e *engine) reportAndDesire(kind specV1.MessageKind, body []byte) (*specV1.Message, error) {
    var resp *http.Response
    var err error
    res := &specV1.Message{Kind: kind}
    start := time.Now()
    switch kind {
    case specV1.MessageReport:
        resp, err = e.syncService.Report(body)
    case specV1.MessageDesire:
        resp, err = e.syncService.Desire(body)
    default:
        return nil, errors.Errorf("unsupported message kind: %s", kind)
    }
    e.log.Debug("request finish", log.Any("kind", kind), log.Any("start", start.Local()),
        log.Any("cost", time.Since(start).Seconds()))
    //fmt.Printf("report, start: %s, cost: %s \n", start.Local(), time.Since(start))
    if err != nil {
        log.L().Error("request fail", log.Any("error", err))
    } else {
        if resp.StatusCode == 200 {
            log.L().Info("response", log.Any("code", resp.StatusCode),
                log.Any("status", resp.Status))
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                log.L().Error("read response fail", log.Any("error", err))
            } else {

                data, err := utils.ParseEnv(body)
                if err != nil {
                    log.L().Error("parse env fail", log.Any("error", errors.WithStack(err)))
                }
                res.Content.SetJSON(data)
                log.L().Info("receive response", log.Any("info", res))
            }
        } else {
            log.L().Error("request fail", log.Any("code", resp.StatusCode),
                log.Any("status", resp.Status))
        }
    }
    Close(resp.Body)

    return res, nil
}

func (e *engine) getParams() map[string]interface{} {
    params := make(map[string]interface{}, 3)
    sapps := e.node.Desire.AppInfos(true)
    for _, app := range sapps {
        name, val := app.Name, app.Version
        if strings.Contains(name, constants.BaetylInit) {
            params[constants.TmplVarBaetylInit] = name
            params[constants.TmplVarBaetylInitVer] = val
        } else if strings.Contains(name, constants.BaetylCore) {
            params[constants.TmplVarBaetylCore] = name
            params[constants.TmplVarBaetylCoreVer] = val
        } else if strings.Contains(name, constants.BaetylBroker) {
            params[constants.TmplVarBaetylBroker] = name
            params[constants.TmplVarBaetylBrokerVer] = val
        }
    }
    return params
}

func Close(body io.ReadCloser) {
    err := body.Close()
    if err != nil {
        fmt.Println("client reader close fail", err)
    }
}
