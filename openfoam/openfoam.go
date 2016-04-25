/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package openfoam

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	// Name of plugin
	Name = "openfoam"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

// Meta returns plugin meta data info
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(Name, Version, Type, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// OpenFoam struct
type OpenFoam struct {
}

// NewOpenFoamCollector returns new Collector instance
func NewOpenFoamCollector() *OpenFoam {
	return &OpenFoam{}

}

func joinNamespace(ns []string) string {
	return "/" + strings.Join(ns, "/")
}

// OpenFoamMetrics contains list of available metrics
var OpenFoamMetrics = []string{"k", "p", "Ux", "Uy", "Uz", "omega"}

// CollectMetrics returns collected metrics
func (p *OpenFoam) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	metrics := make([]plugin.PluginMetricType, len(mts))

	webServerIP := mts[0].Config().Table()["webServerIP"].(ctypes.ConfigValueStr).Value
	webServerPort := mts[0].Config().Table()["webServerPort"].(ctypes.ConfigValueInt).Value
	webServerFilePath := mts[0].Config().Table()["webServerFilePath"].(ctypes.ConfigValueStr).Value
	webServerURL := openFoamURL(webServerIP, webServerPort)

	var tags map[string]string
	tags = make(map[string]string)
	tags["hostname"] = webServerIP

	for i, p := range mts {

		metric, err := openFoamStat(p.Namespace(), webServerURL, webServerFilePath)
		if err != nil {
			return nil, err
		}
		metrics[i] = *metric
		metrics[i].Tags_ = tags

	}
	return metrics, nil
}

// GetConfigPolicy returns a config policy
func (p *OpenFoam) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	webServerIP, err := cpolicy.NewStringRule("webServerIP", true)
	handleErr(err)
	webServerIP.Description = "OpenFoam hostname/ip address"
	config.Add(webServerIP)
	webServerPort, err := cpolicy.NewIntegerRule("webServerPort", false, 8000)
	handleErr(err)
	webServerPort.Description = "webServerger port / default 8000"
	config.Add(webServerPort)
	webServerFilePath, err := cpolicy.NewStringRule("webServerFilePath", true)
	handleErr(err)
	webServerIP.Description = "File location"
	config.Add(webServerFilePath)

	cp.Add([]string{""}, config)
	return cp, nil

}

// GetMetricTypes returns metric types that can be collected
func (p *OpenFoam) GetMetricTypes(cfg plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	var metrics []plugin.PluginMetricType
	for _, metricType := range OpenFoamMetrics {
		metrics = append(metrics, plugin.PluginMetricType{Namespace_: []string{"intel", "openfoam", metricType, "initial"}})
		metrics = append(metrics, plugin.PluginMetricType{Namespace_: []string{"intel", "openfoam", metricType, "final"}})
	}
	return metrics, nil
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}

func openFoamStat(ns []string, webServerURL string, path string) (*plugin.PluginMetricType, error) {
	data, err := getOpenFoamLog(webServerURL, path)
	if err != nil {
		return nil, err
	}
	openFoamLog := string(data[:])
	value, err := getLastValue(ns, strings.Split(openFoamLog, "\n"))
	if err != nil {
		return nil, err
	}

	return &plugin.PluginMetricType{
		Namespace_: ns,
		Data_:      value,
		Timestamp_: time.Now(),
	}, nil

}

func getLastValue(ns []string, data []string) (float64, error) {
	searchFor := ns[2]
	switch searchFor {
	case "k":
		searchFor = "for k"
	case "p":
		searchFor = "for p"
	}

	for i := len(data) - 1; i >= 0; i-- {
		if strings.Contains(data[i], searchFor) {
			value, err := getPositionalValue(ns[3], data[i])

			if err != nil {
				log.Fatal(err)
				return 0.0, err
			}
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				log.Fatal(err)
				return 0.0, err
			}
			return floatValue, nil

		}

	}
	return 0.0, fmt.Errorf("Can't find data in OpenFoamLog")

}

func getPositionalValue(position string, data string) (string, error) {
	line := regexp.MustCompile("^(.+):  Solving for (.+), Initial residual = (?P<initial>.+), Final residual = (?P<final>.+), No Iterations (.+)$")
	switch position {
	case "initial":

		values := line.FindAllStringSubmatch(data, -1)[0]
		if len(values) >= 4 {
			return values[3], nil
		}

	case "final":

		values := line.FindAllStringSubmatch(data, -1)[0]
		if len(values) >= 4 {
			return values[4], nil
		}

	}
	return "", fmt.Errorf("Can't find initial or final residual")

}

func getOpenFoamLog(webServerURL string, path string) ([]byte, error) {
	response, err := openFoamWebCall(webServerURL, path)
	if err != nil {
		log.Println("Error in call getOpenFoamLog", err)
		return []byte{}, err
	}

	return response, nil
}
