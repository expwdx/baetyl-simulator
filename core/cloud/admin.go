package cloud

import (
	"baetyl-simulator/constants"
	"baetyl-simulator/core/utils"
	"baetyl-simulator/errors"
	"baetyl-simulator/model"
	"context"
	"encoding/json"
	"gopkg.in/tomb.v2"
	"strings"
	"time"

	"baetyl-simulator/config"
	"baetyl-simulator/middleware/log"
	"baetyl-simulator/service"
)

type Admin interface {
	Start(ctx context.Context)
}

type adminImpl struct {
	app      	service.ApplicationService
	node 		service.NodeService
	tomb        tomb.Tomb
	cfg         *config.Config
	log         *log.Logger
	tmpl		service.TemplateService
}

func NewAdmin (config *config.Config) (Admin, error) {
	as, err := service.NewApplicationService(&config.Cloud.Admin)
	if err != nil {
		return nil, err
	}

	nodeService, err := service.NewNodeService(&config.Cloud.Admin)
	if err != nil {
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

	return &adminImpl{
		app: as,
		node: nodeService,
		cfg: config,
		log: log.With(log.Any("core", "cloud")),
		tmpl: t,
	}, nil
}

func (a *adminImpl) Start(ctx context.Context) {
	a.tomb.Go(a.deploying)
	a.tomb.Go(a.reading)

	for {
		select {
		case <- ctx.Done():
			a.log.Info("node stop", log.Any("error", ctx.Err()))
			a.Close()
			return
		default:
			a.log.Debug("engine running")
			time.Sleep(30 * time.Second)
		}
	}
}

func (a *adminImpl) Close() {
	a.tomb.Kill(nil)
	if err := a.tomb.Wait(); err != nil {
		a.log.Error("engine close error", log.Any("error", err))
	}
}

func (a *adminImpl) deploying() error {
	defer utils.HandlePanic()
	var appView *model.ApplicationView
	primaryAppView, err := a.loadApp(constants.APP_MYSQL_TEMPLATE)
	if err != nil {
		return err
	}
	modAppView := primaryAppView
	modAppView.Services[0].Image = "harbor.sz.yingzi.com/base/mysql:6"

	_, err = a.createApp(primaryAppView)
	if err != nil {
		return err
	}

	flag := 0

	t := time.NewTicker(a.cfg.User.Deploy.Interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			flag = ^flag
			switch flag {
			case 0:
				appView = primaryAppView
			case -1:
				appView = modAppView
			default:
				a.log.Error("unsupported flag", log.Any("flag", flag))
				return errors.Errorf("unsupported flag %d", flag)
			}
			_, err := a.updateApp(appView.Name, appView)
			if err != nil {
				log.L().Error("update app fail", log.Any("error", err))
			}
		case <-a.tomb.Dying():
			log.L().Info("stop deploying")
			return nil
		}
	}
}

func (a *adminImpl) reading() error {
	defer utils.HandlePanic()
	t := time.NewTicker(a.cfg.User.Read.Interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			nodes, err := a.node.List("")
			if err != nil {
				log.L().Error("get node list fail", log.Any("error", errors.Trace(err)))
			}

			log.L().Info("get node list", log.Any("info", nodes))
		case <-a.tomb.Dying():
			log.L().Info("stop reading")
			return nil
		}
	}
}

func (a *adminImpl) getParams() (map[string]interface{}, error) {
	params := make(map[string]interface{}, 0)
	labels := a.cfg.Mock.NodeLabels
	if len(labels) == 0 {
		return nil, errors.New("no label included in config")
	}

	var mutilLabelsStr string

	for k, v := range labels {
		labelStr := strings.Join([]string{k, v}, "=")
		mutilLabelsStr = strings.Join([]string{mutilLabelsStr, labelStr}, ",")
	}
	params["NodeSelector"] = strings.TrimLeft(mutilLabelsStr, ",")
	params["AppName"] = a.cfg.Mock.AppName

	return params, nil
}

func (a *adminImpl) loadApp(tmpl string) (*model.ApplicationView, error) {
	params, err := a.getParams()
	if err != nil {
		return nil, err
	}

	tmplApp, err := a.tmpl.ParseTemplate(tmpl, params)
	if err != nil {
		return nil, err
	}

	var app *model.ApplicationView
	if err := json.Unmarshal(tmplApp, &app); err != nil {
		return nil, err
	}

	return app, nil

}

func (a *adminImpl) createApp (appView *model.ApplicationView) (*model.ApplicationView, error) {
	log.L().Info("create app request", log.Any("info", appView))

	res, err := a.app.Create(nil, appView)
	if err != nil {
		return nil, err
	}

	log.L().Info("create app response", log.Any("info", res))

	return res, nil
}

func (a *adminImpl) updateApp (name string, appView *model.ApplicationView) (*model.ApplicationView, error) {
	log.L().Info("update app request", log.Any("info", appView))

	res, err := a.app.Update("", name, appView)
	if err != nil {
		return nil, err
	}

	log.L().Info("update app response", log.Any("info", res))

	return res, err
}