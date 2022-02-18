##
# GraphPlot
#
# @file
# @version 0.1

all: build

build:
	go build ./cmd/plotter.go

clear:
	rm -rf plotter results/*.png */__debug_bin __debug_bin results

.PHONY: all clear

# end
