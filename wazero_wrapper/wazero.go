package wazero_wrapper

import (
	"context"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func CreateRuntime(ctx context.Context) wazero.Runtime {
	// Create a new WebAssembly Runtime.
	runtime := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigCompiler())

	// Instantiate WASI, which implements host functions needed for TinyGo to
	// implement `panic`.
	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	return runtime
}

func CreateWasmMod(ctx context.Context, runtime wazero.Runtime, src []byte) (api.Module, error) {
	module, err := runtime.CompileModule(ctx, src)
	if err != nil {
		return nil, err
	}
	return runtime.InstantiateModule(ctx, module, wazero.NewModuleConfig().WithStartFunctions("_initialize"))
}

func WasmRun(ctx context.Context, fun api.Function, args ...uint64) ([]uint64, error) {
	return fun.Call(ctx, args...)
}
