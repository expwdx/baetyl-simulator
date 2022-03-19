package service

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type TemplateService interface {
	Execute(name, text string, params map[string]interface{}) ([]byte, error)
	GetTemplate(filename string) (string, error)
	ParseTemplate(filename string, params map[string]interface{}) ([]byte, error)
	UnmarshalTemplate(filename string, params map[string]interface{}, out interface{}) error
}

// TemplateServiceImpl is a service to read and parse template files.
type TemplateServiceImpl struct {
	path  string
	funcs map[string]interface{}
}

func NewTemplateService(path string, funcs map[string]interface{}) (TemplateService, error) {
	return &TemplateServiceImpl{
		path:  path,
		funcs: funcs,
	}, nil
}

func (s *TemplateServiceImpl) Execute(name, text string, params map[string]interface{}) ([]byte, error) {
	t, err := template.New(name).Option("missingkey=error").Funcs(s.funcs).Parse(text)
	if err != nil {
		return nil, errors.Wrap(err, "execute template error")
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, params)
	if err != nil {
		return nil, errors.Wrap(err, "execute template error")
	}
	return buf.Bytes(), nil
}

func (s *TemplateServiceImpl) GetTemplate(filename string) (string, error) {
	file := s.path + filename
	res, err := ioutil.ReadFile(file)
	if err != nil {
		return "", errors.Wrap(err, "read template error")
	}

	return string(res), nil
}

func (s *TemplateServiceImpl) ParseTemplate(filename string, params map[string]interface{}) ([]byte, error) {
	tl, err := s.GetTemplate(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := s.Execute(filename, tl, params)
	if err != nil {
		return nil, errors.Wrap(err, "execute template error")
	}
	return data, nil
}

func (s *TemplateServiceImpl) UnmarshalTemplate(filename string, params map[string]interface{}, out interface{}) error {
	tp, err := s.ParseTemplate(filename, params)
	if err != nil {
		return errors.WithStack(err)
	}
	err = yaml.Unmarshal(tp, out)
	if err != nil {
		return errors.Wrap(err, "template unmarshal error")
	}
	return nil
}
