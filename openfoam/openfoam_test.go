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
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"
	"github.com/jarcoal/httpmock"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOpenFoamPlugin(t *testing.T) {
	Convey("Meta should return metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, Name)
		So(meta.Version, ShouldResemble, Version)
		So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
	})

	Convey("Create OpenFoam Collector", t, func() {
		openFoamCol := NewOpenFoamCollector()
		Convey("So psCol should not be nil", func() {
			So(openFoamCol, ShouldNotBeNil)
		})
		Convey("So psCol should be of OpenFoam type", func() {
			So(openFoamCol, ShouldHaveSameTypeAs, &OpenFoam{})
		})
		Convey("openFoamCol.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := openFoamCol.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
		})
	})
	Convey("Join namespace ", t, func() {
		namespace1 := []string{"intel", "openfoam", "one"}
		namespace2 := []string{}
		Convey("So namespace should equal intel/openfoam/one", func() {
			So("/intel/openfoam/one", ShouldResemble, joinNamespace(namespace1))
		})
		Convey("So namespace should equal slash", func() {
			So("/", ShouldResemble, joinNamespace(namespace2))
		})

	})
	Convey("Get URI ", t, func() {
		Convey("So should return 10.1.0.1:8000", func() {
			webServerIP := "10.1.0.1"
			webServerPort := 8000
			uri := openFoamURL(webServerIP, webServerPort)
			So("http://10.1.0.1:8000", ShouldResemble, uri)
		})
	})
	Convey("Get Metrics Types", t, func() {
		openFoamCol := NewOpenFoamCollector()
		cfgNode := cdata.NewNode()
		var cfg = plugin.ConfigType{
			ConfigDataNode: cfgNode,
		}
		Convey("So should return 12 types of metrics", func() {
			metrics, err := openFoamCol.GetMetricTypes(cfg)
			So(12, ShouldResemble, len(metrics))
			So(err, ShouldBeNil)
		})
		Convey("So should check namespace", func() {
			metrics, err := openFoamCol.GetMetricTypes(cfg)
			waitNamespace := joinNamespace(metrics[0].Namespace().Strings())
			wait := regexp.MustCompile(`^/intel/openfoam/k/initial`)
			So(true, ShouldEqual, wait.MatchString(waitNamespace))
			So(err, ShouldBeNil)

		})

	})
	Convey("Collect Metrics", t, func() {
		openFoamCol := NewOpenFoamCollector()
		cfgNode := cdata.NewNode()
		cfgNode.AddItem("webServerIP", ctypes.ConfigValueStr{Value: "192.168.192.200"})
		cfgNode.AddItem("webServerPort", ctypes.ConfigValueInt{Value: 8000})
		cfgNode.AddItem("webServerFilePath", ctypes.ConfigValueStr{Value: "test.log"})

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
			metrics := []plugin.MetricType{{
				Namespace_: core.NewNamespace("intel", "openfoam", "k", "initial"),
				Config_:    cfgNode,
			}}
			collect, err := openFoamCol.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(collect[0].Data_, ShouldNotBeNil)
			So(len(collect), ShouldResemble, 1)

		})
		Convey("So should collect Uz metrics", func() {
			metrics := []plugin.MetricType{{
				Namespace_: core.NewNamespace("intel", "openfoam", "Uz", "final"),
				Config_:    cfgNode,
			}}
			collect, err := openFoamCol.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(collect[0].Data_, ShouldNotBeNil)
			So(collect[0].Data_, ShouldResemble, 2.13410839e-06)
			So(len(collect), ShouldResemble, 1)

		})
		Convey("So should collect Ux metrics", func() {
			metrics := []plugin.MetricType{{
				Namespace_: core.NewNamespace("intel", "openfoam", "Ux", "final"),
				Config_:    cfgNode,
			}}
			collect, err := openFoamCol.CollectMetrics(metrics)
			So(err, ShouldBeNil)
			So(collect[0].Data_, ShouldNotBeNil)
			So(collect[0].Data_, ShouldResemble, 2.36534655e-07)
			So(len(collect), ShouldResemble, 1)

		})
		Convey("So should return error if value dosn't exist", func() {
			metrics := []plugin.MetricType{{
				Namespace_: core.NewNamespace("intel", "openfoam", "fUx", "final"),
				Config_:    cfgNode,
			}}
			_, err := openFoamCol.CollectMetrics(metrics)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "Can't find data in OpenFoamLog")

		})
	})

}
