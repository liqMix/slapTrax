package main

import (
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/kettek/gobl"
)

func main() {
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	runArgs := append([]interface{}{}, "./slapTrax"+exe)
	var wasmSrc string

	Task("build").
		Exec("go", "build", "-o", "slapTrax"+exe, "./cmd/game")
	Task("run").
		Exec(runArgs...)
	Task("watch").
		Watch("cmd/game/*", "internal/*", "internal/render/*").
		Signaler(SigQuit).
		Run("build").
		Run("run")
	Task("build-service").
		Env("GOOS=linux", "GOARCH=amd64").
		Exec("go", "build", "-o", "slapboard", "./cmd/service")
	Task("build-web").
		Env("GOOS=js", "GOARCH=wasm").
		Exec("go", "build", "-o", "web/slapTrax.wasm", "./cmd/game").
		Exec("go", "env", "GOROOT").
		Result(func(i interface{}) {
			goRoot := strings.TrimSpace(i.(string))
			wasmSrc = filepath.Join(goRoot, "misc/wasm/wasm_exec.js")
		}).
		Exec("cp", &wasmSrc, "web/")
	Go()
}
