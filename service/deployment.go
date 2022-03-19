package service

import (
    "baetyl-simulator/baetyl"
    "baetyl-simulator/config"
    "baetyl-simulator/constants"
    "github.com/pkg/errors"
    "io/ioutil"
)

type DeploymentService interface {
    GetInitDeployment (token string) (string, error)
}

type deploymentService struct {
    client      baetyl.BaetylHttpClient
}

func NewDeploymentService (config *config.ServerConfig) (DeploymentService, error) {
    adminClient, err := baetyl.NewBaetylClient(nil, config)
    if err != nil {
        return nil, err
    }

    return &deploymentService{
        client: adminClient,
    }, nil
}

func (s *deploymentService) GetInitDeployment (token string) (string, error) {
    params := baetyl.QueryParam{"token": token}
    resp, err := s.client.Get(constants.URL_PATH_BAETYL_INIT_DEPLOYMENT, nil, params)
    if err != nil {
        return "", err
    }

    resBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", errors.Wrap(err, "read response fail.")
    }

    return string(resBytes), nil
}

