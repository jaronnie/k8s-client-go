package util

import (
	"regexp"

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
