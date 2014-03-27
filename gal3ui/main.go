//Copyright (c) 2012 Tim Shannon
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in
//all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//THE SOFTWARE.

package main

import (
	"bitbucket.org/tshannon/config"
	rest "bitbucket.org/tshannon/gal3upload/gal3rest"
	"gopkg.in/v0/qml"
	"os"
)

var client *rest.Client
var dialog qml.Object

func main() {
	qml.Init(nil)
	engine := qml.NewEngine()

	component, err := engine.LoadFile("main.qml")
	panic(err)

	defer os.Exit(0)

	LoadConfig()

	if client != nil {
		engine.Context().SetVar("client", client)
	}

	win := component.CreateWindow(nil)

	win.Show()
	win.Wait()
}

func errHandle(err error) bool {
	//TODO: show dialog
	if err != nil {
		dialog.Call("show", "An error occurred: "+err.Error())
		return true
	}
	return false
}

func LoadConfig() {
	cfg, err := config.LoadOrCreate("settings.json")
	if err != nil {
		panic(err)
	}

	if cfg.String("url", "") == "url" || cfg.String("apikey", "") == "" {
		//TODO: Show input dialog

	}
}
