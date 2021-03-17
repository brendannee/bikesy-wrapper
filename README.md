# blinktag.com/bikesy-wrapper

Wraps request to biksey-api, including routing for different safety profiles and matching to elevation data not provided inside OSRM docker.

## Install

Make sure code is in
$GOPATH/src/blinktag.com/bikesy-wrapper

## Local Redis Install for development
See bikesy-api for examples of how to populate elevation data.  By default will talk to heroku, but can use local DB if needed.
```
brew install redis
launchctl load ~/Library/LaunchAgents/homebrew.mxcl.redis.plist
```

## Configure
Create a .env file with the following:
```
REDIS_URL=redis://127.0.0.1:6379
CONFIG=./config/config.yaml
PORT=8888
```

## Run

```
go build
export CONFIG=./config/config.yaml; ./bikesy-wrapper
```

## Response
TO BE COMPLETED

## Sample Request
Assumes that OSRM services is properly hosted on heroku (see config).  Required params: lat1/lng1/lat2/lat2,hills,safety.  Hills and safety can be low, med, or high.
```
curl "http://localhost:8888/route?lng1=-122.424474&lat1=37.766237&lng2=-122.443049&lat2=37.775325&hills=low&safety=low"
```

## Tests and Linting
Ensure that $GOPATH/bin/ is in your $PATH variable to execute command.
```
brew install vektra/tap/mockery
go get -u golang.org/x/lint/golint
sh test.sh
```

## Heroku
```
heroku git:remote -a bikesy-wrapper
heroku container:login
heroku container:push web
heroku container:release web
```

## Digital Ocean

If needed, ![create your registry in digital ocean](https://www.digitalocean.com/docs/container-registry/quickstart/registry).

Install doctl and init repo (if needed).
```
brew install doctl
doctl auth init --context bikesy
doctl auth switch --context bikesy
```

Add new registry (one time)
```
doctl registry create bikesy
doctl registry login
```

Deploy.
```
env GOOS=linux GOARCH=amd64 docker build -t bikesy-wrapper .
docker tag bikesy-wrapper registry.digitalocean.com/bikesy/bikesy-wrapper
docker push registry.digitalocean.com/bikesy/bikesy-wrapper
```

