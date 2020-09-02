# blinktag.com/bikesy-wrapper

Wraps request to biksey-api, including routing for different safety profiles and matching to elevation data not provided inside OSRM docker.

## Install

Make sure code is in
$GOPATH/blinktag.com/bikesy-wrapper

```
dep ensure -update
```

## Local Redis Install for development
See bikesy-api for examples of how to populate elevation data.  By default will talk to heroku, but can use local DB if needed.
```
brew install redis
launchctl load ~/Library/LaunchAgents/homebrew.mxcl.redis.plist
```

## Run

```
go build
export CONFIG=./config/development.yaml; ./bikesy-wrapper
```

## Response
TO BE COMPLETED

## Sample Request
Assumes that OSRM services is properly hosted on heroku (see config).
```
curl "http://localhost:8888/route?lng1=-122.424474&lat1=37.766237&lng2=-122.443049&lat2=37.775325"
```

## Tests and Linting
Ensure that $GOPATH/bin/ is in your $PATH variable to execute command.
```
brew install vektra/tap/mockery
go get -u golang.org/x/lint/golint
sh test.sh
```
