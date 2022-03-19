package service

import (
    "baetyl-simulator/baetyl"
    "baetyl-simulator/config"
    "baetyl-simulator/constants"
    "baetyl-simulator/model"
    specV1 "baetyl-simulator/spec/v1"
    "github.com/pkg/errors"
)

type NodeService interface {
    Get(namespace, name string) (*specV1.Node, error)
    List(namespace string) (*model.NodeViewList, error)
    Create (namespace, name string, labels map[string]string) (*specV1.Node, error)
    Delete (namespace, name string) error
    GetInstallScript (namespace, name string) (string, error)
}

type nodeService struct {
    client      baetyl.BaetylHttpClient
}

func NewNodeService (config *config.ServerConfig) (NodeService, error) {
    adminClient, err := baetyl.NewBaetylClient(nil, config)
    if err != nil {
        return nil, err
    }

    return &nodeService{
        client: adminClient,
    }, nil
}

func (s *nodeService) Get(namespace, name string) (*specV1.Node, error) {
    var node *specV1.Node
    pathField := baetyl.PathField{"name": name}

    resp, err := s.client.Get(constants.URL_PATH_BAETYL_CLOUD_NODE, pathField, nil)
    if err != nil {
        return nil, err
    }

    err = s.client.Read(resp, &node)
    if err != nil {
        return nil, err
    }

    return node, nil
}

func (s *nodeService) List(namespace string) (*model.NodeViewList, error) {
    var nodes *model.NodeViewList

    resp, err := s.client.Get(constants.URL_PATH_BAETYL_CLOUD_NODES, nil, nil)
    if err != nil {
        return nil, err
    }

    err = s.client.Read(resp, &nodes)
    if err != nil {
        return nil, err
    }

    return nodes, nil
}

func (s *nodeService) Create (namespace, name string, labels map[string]string) (*specV1.Node, error) {
    node := &specV1.Node{
        Name: name,
        Namespace: namespace,
        Labels: labels,
    }

    resp, err := s.client.Post(constants.URL_PATH_BAETYL_CLOUD_NODES, nil, node)
    if err != nil {
        return nil, err
    }

    err = s.client.Read(resp, &node)
    if err != nil {
        return nil, err
    }

    return node, nil
}

func (s *nodeService) Delete (namespace, name string) error {
    pathField := baetyl.PathField{"name": name}
    resp, err := s.client.Delete(constants.URL_PATH_BAETYL_CLOUD_NODE, pathField)
    if err != nil {
        return err
    }

    if resp.StatusCode != 200 {
        return errors.Errorf("delete node fail. namespace: %s, node: %s", namespace, name)
    }
    _, err = s.client.ReadMap(resp)
    if err != nil {
        return err
    }

    return nil
}

func (s *nodeService) GetInstallScript (namespace, name string) (string, error) {
    pathFields := baetyl.PathField{"name": name}
    resp, err := s.client.Get(constants.URL_PATH_BAETYL_CLOUD_NODE_INIT, pathFields, nil)
    if err != nil {
        return "", err
    }

    var scriptInfo map[string]interface{}
    scriptInfo, err = s.client.ReadMap(resp)
    if err != nil {
        return "", err
    }

    res, ok := scriptInfo["cmd"]
    if !ok {
        return "", errors.New("cmd not found in response.")
    }

    return res.(string), nil
}
