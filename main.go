package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"

	wasm_plugin "github.com/bytebard.software/bazel-experiments/plugin"
	"github.com/bytebard.software/bazel-experiments/wazero_wrapper"
)

// main is an example of how to extend a Go application with an addition
// function defined in WebAssembly.
//
// Since addWasm was compiled with TinyGo's `wasi` target, we need to configure
// WASI host imports.
func main() {
	// Parse positional arguments.
	flag.Parse()

	ctx := context.Background()

	runtime := wazero_wrapper.CreateRuntime(ctx)
	defer runtime.Close(ctx)

	mod, err := wazero_wrapper.CreateWasmMod(ctx, runtime, wasm_plugin.LanguageSrc)
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}

	op, x, y, err := readOpAndTwoArgs(flag.Arg(0), flag.Arg(1), flag.Arg(2))
	if err != nil {
		log.Panicf("failed to read arguments: %v", err)
	}

	fun := mod.ExportedFunction(op)
	results, err := wazero_wrapper.WasmRun(ctx, fun, x, y)
	if err != nil {
		log.Panicf("failed to call %s: %v", op, err)
	}

	fmt.Printf("%d %s %d = %d\n", x, op, y, results[0])
}

func readOpAndTwoArgs(op, xs, ys string) (string, uint64, uint64, error) {
	if op == "" || xs == "" || ys == "" {
		return "", 0, 0, errors.New("must specify an operation and two command line arguments")
	}

	x, err := strconv.ParseUint(xs, 10, 64)
	if err != nil {
		return "", 0, 0, fmt.Errorf("argument X: %v", err)
	}

	y, err := strconv.ParseUint(ys, 10, 64)
	if err != nil {
		return "", 0, 0, fmt.Errorf("argument Y: %v", err)
	}

	return op, x, y, nil
}
