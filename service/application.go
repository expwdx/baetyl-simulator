package service

import (
	"baetyl-simulator/baetyl"
	"baetyl-simulator/config"
	"baetyl-simulator/constants"
	"baetyl-simulator/model"
)

type ApplicationService interface {
	Get(namespace, name string) (*model.ApplicationView, error)
	Create (namespace, app *model.ApplicationView) (*model.ApplicationView, error)
	Update (namespace, name string, app *model.ApplicationView) (*model.ApplicationView, error)
}

type applicationService struct {
	client      baetyl.BaetylHttpClient
}

func NewApplicationService (config *config.ServerConfig) (ApplicationService, error) {
	adminClient, err := baetyl.NewBaetylClient(nil, config)
	if err != nil {
		return nil, err
	}

	return &applicationService{
		client: adminClient,
	}, nil
}

func (s *applicationService) Get(namespace, name string) (*model.ApplicationView, error) {
	var app *model.ApplicationView
	pathField := baetyl.PathField{"name": name}

	resp, err := s.client.Get(constants.URL_PATH_BAETYL_CLOUD_APP, pathField, nil)
	if err != nil {
		return nil, err
	}

	err = s.client.Read(resp, &app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *applicationService) Create (namespace, app *model.ApplicationView) (*model.ApplicationView, error) {
	resp, err := s.client.Post(constants.URL_PATH_BAETYL_CLOUD_APPS, nil, app)
	if err != nil {
		return nil, err
	}

	err = s.client.Read(resp, &app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (s *applicationService) Update(namespace, name string, app *model.ApplicationView) (*model.ApplicationView, error) {
	pathField := baetyl.PathField{"name": name}
	resp, err := s.client.Put(constants.URL_PATH_BAETYL_CLOUD_APP, pathField, app)
	if err != nil {
		return nil, err
	}

	err = s.client.Read(resp, &app)
	if err != nil {
		return nil, err
	}

	return app, nil
}