# blinktag.com/bikesy-wrapper

Wraps request to biksey-api, including routing for different safety profiles and matching to elevation data not provided inside OSRM docker.

## Install

Make sure code is in
$GOPATH/blinktag.com/bikesy-wrapper

```
dep ensure -update
```

## Run

```
go build
export CONFIG=./config/development.yaml; ./bikesy-wrapper
```

## Response
TO BE COMPLETED

## Sample Request
TO BE COMPLETED

## Tests and Linting
TO BE COMPLETED 
Install golint (https://github.com/golang/lint) and mockery (https://github.com/vektra/mockery).
Ensure that $GOPATH/bin/ is in your $PATH variable to execute command.

