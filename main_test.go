package main

import (
	"context"
	"testing"

	"github.com/bytebard.software/bazel-experiments/host"
	wasm_plugin "github.com/bytebard.software/bazel-experiments/plugin"
	"github.com/bytebard.software/bazel-experiments/wazero_wrapper"
)

func BenchmarkRuntimeCreation(b *testing.B) {
	ctx := context.Background()
	for b.Loop() {
		wazero_wrapper.CreateRuntime(ctx)
	}
}

func BenchmarkModuleLoading(b *testing.B) {
	ctx := context.Background()
	runtime := wazero_wrapper.CreateRuntime(ctx)
	defer runtime.Close(ctx)

	for b.Loop() {
		wazero_wrapper.CreateWasmMod(ctx, runtime, wasm_plugin.LanguageSrc)
	}
}

func BenchmarkSimpleFunctionCallWasm(b *testing.B) {
	ctx := context.Background()
	runtime := wazero_wrapper.CreateRuntime(ctx)
	defer runtime.Close(ctx)
	mod, err := wazero_wrapper.CreateWasmMod(ctx, runtime, wasm_plugin.LanguageSrc)
	if err != nil {
		panic(err)
	}
	add := mod.ExportedFunction("add")

	// Declare variable outside the loop to avoid optimizations
	var i uint64 = 0
	for b.Loop() {
		res, _ := wazero_wrapper.WasmRun(ctx, add, i, 1)
		// Ignore error
		i = res[0]
	}
}

func add(left uint64, right uint64) uint64 {
	var res = left + right
	return res
}

func BenchmarkSimpleFunctionCallGo(b *testing.B) {
	// Declare variable outside the loop to avoid optimizations
	var i uint64 = 0
	for b.Loop() {
		i = add(i, 1)
	}
}

var iterationsInner uint32 = 10_000
var iterationsOuter uint32 = 10_000

func BenchmarkExpensiveFunctionCallWasm(b *testing.B) {
	ctx := context.Background()
	runtime := wazero_wrapper.CreateRuntime(ctx)
	defer runtime.Close(ctx)
	mod, err := wazero_wrapper.CreateWasmMod(ctx, runtime, wasm_plugin.LanguageSrc)
	if err != nil {
		panic(err)
	}
	iterations := mod.ExportedFunction("inner_loop_iterations")

	inner := uint64(iterationsInner)
	outer := uint64(iterationsOuter)

	for b.Loop() {
		_, err = wazero_wrapper.WasmRun(ctx, iterations, outer, inner)
		if err != nil {
			panic(err)
		}
	}
}

func inner_loop_iterations(outer uint32, inner uint32) uint32 {
	var res uint32 = 0
	var i uint32 = 0
	for i < outer {
		var u uint32 = 0
		for u < inner {
			res = (u % 1000) + (i % 1000)
			u += 1
		}
		i += 1
	}
	return res
}

func BenchmarkExpensiveFunctionCallGo(b *testing.B) {
	for b.Loop() {
		inner_loop_iterations(iterationsOuter, iterationsInner)
	}
}

func TestHostLang(t *testing.T) {
	langHost := host.NewHost()
	defer langHost.Close()

	name := langHost.Name()
	expected := "go"
	if name != expected {
		t.Errorf("Language name was %s, expected %s", name, expected)
	}
}
