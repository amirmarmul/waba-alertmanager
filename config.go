package main

import (
	"io/ioutil"
	
	"github.com/prometheus/alertmanager/template"
	"gopkg.in/yaml.v2"

	"waba-alertmanager/providers/acs"
)

type ReceiverConf struct {
	Name		string
	Provider	string
	To			[]string
	From		string
	Text		string
	Type		string
}

var config struct {
	Providers	struct {
		Acs		acs.Config
	}
	Receivers	[]ReceiverConf
	Templates	[]string
}

var tmpl *template.Template

func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	tmpl, err = template.FromGlobs(config.Templates)
	return err
}