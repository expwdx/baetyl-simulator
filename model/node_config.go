package model

import (
	"time"
)

// NodeConfig node config
type NodeConfig struct {
	NodeName   string    `json:"nodeName,omitempty"`
	Namespace  string    `json:"namespace,omitempty"`
	Name       string    `json:"key,omitempty"`
	Value      string    `json:"value,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
}

type NodeConfigFilter struct {
	*Filter
	Name     string `form:"lookupField,omitempty" json:"lookupField,omitempty"`
}
