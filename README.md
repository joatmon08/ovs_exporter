# Open vSwitch Prometheus Exporter

A Prometheus exporter for Open vSwitch.

## Getting Started

### Prerequisites

* Golang
* Open vSwitch 2.5.0
* [OVSDB Schema 7.12.1](https://tools.ietf.org/html/rfc7047)
* [socketplane/libovsdb](https://github.com/socketplane/libovsdb)

### Installing

* Clone the repository.
* Run the exporter: `go run ovs_exporter.go`.
* To Do: Vagrantfile with Open vSwitch & exporter.

## Running the tests

* Unit tests can be run with `go test ./...`.
* To Do: Integration tests

## Deployment

* To deploy, build the binary: `go build ovs_exporter.go`.
* To run, you can specify many options, including:
```
  -listen-port string
        Address to listen on for web interface and telemetry. (default ":9107")
  -metrics-path string
        Path under which to expose metrics. (default "/metrics")
  -uri string
        URI to connect to Open vSwitch (default "/var/run/openvswitch/db.sock")
``` 

## Authors

* **Rosemary Wang** - *Initial work* - [joatmon08](https://github.com/joatmon08)

## License

This project is licensed under the MIT License - see the 
[LICENSE.md](LICENSE.md) file for details
