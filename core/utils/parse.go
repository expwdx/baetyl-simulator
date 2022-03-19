package utils

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func ParseEnv(data []byte) ([]byte, error) {
	text := string(data)
	envs := os.Environ()
	envMap := make(map[string]string)
	for _, s := range envs {
		t := strings.Split(s, "=")
		envMap[t[0]] = t[1]
	}
	tmpl, err := template.New("template").Option("missingkey=error").Parse(text)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	buffer := bytes.NewBuffer(nil)
	err = tmpl.Execute(buffer, envMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return buffer.Bytes(), nil
}
