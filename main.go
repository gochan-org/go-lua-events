package main

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

var (
	state *lua.LState
)

func main() {
	var err error
	state = lua.NewState(lua.Options{
		SkipOpenLibs: false,
	})
	defer state.Close()

	loadObjects()
	if err = filepath.Walk("./plugins/", loadPlugins); err != nil {
		panic(err.Error())
	}
}
