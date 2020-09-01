ls -d */ | grep -v vendor | while read d; do
    ${GOPATH}/bin/golint ./${d}*.go
    go test blinktag.com/bikesy-wrapper/${d} -cover
done
