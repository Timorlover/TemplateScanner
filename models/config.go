package models

import (
	"TemplateScanner/util"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	LoadTemplate []string `yaml:"load_template,omitempty"`
	Timeout      int      `yaml:"timeout,omitempty"`
	Telegram     struct {
		Token  string `yaml:"token,omitempty"`
		ChatID int    `yaml:"chatid,omitempty"`
	} `yaml:"telegram,omitempty"`
	InputArguments []map[string]string `yaml:"input_arguments,omitempty"`
	Thread         int                 `yaml:"thread"`
}

func (c Config) Parse(data []byte) (*Config, error) {
	err := yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal Data Error:%s", err)
	}
	return &c, nil
}

func LoadConfigYaml() (*Config, error) {
	if util.IsEmptyFile("config.yaml") {
		data, err := os.ReadFile("config.yaml")
		if err != nil {
			return nil, fmt.Errorf(color.RedString("Load ConfigYaml Err:", err))
		}
		var parse Config
		parseData, err := parse.Parse(data)
		if err != nil {
			fmt.Println(color.RedString("Parse ConfigYaml Err:", err))
			return nil, fmt.Errorf(color.RedString("Parse ConfigYaml Err:", err))
		} else {
			return parseData, nil
		}
	}
	return nil, fmt.Errorf(color.RedString("No ConfigYaml!"))
}

func GenerateConfigYaml() (*Config, error) {
	p := Config{
		LoadTemplate: nil,
		Timeout:      300,
		Telegram: struct {
			Token  string `yaml:"token,omitempty"`
			ChatID int    `yaml:"chatid,omitempty"`
		}{},
		InputArguments: nil,
		Thread:         3,
	}

	var templateList []string
	data := LoadAllScannerYaml()
	for _, v := range data {
		//fmt.Println(v.ScannerName)
		templateList = append(templateList, v.ScannerName)
	}

	//p.LoadTemplate = []string{"xray", "nuclei", "bbscan"}
	p.LoadTemplate = templateList
	p.InputArguments = []map[string]string{}
	y, err := yaml.Marshal(&p)
	if err != nil {
		log.Fatalf(color.RedString("Marshal ConfigYaml Err:%v"), err)
	}

	// Write YAML to file
	err = ioutil.WriteFile("config.yaml", y, 0644)
	if err != nil {
		log.Fatalf(color.RedString("Write Data To ConfigYaml Err: %v"), err)
	} else {
		fmt.Println(color.GreenString("New ConfigYaml Will Generate!"))
	}
	return &p, fmt.Errorf(color.GreenString("New ConfigYaml Will Generate!"))
}
