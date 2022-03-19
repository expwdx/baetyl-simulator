package model

import (
	"reflect"
	"time"

	specV1 "baetyl-simulator/spec/v1"
)

type SecretList struct {
	Total        int `json:"total"`
	*ListOptions `json:",inline"`
	Items        []specV1.Secret `json:"items"`
}

type SecretView struct {
	Name              string            `json:"name,omitempty" validate:"omitempty,resourceName"`
	Namespace         string            `json:"namespace,omitempty"`
	Data              map[string]string `json:"data,omitempty" binding:"required"`
	CreationTimestamp time.Time         `json:"createTime,omitempty"`
	UpdateTimestamp   time.Time         `json:"updateTime,omitempty"`
	Description       string            `json:"description"`
	Version           string            `json:"version,omitempty"`
}

func (s *SecretView) Equal(target *SecretView) bool {
	return reflect.DeepEqual(s.Data, target.Data) &&
		reflect.DeepEqual(s.Description, target.Description)
}

type SecretViewList struct {
	Total        int `json:"total"`
	*ListOptions `json:",inline"`
	Items        []SecretView `json:"items"`
}
