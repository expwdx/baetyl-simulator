package service


import (
    "baetyl-simulator/baetyl"
    "baetyl-simulator/config"
    "baetyl-simulator/constants"
    "context"
    "net/http"
)

type SyncService interface {
    Report (body interface{}) (*http.Response, error)
    Desire (body interface{}) (*http.Response, error)
}

type syncService struct {
    client      baetyl.BaetylHttpClient
}

func NewSyncService (ctx context.Context, config *config.ServerConfig) (SyncService, error) {
    syncClient, err := baetyl.NewBaetylClient(ctx, config)
    if err != nil {
        return nil, err
    }

    return &syncService{
        client: syncClient,
    }, nil
}

func (s *syncService) Report (body interface{}) (*http.Response, error) {
    resp, err := s.client.Post(constants.URL_PATH_BAETYL_SYNC_REPORT, nil, body)
    if err != nil {
        return nil, err
    }

    return resp, nil
}

func (s *syncService) Desire (body interface{}) (*http.Response, error)  {
    resp, err := s.client.Post(constants.URL_PATH_BAETYL_SYNC_REPORT, nil, body)
    if err != nil {
        return nil, err
    }

    return resp, nil
}
