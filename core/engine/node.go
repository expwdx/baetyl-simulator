package engine

import (
	"baetyl-simulator/common"
	"baetyl-simulator/errors"
	specV1 "baetyl-simulator/spec/v1"
)

func (e *engine) GetShadowDelta() (specV1.Delta, error) {
	var node *specV1.Node
	err := common.Cache.Get(e.node.Name, &node)
	if err != nil {
		return nil, err
	}

	delta, err := node.Desire.DiffWithNil(node.Report)
	if err != nil {
		return nil, err
	}

	return delta, nil
}

func (e *engine) MergeReport(reported specV1.Report) error {
	var node *specV1.Node
	err := common.Cache.Get(e.node.Name, &node)
	if err != nil {
		return err
	}

	if node == nil {
		return errors.New("node not found in cache")
	}

	if node.Report == nil {
		node.Report = reported
	} else {
		err = node.Report.Merge(reported)
		if err != nil {
			return errors.Trace(err)
		}
	}

	err = common.Cache.Set(e.node.Name, node, -1)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (e *engine)MergeDesire(desired specV1.Desire) error {
	var node *specV1.Node
	err := common.Cache.Get(e.node.Name, &node)
	if err != nil {
		return err
	}

	if node == nil {
		return errors.New("node not found in cache")
	}

	if node.Desire == nil {
		node.Desire = desired
	} else {
		err = node.Desire.Merge(desired)
		if err != nil {
			return errors.Trace(err)
		}
	}

	err = common.Cache.Set(e.node.Name, node, -1)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}