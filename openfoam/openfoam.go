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

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	// Plugin Name
	Plugin = "openfoam"
	// Version of plugin
	Version = 3
	// NsMetricPosition from openfoam ns
	NsMetricPosition = 2
	// NsSubMetricPosition from openfoam ns
	NsSubMetricPosition = 3
	// Vendor Name
	Vendor = "intel"
)

// OpenFoam struct
type OpenFoam struct {
}

// OpenFoamMetrics contains list of available metrics
var OpenFoamMetrics = []string{"k", "p", "Ux", "Uy", "Uz", "omega"}

// CollectMetrics returns collected metrics
func (OpenFoam) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	metrics := make([]plugin.Metric, len(mts))
	webServerIP, err := mts[0].Config.GetString("webServerIP")
	if err != nil {
		return nil, err
	}
	webServerPort, err := mts[0].Config.GetInt("webServerPort")
	if err != nil {
		return nil, err
	}
	webServerFilePath, err := mts[0].Config.GetString("webServerFilePath")
	if err != nil {
		return nil, err
	}
	timeOut, err := mts[0].Config.GetInt("timeOut")
	if err != nil {
		return nil, err
	}
	webServerURL := openFoamURL(webServerIP, webServerPort)

	var tags map[string]string
	tags = make(map[string]string)
	tags["hostname"] = webServerIP
	for i, p := range mts {

		metric, err := openFoamStat(p.Namespace, webServerURL, webServerFilePath, timeOut)
		if err != nil {
			return nil, err
		}
		metrics[i] = *metric
		metrics[i].Tags = tags

	}
	return metrics, nil
}

// GetConfigPolicy returns a config policy
func (OpenFoam) GetConfigPolicy() (plugin.ConfigPolicy, error) {

	policy := plugin.NewConfigPolicy()
	ns := []string{"intel", "openfoam"}
	policy.AddNewStringRule(ns, "webServerIP", true)
	policy.AddNewIntRule(ns, "webServerPort", true, plugin.SetDefaultInt(8000))
	policy.AddNewStringRule(ns, "webServerFilePath", true)
	policy.AddNewIntRule(ns, "timeOut", false, plugin.SetDefaultInt(2))

	return *policy, nil

}

// GetMetricTypes returns metric types that can be collected
func (OpenFoam) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	var metrics []plugin.Metric
	for _, metricType := range OpenFoamMetrics {
		for _, m := range []string{"initial", "final"} {
			ns := plugin.NewNamespace(Vendor, Plugin, metricType, m)
			metric := plugin.Metric{Namespace: ns, Version: Version}
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

func openFoamStat(ns plugin.Namespace, webServerURL string, path string, timeOut int64) (*plugin.Metric, error) {
	data, err := getOpenFoamLog(webServerURL, path, timeOut)
	if err != nil {
		return nil, err
	}
	openFoamLog := string(data[:])
	value, err := getLastValue(ns, strings.Split(openFoamLog, "\n"))
	if err != nil {
		return nil, err
	}

	return &plugin.Metric{
		Namespace: ns,
		Data:      value,
		Timestamp: time.Now(),
	}, nil

}

func getLastValue(ns plugin.Namespace, data []string) (float64, error) {
	searchFor := ns.Strings()[NsMetricPosition]
	switch searchFor {
	case "k":
		searchFor = "for k"
	case "p":
		searchFor = "for p"
	}

	for i := len(data) - 1; i >= 0; i-- {
		if strings.Contains(data[i], searchFor) {
			value, err := getPositionalValue(ns.Strings()[NsSubMetricPosition], data[i])

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

func getOpenFoamLog(webServerURL string, path string, timeOut int64) ([]byte, error) {
	response, err := openFoamWebCall(webServerURL, path, timeOut)
	if err != nil {
		log.Println("Error in call getOpenFoamLog", err)
		return []byte{}, err
	}

	return response, nil
}
