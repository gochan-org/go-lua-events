package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/jiyeyuran/go-eventemitter"
	lua "github.com/yuin/gopher-lua"
)

var (
	em = eventemitter.NewEventEmitter()
)

type testPost struct {
	ID      int
	Name    string
	Email   string
	Subject string
	Message string
}

func loadObjects() {
	state.SetGlobal("addListener", state.NewFunction(func(l *lua.LState) int {
		name := l.ToString(1)
		fmt.Println("Adding listener for", name)
		cb := l.ToFunction(2)

		em.On(name, func(interfaces ...interface{}) {
			if len(interfaces) != 1 {
				return
			}
			data := l.NewTable()
			// there might be a simpler/automatic way to do this, but for now at least this'll have to do
			postReflect := reflect.ValueOf(interfaces[0])
			postType := postReflect.Type()
			numFields := postReflect.NumField()
			for i := 0; i < numFields; i++ {
				field := postReflect.Field(i)
				switch field.Kind() {
				case reflect.Int:
					fallthrough
				case reflect.Float32:
					fallthrough
				case reflect.Float64:
					data.RawSetString(postType.Field(i).Name, lua.LNumber(field.Int()))
				case reflect.String:
					data.RawSetString(postType.Field(i).Name, lua.LString(field.String()))
				}
			}

			err := l.CallByParam(lua.P{
				Fn: cb,
			}, data)
			if err != nil {
				panic(err.Error())
			}
		})

		return 0
	}))

	state.SetGlobal("doEvents", state.NewFunction(func(l *lua.LState) int {
		fmt.Println("calling doEvents()")
		em.Emit("post-received", testPost{
			ID:      99,
			Name:    "Poster",
			Email:   "noko",
			Subject: "Test post",
			Message: "Blah blah blah",
		})
		return 0
	}))
}

func loadPlugins(fullpath string, info os.FileInfo, err error) error {
	if info.IsDir() || strings.ToLower(filepath.Ext(fullpath)) != ".lua" {
		return nil
	}

	fmt.Println("Running file", fullpath)
	return state.DoFile(fullpath)
}
