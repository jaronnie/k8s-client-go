package util

import (
	"bytes"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/flosch/pongo2/v4"
	"github.com/pkg/errors"
	yaml2 "k8s.io/apimachinery/pkg/util/yaml"
)

func Yaml2Jsons(data []byte) ([][]byte, error) {
	jsons := make([][]byte, 0)
	yamls := RegSplit(data, "\n---\n")
	for _, v := range yamls {
		if len(v) == 0 {
			continue
		}
		obj, err := yaml2.ToJSON(v)
		if err != nil {
			return nil, err
		}
		jsons = append(jsons, obj)
	}
	return jsons, nil
}

func RegSplit(t []byte, reg string) [][]byte {
	re := regexp.MustCompile(reg)
	split := re.Split(string(t), -1)
	set := [][]byte{}
	for i := range split {
		set = append(set, []byte(split[i]))
	}
	return set
}

// template
// {{ .Namespace }} --> default
func ParseTemplateYAML(data interface{}, tplT []byte) ([]byte, error) {
	t := template.New("temp").Funcs(sprig.TxtFuncMap())
	t, err := t.Parse(string(tplT))
	if err != nil {
		return nil, err
	}
	var buf = new(bytes.Buffer)
	err = t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

// template
// {{ namespace }} --> default
func ParseTemplateYAML2(data map[string]interface{}, tpl []byte) ([]byte, error) {
	var (
		t   *pongo2.Template
		out string
		err error
	)
	if t, err = pongo2.FromString(string(tpl)); err != nil {
		goto FAIL
	}
	if out, err = t.Execute(data); err != nil {
		goto FAIL
	}
	return []byte(out), nil
FAIL:
	return []byte(""), errors.Wrap(err, "fail to parse template")
}
