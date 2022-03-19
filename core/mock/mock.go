package mock

import (
    "baetyl-simulator/constants"
    "baetyl-simulator/middleware/log"
    "baetyl-simulator/service"
    "encoding/base64"
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/pkg/errors"

    "baetyl-simulator/config"
)

type MockService interface {
    InitData() error
    CleanData() error
    CrawlNodeCertToLocal (namespace, name string) error
}

type mockService struct {
    mockCfg           *config.MockConfig
    nodeService       service.NodeService
    deploymentService service.DeploymentService
}

func NewMockService(config *config.Config) (MockService, error) {
    nodeService, err := service.NewNodeService(&config.Cloud.Admin)
    if err != nil {
        return nil, err
    }

    deploymentService, err := service.NewDeploymentService(&config.Cloud.Init)
    if err != nil {
        return nil, err
    }

    return &mockService{
        mockCfg: &config.Mock,
        nodeService:       nodeService,
        deploymentService: deploymentService,
    }, nil
}

func (s *mockService) InitData() error {
    defer func() {
        log.L().Info("init resource data finished.")
    }()

    nodePrefix := s.mockCfg.NodeNamePrefix
    total := s.mockCfg.NodeCount
    namespace := s.mockCfg.Namespace
    nodeLabels := s.mockCfg.NodeLabels
    startNo := s.mockCfg.NodeStartNo

    if total == 0 {
        return errors.New("node count cannot be 0")
    }

    groupElemNum := 100
    group, tail := total/groupElemNum, total%groupElemNum
    if tail > 0 {
        group += 1
    }

    wg := sync.WaitGroup{}

    for j := 0 ; j < group; j++ {
        wg.Add(1)
        gi := j
        go func() {
            defer wg.Done()

            var gn int
            if tail > 0 && gi == (group-1) {
                gn = tail
            } else {
                gn = groupElemNum
            }

            for i := 0; i < gn; i++ {
                nodeNo := startNo + gi*groupElemNum + i
                nodeName := fmt.Sprintf("%s-%d", nodePrefix, nodeNo)
                if _, err := s.nodeService.Create(namespace, nodeName, nodeLabels); err != nil {
                    log.L().Error("create node fail", log.Any("error", err))
                    return
                }
                log.L().Info("create node", log.Any("namespace", namespace), log.Any("node", nodeName))

                if err := s.CrawlNodeCertToLocal(namespace, nodeName); err != nil {
                    log.L().Error("crawling node cert fail", log.Any("error", err))
                    return
                }
                time.Sleep(5*time.Millisecond)
            }
        }()
    }

    wg.Wait()

    return nil
}

func (s *mockService) CleanData() error {
    nodePrefix := s.mockCfg.NodeNamePrefix
    startNo := s.mockCfg.NodeStartNo
    total := s.mockCfg.NodeCount
    namespace := s.mockCfg.Namespace

    for i := startNo; i < (startNo + total); i++ {
        nodeName := fmt.Sprintf("%s-%d", nodePrefix, i)
        if err := s.nodeService.Delete(namespace, nodeName); err != nil {
            return errors.Wrap(err, "delete node fail")
        }
    }

    return nil
}

func (s *mockService) CrawlNodeCertToLocal(namespace, name string) error {
    initScript, err := s.nodeService.GetInstallScript(namespace, name)
    if err != nil {
        return errors.Wrap(err, "get install script error")
    }

    token := strings.Split(strings.Split(initScript, "token=")[1], "&mode=")[0]
    initDeployment, err := s.deploymentService.GetInitDeployment(token)
    if err != nil {
        return errors.Wrap(err, "get init deployment yaml error")
    }

    lineChar, err := getNewlineChar(initDeployment)
    if err != nil {
        panic(err)
    }

    certKeyList := strings.Split(initDeployment, fmt.Sprintf("%s  client.key: '", lineChar))
    certKey := strings.Split(certKeyList[1], fmt.Sprintf("'%s  ca.pem: '", lineChar))[0]
    certKeyBytes, err := base64.StdEncoding.DecodeString(certKey)
    if err != nil {
        return errors.Wrap(err, "decode cert key error")
    }

    certPemList := strings.Split(initDeployment, fmt.Sprintf("%s  client.pem: '", lineChar))
    certPem := strings.Split(certPemList[1], fmt.Sprintf("'%s  client.key: '", lineChar))[0]
    certPemBytes, err := base64.StdEncoding.DecodeString(certPem)
    if err != nil {
        return errors.Wrap(err, "decode cert key error")
    }

    log.L().Debug("certs", log.Any("node", name), log.Any("client.key", certKeyBytes),
        log.Any("client.pem", certPemBytes))

    if err = writeCerts(name, certKeyBytes, certPemBytes); err != nil {
        return errors.Wrap(err, "write certs to local fail")
    }

    return nil
}

func writeCerts(name string, certKey, certPem []byte) error {
    keyFile := fmt.Sprintf(constants.NODE_CERT_KEY, name)
    pemFile := fmt.Sprintf(constants.NODE_CERT_CERT, name)
    if err := os.MkdirAll(fmt.Sprintf(constants.NODE_CERTS_DIR, name), 0700); err != nil {
        return errors.Wrap(err, "mkdir error")
    }

    if err := ioutil.WriteFile(keyFile, certKey, 0600); err != nil {
        return errors.Wrap(err, "write client.key to local file error")
    }

    if err := ioutil.WriteFile(pemFile, certPem, 0600); err != nil {
        return errors.Wrap(err, "write client.key to local file error")
    }

    return nil
}

func getNewlineChar(s string) (string, error) {
    l := strings.Split(s, "\r\n  client.key: '")
    if len(l) == 2 {
        return "\r\n", nil
    }
    l = strings.Split(s, "\n  client.pem: '")
    if len(l) == 2 {
        return "\n", nil
    }

    return "", errors.Errorf("cannot parse %s", s)
}


