package main

import (
	"mb-feedback/internal/app"
)

func main() {
	a := &app.App{}

	a.Init()
	a.Start()
	a.Listen()
	a.Stop()
	a.Exit()
}

/**
 * Your TimeMap object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Set(key,value,timestamp);
 * param_2 := obj.Get(key,timestamp);
 */
