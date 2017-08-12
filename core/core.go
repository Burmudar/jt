package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// DefaultTmplFuncMap are the additional functions that are available to the templates that are created by this package
// The following functions are added:
// - json: takes a key and a jsonRawMessage. The key is used to traverse the given json.RawMessage
// - str: converts bytes to string
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
		return fmt.Sprintf("%s", string(trimQuotes(v)))
	},
}

func trimQuotes(v []byte) []byte {
	if bytes.IndexByte(v, '"') == 0 {
		v = bytes.Trim(v, "\"")
	}
	return v
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

// ToJSONMap unmarshalls the given data into a map where the keys are strings and the values are json.RawMessage.
// Data is expected to be valid JSON
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

func arrayJSON(key string, jsonMap map[string]json.RawMessage) json.RawMessage {
	index, err := strconv.Atoi(string(key[len(key)-2]))
	if err != nil {
		panic(err)
	}
	key = string(key[:len(key)-3])
	v := jsonMap[key]
	arr := make([]json.RawMessage, 0)
	json.Unmarshal(v, &arr)
	return arr[index]
}

func jsonKey(key string, jsonMap map[string]json.RawMessage) json.RawMessage {
	cmp := strings.Split(key, ".")
	isArrIndexKey := func(k string) bool {
		size := len(k)
		leftB := strings.LastIndex(k, "[")
		rightB := strings.LastIndex(k, "]")
		//are the brackets at the end of the key or somewhere before that... which doesn't make sense
		if leftB == size-3 && rightB == size-1 {
			return true
		}
		return false
	}
	for i := 0; i < len(cmp)-1; i++ {
		if isArrIndexKey(cmp[i]) {
			if m, err := ToJSONMap(arrayJSON(cmp[i], jsonMap)); err != nil {
				panic(err)
			} else {
				jsonMap = m
			}
		} else {
			jsonMap = traverseJSON(cmp[i], jsonMap)
		}
	}
	key = cmp[len(cmp)-1]
	if isArrIndexKey(key) {
		return arrayJSON(key, jsonMap)
	}
	return jsonMap[key]
}

// NewTemplate creates a new template using the provided data as the template data to be parsed.
// If no name is provided the default name is "json"
func NewTemplate(name, data string) (*template.Template, error) {
	if name == "" {
		name = "json"
	}
	tmpl, err := template.New(name).Funcs(DefaultTmplFuncMap).Parse(data)
	return tmpl, err
}

// ApplyTemplate takes the given template and applies the map to the template. No additional processing is done to the map.
// If writer is nil, the template is executed to Stdout
func ApplyTemplate(tmpl *template.Template, w io.Writer, ctx map[string]json.RawMessage) error {
	if w == nil {
		w = os.Stdout
	}
	return tmpl.Execute(w, ctx)
}
