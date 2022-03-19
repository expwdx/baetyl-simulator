package engine

import specV1 "baetyl-simulator/spec/v1"


func (e *engine) genResourceInfos(kind specV1.Kind, infos map[string]string) []specV1.ResourceInfo {
	var crds []specV1.ResourceInfo
	for name, version := range infos {
		crds = append(crds, specV1.ResourceInfo{
			Kind:    kind,
			Name:    name,
			Version: version,
		})
	}
	return crds
}
