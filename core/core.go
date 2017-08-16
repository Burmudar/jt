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
	"json": jsonKey,
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

func mapJSON(key string, data json.RawMessage) (json.RawMessage, error) {
	if isArrIndexKey(key) {
		key = key[:strings.LastIndex(key, "[")]
	}
	m := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return m[key], nil
}

func arrayJSON(key string, data json.RawMessage) (json.RawMessage, error) {
	index, err := strconv.Atoi(string(key[len(key)-2]))
	if err != nil {
		return nil, err
	}
	arr := make([]json.RawMessage, 0)
	json.Unmarshal(data, &arr)
	if index >= len(data) {
		return nil, fmt.Errorf("key[%s] index out of range: %v", key, string(data))
	}
	return arr[index], nil
}

func isArrIndexKey(k string) bool {
	size := len(k)
	leftB := strings.LastIndex(k, "[")
	rightB := strings.LastIndex(k, "]")
	//are the brackets at the end of the key or somewhere before that... which doesn't make sense
	if leftB == size-3 && rightB == size-1 {
		return true
	}
	return false
}

func ParseHandler(k string) func(v json.RawMessage) (json.RawMessage, error) {
	if isArrIndexKey(k) {
		return func(v json.RawMessage) (json.RawMessage, error) {
			data, err := mapJSON(k[:strings.LastIndex(k,"[")], v)
			if err == nil {
				v = data
			}
			return arrayJSON(k, v)
		}
	}
	return func(v json.RawMessage) (json.RawMessage, error) {
		return mapJSON(k, v)
	}
}

func jsonKey(key string, v json.RawMessage) json.RawMessage {
	cmp := strings.Split(key, ".")
	type parseFn func(data json.RawMessage) (json.RawMessage, error)
	parsers := make([]parseFn, 0)
	for i := 0; i < len(cmp); i++ {
		parsers = append(parsers, ParseHandler(cmp[i]))
	}
	for _, fn := range parsers {
		if len(v) == 0 {
			fmt.Fprint(os.Stderr, fmt.Errorf("no json to parse"))
			break
		}
		data, err := fn(v)
		if err != nil {
			panic(err)
		}
		v = data
	}
	return v
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
func ApplyTemplate(tmpl *template.Template, w io.Writer, ctx map[string]json.RawMessage ) error {
	if w == nil {
		w = os.Stdout
	}
	return tmpl.Execute(w, ctx)
}
