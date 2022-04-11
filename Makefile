.PHONY : clean

buildTime = `date +%Y-%m-%dT%T%z`
target = main.go
ldflags = -ldflags="-s -w -X main.buildTime=${buildTime}"
gcflags = -gcflags="-trimpath=${PWD}"
output = -o=synod

build:
	CGO_ENABLED=0 go build ${ldflags} ${gcflags} ${output} ${target}

clean:
	rm -rf var/disk/*
	rm -rf var/temp/*
	curl -XDELETE localhost:9200/metas -v