package model

import (
	"time"

	specV1 "baetyl-simulator/spec/v1"
)

type Node struct {
	Namespace         string        `json:"namespace,omitempty"`
	Name              string        `json:"name,omitempty" validate:"omitempty,resourceName"`
	CreateTime 		  time.Time		`json:"createTime,omitempty"`
	UpdateTime 		  time.Time		`json:"updateTime,omitempty"`
}

// NodeViewList node view list
type NodeViewList struct {
	Total        int `json:"total"`
	*ListOptions `json:",inline"`
	Items        []specV1.NodeView `json:"items"`
}

// NodeList node list
type NodeList struct {
	Total        int `json:"total"`
	*ListOptions `json:",inline"`
	Items        []specV1.Node `json:"items"`
}

type NodeNames struct {
	Names []string `json:"names,"validate:"maxLength=20"`
}

type NodeProperties struct {
	State NodePropertiesState    `yaml:"state,omitempty" json:"state,omitempty"`
	Meta  NodePropertiesMetadata `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type NodePropertiesState struct {
	Report map[string]interface{} `yaml:"report,omitempty" json:"report,omitempty"`
	Desire map[string]interface{} `yaml:"desire,omitempty" json:"desire,omitempty"`
}

type NodeProperty struct {
	Name    string        `yaml:"name,omitempty" json:"name,omitempty"`
	Type    string        `yaml:"type,omitempty" json:"type,omitempty"`
	Current PropertyValue `yaml:"current,omitempty" json:"current,omitempty"`
	Expect  PropertyValue `yaml:"expect,omitempty" json:"expect,omitempty"`
}

type PropertyValue struct {
	Value      string `yaml:"value,omitempty" json:"value,omitempty"`
	UpdateTime string `yaml:"updateTime,omitempty" json:"updateTime,omitempty"`
}

type NodeMode struct {
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
}

type NodePropertiesMetadata struct {
	ReportMeta map[string]interface{} `yaml:"report,omitempty" json:"report,omitempty"`
	DesireMeta map[string]interface{} `yaml:"desire,omitempty" json:"desire,omitempty"`
}

type NodeCoreConfigs struct {
	Version string `yaml:"version,omitempty" json:"version,omitempty"`
	// unit: seconds
	Frequency int `yaml:"frequency,omitempty" json:"frequency,omitempty"`
	APIPort   int `yaml:"apiport,omitempty" json:"apiport,omitempty"`
}

type NodeCoreVersions struct {
	Versions []string `yaml:"versions,omitempty" json:"versions,omitempty"`
}

type NodeSysAppView struct {
	Name        string `yaml:"name,omitempty" json:"name,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type NodeOptionalSysApps struct {
	Apps []NodeSysAppView `yaml:"apps,omitempty" json:"apps,omitempty"`
}

type NodeSysAppInfo struct {
	Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
	Image       string            `yaml:"image,omitempty" json:"image,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	Programs    map[string]string `yaml:"programs,omitempty" json:"programs,omitempty"`
}