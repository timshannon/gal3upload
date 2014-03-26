package main

import (
	"bitbucket.org/tshannon/config"
	"fmt"
	"gopkg.in/v0/qml"
	"os"
)

var cfg *config.Cfg

func main() {
	qml.Init(nil)
	engine := qml.NewEngine()

	component, err := engine.LoadFile("main.qml")
	handleError(err)

	defer os.Exit(0)

	handleError(LoadConfig())

	win := component.CreateWindow(nil)

	win.Show()
	win.Wait()
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func LoadConfig() error {
	var err error
	cfg, err = config.LoadOrCreate("settings.json")
	return err
}
