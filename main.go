package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Burmudar/jt/core"
)

var BuildDate string = "UNKOWN"
var GitCommit string = "UNKOWN"
var GitBranch string = "UNKOWN"
var Version string = "UNKOWN"

type config struct {
	input          *string
	template       *string
	inlineTemplate *string
}

func readInput(path string) ([]byte, error) {
	var f *os.File
	if path == "" || path == "stdin" {
		f = os.Stdin
	} else if fp, err := os.Open(path); err != nil {
		return nil, err
	} else {
		f = fp
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func readTemplate(cfg config) (string, error) {
	var tmplData string
	if *cfg.inlineTemplate == "" {
		data, err := ioutil.ReadFile(*cfg.template)
		if err != nil {
			return "", err
		}
		tmplData = string(data)
	} else {
		tmplData = *cfg.inlineTemplate
	}
	return tmplData, nil
}

func innerMain(cfg config) error {
	content, err := readInput(*cfg.input)
	if err != nil {
		return err
	}

	jsonMap, err := core.ToJSONMap(content)
	if err != nil {
		fmt.Printf("failed to decode json: %v\n", err)
		return fmt.Errorf("failed to decode json")
	}
	data, err := readTemplate(cfg)
	if err != nil {
		return fmt.Errorf("failed to read template")
	}
	template, err := core.NewTemplate("",data)
	if err != nil {
		return err
	}
	return core.ApplyTemplate(template, nil, jsonMap)
}

func main() {
	cfg := config{}
	cfg.input = flag.String("input-file", "stdin", "specify the file to read input from")
	cfg.template = flag.String("template", "", "Specify the template file to apply")
	cfg.inlineTemplate = flag.String("inline-template", "", "Specify the template inline to apply")
	version := flag.Bool("version", false, "print the version")
	flag.Parse()

	if *version {
		fmt.Printf("JSON-Template - Apply a Golang template to JSON data\nDate: %s\nCommit: %s\nBranch: %s\nVersion: %s\n", BuildDate, GitCommit, GitBranch, Version)
	} else if err := innerMain(cfg); err != nil {
		fmt.Printf("unexpected error: %v\n", err)
	}
}
