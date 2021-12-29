##
# GraphPlot
#
# @file
# @version 0.1

all: build

build:
	go build ./cmd/plotter.go

clear:
	rm -rf plotter *.png

.PHONY: all clear

# end
