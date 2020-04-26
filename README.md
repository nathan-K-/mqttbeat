# Mqttbeat

[![Go Report Card](https://goreportcard.com/badge/github.com/nathan-K-/mqttbeat)](https://goreportcard.com/report/github.com/nathan-K-/mqttbeat)

Welcome to Mqttbeat.

This beat will allow you to put MQTT messages in an elasticsearch instance.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/nathan-K-/mqttbeat`

## Support

Hello, this project is no longer under work. Feel free to fork it, adapt it to your need, and improve it !
Cheers,

## Getting Started with Mqttbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7 (1.8 recommended)

### Init Project
To get running with Mqttbeat and also install the
dependencies, run the following command, with [glide](https://github.com/Masterminds/glide) installed:

```
glide install
pip install virtualenv (for the testsuite)
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Mqttbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/nathan-K-/mqttbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Mqttbeat run the command below. This will generate a binary
in the same directory with the name mqttbeat.

```
make
```


### Run

To run Mqttbeat with debugging output enabled, run:

```
./mqttbeat -c mqttbeat.yml -e -d "*"
```


### Test

To test Mqttbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/mqttbeat.template.json and etc/mqttbeat.asciidoc

```
make update
```


### Cleanup

To clean  Mqttbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Mqttbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/nathan-K-/
cd ${GOPATH}/github.com/nathan-K-/
git clone https://github.com/nathan-K-/mqttbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
