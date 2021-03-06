// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"text/template"
)

const yamlAppsTemplate = "{{ printf \"#ONOS Helm Chart Releases\"}}\n\n" +
	"{{range $key, $value := .Entries }}" +
	"{{ printf \"## %s\" $key }}\n\n" +
	"{{range $value}}" +
	"{{printf \"#### Version **%s**\" .Version}}\n" +
	"{{printf \"> Generated %s\" .Created}}\n\n" +
	"{{printf \"App Version **%s**\" .AppVersion}}\n\n" +
	"{{range .Urls}}" +
	"{{printf \"[%s](%s)\" . .}}\n" +
	"{{end}}\n\n" +
	"{{end}}\n" +
	"{{end}}\n"

type Chart struct {
	ApiVersion string `yaml:"apiVersion"`
	AppVersion string `yaml:"appVersion"`
	Version string `yaml:"version"`
	Created string `yaml:"created"`
	Description string `yaml:"description"`
	Urls []string `yaml:"urls"`
}

type IndexYaml struct {
	ApiVersion string `yaml:"apiVersion"`
	Entries    map[string][]Chart `yaml:"entries"`
	Generated string `yaml:"generated"`
}

/**
 * A simple application that takes the generated index.yaml and outputs it in
 * a Markdown format - usually we pipe this to README.md when in the gh-pages branch
 */
func main() {
	file := flag.String("file", "index", "name of YAML file to parse (without extension or path)")
	flag.Parse()
	indexYaml, err := getIndexYaml(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load %s.yaml %s\n", *file, err)
		os.Exit(1)
	}

	var tmplAppsList, _ = template.New("yamlAppsTemplate").Parse(yamlAppsTemplate)
	tmplAppsList.Execute(os.Stdout, indexYaml)
}

func getIndexYaml(location string) (IndexYaml, error) {
	indexYaml := &IndexYaml{}
	viper.SetConfigName(location)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return IndexYaml{}, err
	}

	if err := viper.Unmarshal(indexYaml); err != nil {
		return IndexYaml{}, err
	}

	return *indexYaml, nil
}
