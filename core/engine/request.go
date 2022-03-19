package engine

import (
	"baetyl-simulator/constants"
	"io/ioutil"
	"net/http"
	"os"

	"baetyl-simulator/core/utils"
	"encoding/json"
	"time"

	"baetyl-simulator/errors"
	"baetyl-simulator/middleware/log"
	specV1 "baetyl-simulator/spec/v1"
)

func (e *engine) Report(r specV1.Report) (specV1.Desire, error) {
	msg := &specV1.Message{
		Kind:     specV1.MessageReport,
		Metadata: map[string]string{"source": os.Getenv(constants.KEY_SRV_NAME)},
		Content:  specV1.LazyValue{Value: r},
	}
	msg.Metadata["node"] = e.node.Name
	res, err := e.Request(msg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	e.log.Debug("sync reports cloud shadow", log.Any("report", msg))
	var desire specV1.Desire
	err = res.Content.Unmarshal(&desire)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return desire, nil
}



func (e *engine) Request(msg *specV1.Message) (*specV1.Message, error) {
	e.log.Debug("http link send request", log.Any("message", msg))
	pld, err := json.Marshal(msg.Content)
	if err != nil {
		return nil, errors.Trace(err)
	}
	var data []byte
	var resp *http.Response
	res := &specV1.Message{Kind: msg.Kind}
	switch msg.Kind {
	case specV1.MessageReport:
		start := time.Now()
		resp, err = e.syncService.Report(pld)
		e.log.Debug("report finish", log.Any("start", start.Local()),
			log.Any("cost", time.Since(start).Seconds()))
		if err != nil {
			return nil, errors.Trace(err)
		}
	case specV1.MessageDesire:
		start := time.Now()
		resp, err = e.syncService.Report(pld)
		e.log.Debug("report finish", log.Any("start", start.Local()),
			log.Any("cost", time.Since(start).Seconds()))
		if err != nil {
			return nil, errors.Trace(err)
		}
	default:
		return nil, errors.Errorf("unsupported message kind")
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Trace(err)
	}
	data, err = utils.ParseEnv(data)
	if err != nil {
		return nil, errors.Trace(err)
	}
	res.Content.SetJSON(data)
	e.log.Debug("http link receive response", log.Any("message", res))
	return res, nil
}