package main

import (
	"github.com/GPlaczek/taskmaster/pkg/api"
	"github.com/GPlaczek/taskmaster/pkg/data/mem"
)

func main() {
	a := api.NewApi(mem.NewData())
	a.RunServer()
}
