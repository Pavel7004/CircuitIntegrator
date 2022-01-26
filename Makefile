##
# GraphPlot
#
# @file
# @version 0.1

all: build

build:
	go build ./cmd/plotter.go

clear:
	rm -rf plotter *.png __debug_bin

.PHONY: all clear

# end
