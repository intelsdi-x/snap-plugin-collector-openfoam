//
// +build small

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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/jarcoal/httpmock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOpenFoamPlugin(t *testing.T) {
	Convey("Create OpenFoam Collector", t, func() {
		openFoamCol := OpenFoam{}
		Convey("So OpenFoam should not be nil", func() {
			So(openFoamCol, ShouldNotBeNil)
		})
		Convey("So OpenFoam should be of OpenFoam type", func() {
			So(openFoamCol, ShouldHaveSameTypeAs, OpenFoam{})
		})
		Convey("openFoamCol.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := openFoamCol.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a plugin.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
		})
	})

	Convey("Get URI ", t, func() {
		Convey("So should return 10.1.0.1:8000", func() {
			var webServerPort int64
			webServerIP := "10.1.0.1"
			webServerPort = 8000
			uri := openFoamURL(webServerIP, webServerPort)
			So("http://10.1.0.1:8000", ShouldResemble, uri)
		})
	})
	Convey("Get Metrics Types", t, func() {
		openFoamCol := OpenFoam{}
		config := plugin.Config{}
		Convey("So should return 12 types of metrics", func() {
			metrics, err := openFoamCol.GetMetricTypes(config)
			So(len(metrics), ShouldResemble, 12)
			So(err, ShouldBeNil)
		})
		Convey("So should check namespace", func() {
			var waitNamespace string

			metrics, err := openFoamCol.GetMetricTypes(config)
			for _, m := range metrics[0].Namespace.Strings() {
				waitNamespace = fmt.Sprintf("%s/%s", waitNamespace, m)
			}
			wait := regexp.MustCompile(`^/intel/openfoam/k/initial`)
			So(true, ShouldEqual, wait.MatchString(waitNamespace))
			So(err, ShouldBeNil)

		})

	})
	Convey("Collect Metrics", t, func() {
		openFoamCol := OpenFoam{}
		var port, timeOut int64
		var ip, filePath string
		port = 8000
		ip = "192.168.192.200"
		filePath = "test.log"
		timeOut = 1

		config := plugin.Config{
			"webServerIP":       ip,
			"webServerPort":     port,
			"webServerFilePath": filePath,
			"timeOut":           timeOut,
		}
		buf := bytes.NewBuffer(nil)
		f, _ := os.Open("./openfoam_test.log")
		io.Copy(buf, f)
		f.Close()
		s := string(buf.Bytes())

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://192.168.192.200:8000/test.log",
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, s)
				return resp, nil

			},
		)

		Convey("So should collect k metrics", func() {
			metrics := []plugin.Metric{}
			metrics = append(metrics, plugin.Metric{Namespace: plugin.NewNamespace(Vendor, Plugin, "k", "initial"), Config: config})

			collect, err := openFoamCol.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(collect[0].Data, ShouldNotBeNil)
			So(len(collect), ShouldResemble, 1)

		})
		Convey("So should collect Uz metrics", func() {
			metrics := []plugin.Metric{}

			metrics = append(metrics, plugin.Metric{Namespace: plugin.NewNamespace(Vendor, Plugin, "Uz", "final"), Config: config})

			collect, err := openFoamCol.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(collect[0].Data, ShouldNotBeNil)
			So(collect[0].Data, ShouldResemble, 2.13410839e-06)
			So(len(collect), ShouldResemble, 1)

		})
		Convey("So should collect Ux metrics", func() {
			metrics := []plugin.Metric{}
			metrics = append(metrics, plugin.Metric{Namespace: plugin.NewNamespace(Vendor, Plugin, "Ux", "final"), Config: config})

			collect, err := openFoamCol.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(collect[0].Data, ShouldNotBeNil)
			So(collect[0].Data, ShouldResemble, 2.36534655e-07)
			So(len(collect), ShouldResemble, 1)

		})
		Convey("So should return error if value dosn't exist", func() {
			metrics := []plugin.Metric{}
			metrics = append(metrics, plugin.Metric{Namespace: plugin.NewNamespace(Vendor, Plugin, "fUx", "final"), Config: config})
			_, err := openFoamCol.CollectMetrics(metrics)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "Can't find data in OpenFoamLog")

		})
	})

}
