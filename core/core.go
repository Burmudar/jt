package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

var DefaultTmplFuncMap = template.FuncMap{
	"json": func(k string, v json.RawMessage) json.RawMessage {
		m := make(map[string]json.RawMessage)
		err := json.Unmarshal(v, &m)
		if err != nil {
			panic(err)
		}
		return jsonKey(k, m)
	},
	"str": func(v []byte) string {
		if bytes.IndexByte(v, '"') == 0 {
			v = bytes.Trim(v, "\"")
		}
		return fmt.Sprintf("%s\n", string(v))
	},
}

func readInput(path string) ([]byte, error) {
	var f *os.File
	if path == "" {
		f = os.Stdin
	} else if fp, err := os.Open(path); err != nil {
		return nil, err
	} else {
		f = fp
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func readTemplate(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	return data, err
}

func ToJSONMap(data []byte) (map[string]json.RawMessage, error) {
	contentMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &contentMap); err != nil {
		return nil, err
	}
	return contentMap, nil
}

func traverseJSON(k string, jsonMap map[string]json.RawMessage) map[string]json.RawMessage {
	if msg, ok := jsonMap[k]; !ok {
		panic(fmt.Errorf("not valid key in json: %q", k))
	} else {
		m := make(map[string]json.RawMessage)
		err := json.Unmarshal(msg, &m)
		if err != nil {
			panic(err)
		}
		return m
	}
}

func jsonKey(key string, jsonMap map[string]json.RawMessage) json.RawMessage {
	cmp := strings.Split(key, ".")
	for i := 0; i < len(cmp)-1; i++ {
		jsonMap = traverseJSON(cmp[i], jsonMap)
	}
	return jsonMap[cmp[len(cmp)-1]]
}

func NewTemplate(data string) (*template.Template, error) {
	tmpl, err := template.New("json").Funcs(DefaultTmplFuncMap).Parse(data)
	return tmpl, err
}

func ApplyTemplate(tmpl *template.Template, ctx map[string]json.RawMessage) error {
	return tmpl.Execute(os.Stdout, ctx)
}
