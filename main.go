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
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// main is an example of how to extend a Go application with an addition
// function defined in WebAssembly.
//
// Since addWasm was compiled with TinyGo's `wasi` target, we need to configure
// WASI host imports.
func main() {
	// Parse positional arguments.
	flag.Parse()

	// Choose the context to use for function calls.
	// ctx := context.Background()
	ctx := context.TODO()

	runtime := createRuntime(ctx)
	defer runtime.Close(ctx)

	mod, err := createWasmMod(ctx, runtime)
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}

	fun := mod.ExportedFunction("iterations_squared")
	fmt.Printf("BL: iterations START")
	fun.Call(ctx, 50_000)
	fmt.Printf("BL: iterations DONE")

	// // Read two args to add.
	// op, x, y, err := readOpAndTwoArgs(flag.Arg(0), flag.Arg(1), flag.Arg(2))
	// if err != nil {
	// 	log.Panicf("failed to read arguments: %v", err)
	// }

	// fun := mod.ExportedFunction(op)
	// results, err := wasmRun(ctx, fun, x, y)
	// if err != nil {
	// 	log.Panicf("failed to call add: %v", err)
	// }
	//
	// fmt.Printf("%d %s %d = %d\n", x, op, y, results[0])
}

func createRuntime(ctx context.Context) wazero.Runtime {
	// Create a new WebAssembly Runtime.
	runtime := wazero.NewRuntime(ctx)

	// Instantiate WASI, which implements host functions needed for TinyGo to
	// implement `panic`.
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	return runtime
}

func createWasmMod(ctx context.Context, runtime wazero.Runtime) (api.Module, error) {
	// Instantiate the guest Wasm into the same runtime. It exports the `add`
	// function, implemented in WebAssembly.
	// log.Printf("BL: Plguin src is %s", wasm_plugin.PluginSrc)
	return runtime.InstantiateWithConfig(ctx, wasm_plugin.PluginSrc, wazero.NewModuleConfig().WithStartFunctions("_initialize"))
}

func wasmRun(ctx context.Context, fun api.Function, args ...uint64) ([]uint64, error) {
	return fun.Call(ctx, args...)
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
