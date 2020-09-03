FROM golang:1.14-rc-stretch

WORKDIR $GOPATH/src/blinktag.com/bikesy-wrapper

COPY . .
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure -update
ENV CONFIG=$GOPATH/src/blinktag.com/bikesy-wrapper/config/development.yaml

RUN go build

CMD [ "./bikesy-wrapper" ]


