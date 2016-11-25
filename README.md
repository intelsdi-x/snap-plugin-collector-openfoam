# Snap Collector Plugin - OpenFOAM

[![Build Status](https://api.travis-ci.org/intelsdi-x/snap-plugin-collector-openfoam.svg)](https://travis-ci.org/intelsdi-x/snap-plugin-collector-openfoam)
[![Go Report Card](https://goreportcard.com/badge/intelsdi-x/snap-plugin-collector-openfoam)](https://goreportcard.com/report/intelsdi-x/snap-plugin-collector-openfoam)

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
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

OpenFOAM instance with log file hosted by http server http://www.openfoam.com/ .

#### Download the plugin binary:

You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-openfoam/releases) page.

#### Compile plugin
```
make
```
### Configuration and Usage

This plugin requires these config values to be set (examples are given):
- `webServerIP` (ex. `"192.168.122.89"`)
- `webServerPort` (ex. `8000`)
- `webServerFilePath` (ex. `"run.log"`)

### Documentation

### Examples
Example of running snap OpenFOAM collector and writing data to file.

Ensure [snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `sudo service snap-telemetry start`
* systemd: `sudo systemctl start snap-telemetry`
* command line: `sudo snapteld -l 1 -t 0 &`

Download and load snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-openfoam/latest/linux/x86_64/snap-plugin-collector-openfoam
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ snaptel plugin load snap-plugin-collector-openfoam
$ snaptel plugin load snap-plugin-publisher-file
```

See available metrics for your system.
```
$ snaptel metric list
```

Create a task JSON file:    
```json
{
    "version":1,
    "schedule":{
        "type":"simple",
        "interval":"1s"
    },
    "workflow":{
        "collect":{
            "metrics":{
                "/intel/openfoam/k/initial":{},
                "/intel/openfoam/Ux/final":{},
                "/intel/openfoam/Uy/initial":{}
            },
            "config":{
                "/intel/openfoam":{
                    "webServerIP":"192.168.122.89",
                    "webServerPort":8000,
                    "webServerFilePath":"run.log"
                }
            },
            "process":null,
            "publish":[
                {
                    "plugin_name":"file",
                    "config":{
                        "file":"/tmp/published_openfoam"
                    }
                }
            ]
        }
    }
}
```
Alternatively use provided example manifest:
Change ip address and port of openfoam host in task manifest:
```
vim example/openfoam-file-example.json
```

Create task:
```
$ snaptel task create -t example/openfoam-file-example.json
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
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).
And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
This is Open Source software released under the Apache 2.0 License. Please see the [LICENSE](LICENSE) file for full license details.

* Author: [Marcin Spoczynski](https://github.com/sandlbn/)

This software has been contributed by MIKELANGELO, a Horizon 2020 project co-funded by the European Union. https://www.mikelangelo-project.eu/
