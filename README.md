# Snap Collector Plugin - OpenFOAM

[![Build Status](https://api.travis-ci.org/intelsdi-x/snap-plugin-collector-openfoam.svg)](https://travis-ci.org/intelsdi-x/snap-plugin-collector-openfoam )
[![Go Report Card](http://goreportcard.com/badge/intelsdi-x/snap-plugin-collector-openfoam)](http://goreportcard.com/report/intelsdi-x/snap-plugin-collector-openfoam)
 
 1. [Getting Started](#getting-started)
   * [System Requirements](#system-requirements)
   * [Installation](#installation)
   * [Configuration and Usage](configuration-and-usage)
 2. [Documentation](#documentation)
   * [Collected Metrics](#collected-metrics)
   * [Examples](#examples)
   * [Roadmap](#roadmap)
 3. [Community Support](#community-support)
 4. [Contributing](#contributing)
 5. [License](#license)
 6. [Acknowledgements](#acknowledgements)
 
 ## Getting Started
 
 ### System Requirements
 
 OpenFOAM instance with log file hosted by http server http://www.openfoam.com/
 
 #### Compile plugin
 ```
 make
 ```
 
 ### Documentation
 
 ### Examples
 Example running OpenFoam, passthru processor, and writing data to a file.
 
 In one terminal window, open the snap daemon :
 ```
 $ snapd -t 0 -l 1
 ```
 
 In another terminal window:
 Load OpenFoam plugin
 ```
 $ snapctl plugin load $SNAP_OPENFOAM_PLUGIN_DIR/build/rootfs/snap-plugin-collector-openfoam
 ```
 See available metrics for your system
 ```
 $ snapctl metric list
 ```
 
 Create a task JSON file:    
 ```json
 {
     "version": 1,
     "schedule": {
         "type": "simple",
         "interval": "1s"
     },
     "workflow": {
         "collect": {
             "metrics": {
                 "/intel/openfoam/k/initial": {},
                 "/intel/openfoam/Ux/final": {},
                 "/intel/openfoam/Uy/initial": {}
 
             },
             "config": {
                 "/intel/openfoam": {
                     "webServerIP": "192.168.122.89",
                     "webServerPort": 8000,
                     "webServerFilePath": "run.log"
                 }
             },
             "process": [
                 {
                     "plugin_name": "passthru",
                     "process": null,
                     "publish": [
                         {                         
                             "plugin_name": "file",
                             "config": {
                                 "file": "/tmp/published_openfoam"
                             }
                         }
                     ],
                     "config": null
                 }
             ],
             "publish": null
         }
     }
 }
 ```
 
 Load passthru plugin for processing:
 ```
 $ snapctl plugin load build/rootfs/plugin/snap-processor-passthru
 Plugin loaded
 Name: passthru
 Version: 1
 Type: processor
 Signed: false
 Loaded Time: Fri, 20 Nov 2015 11:44:03 PST
 ```
 
 Load file plugin for publishing:
 ```
 $ snapctl plugin load build/rootfs/plugin/snap-publisher-file
 Plugin loaded
 Name: file
 Version: 3
 Type: publisher
 Signed: false
 Loaded Time: Fri, 20 Nov 2015 11:41:39 PST
 ```
 
 Change ip address and port of openfoam host in task manifest:
 ```
 vim $SNAP_OPENFOAM_PLUGIN_DIR/example/openfoam-file-example.json
 ```
 
 Create task:
 ```
 $ snapctl task create -t $SNAP_OPENFOAM_PLUGIN_DIR/example/openfoam-file-example.json
 Using task manifest to create task
 Task created
 ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
 Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
 State: Running
 ```
 
 
 ### Collected Metrics
 This plugin has the ability to gather the following metrics:
 
 Namespace | Data Type
 ----------|-----------
 /intel/openfoam/k/initial | float64
 /intel/openfoam/k/final | float64 
 /intel/openfoam/p/initial | float64 
 /intel/openfoam/p/final | float64 
 /intel/openfoam/Ux/initial | float64 
 /intel/openfoam/Ux/final | float64 
 /intel/openfoam/Uy/initial | float64 
 /intel/openfoam/Uy/final | float64 
 /intel/openfoam/Uz/initial | float64 
 /intel/openfoam/Uz/final | float64 
 /intel/openfoam/omega/initial | float64 
 /intel/openfoam/omega/final | float64 
 
 ### Roadmap
 As we launch this plugin, we do not have any outstanding requirements for the next release. If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-openfoam/issues).
 
 If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-collector-openfoam/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-collector-openfoam/pulls).
 
 ## Community Support
 This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)
 
 ## Contributing
 We love contributions!
 
 There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).
 
 ## License
 [snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).
 
 ## Acknowledgements
 This is Open Source software released under the Apache 2.0 License. Please see the [LICENSE](LICENSE) file for full license details.
 
 * Author: [Marcin Spoczynski](https://github.com/sandlbn/)
 
 ## Thank You
 And **thank you!** Your contribution, through code and participation, is incredibly important to us.
