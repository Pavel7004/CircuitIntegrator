package main

import "github.com/Pavel7004/GraphPlot/pkg/adapter/http"

func main() {
	s := http.New()

	if err := s.Run(); err != nil {
		panic(err)
	}
}
