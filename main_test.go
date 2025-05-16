package main

import (
	"context"
	"testing"
)

func BenchmarkRuntimeCreation(b *testing.B) {
	ctx := context.Background()
	for b.Loop() {
		createRuntime(ctx)
	}
}

func BenchmarkModuleLoading(b *testing.B) {
	ctx := context.Background()
	runtime := createRuntime(ctx)
	defer runtime.Close(ctx)

	for b.Loop() {
		createWasmMod(ctx, runtime)
	}
}

func BenchmarkSimpleFunctionCallWasm(b *testing.B) {
	ctx := context.Background()
	runtime := createRuntime(ctx)
	defer runtime.Close(ctx)
	mod, err := createWasmMod(ctx, runtime)
	if err != nil {
		panic(err)
	}
	add := mod.ExportedFunction("add")

	for b.Loop() {
		_, err = wasmRun(ctx, add, 1000, 1000)
		if err != nil {
			panic(err)
		}
	}
}

var expensiveIterations uint32 = 1_000

func BenchmarkExpensiveFunctionCallWasm(b *testing.B) {
	ctx := context.Background()
	runtime := createRuntime(ctx)
	defer runtime.Close(ctx)
	mod, err := createWasmMod(ctx, runtime)
	if err != nil {
		panic(err)
	}
	wait := mod.ExportedFunction("iterations_squared")

	for b.Loop() {
		_, err = wasmRun(ctx, wait, uint64(expensiveIterations))
		if err != nil {
			panic(err)
		}
	}
}

func expensive(iterations uint32) uint32 {
	var res uint32 = 0
	var i uint32 = 0
	for i < iterations {
		var u uint32 = 0
		for u < iterations {
			res = (u % 1000) + (i % 1000)
			u += 1
		}
		i += 1
	}
	return res
}

func BenchmarkExpensiveFunctionCallGo(b *testing.B) {
	for b.Loop() {
		expensive(expensiveIterations)
	}
}
