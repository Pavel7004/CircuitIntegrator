##
# GraphPlot
#
# @file
# @version 0.1

all: graph

graph:
	go build -o $@ main.go

jaeger:
	docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest

lint:
	golangci-lint run ./...

clear:
	rm -rf graph */__debug_bin __debug_bin results

.PHONY: all clear jaeger lint graph

# end
