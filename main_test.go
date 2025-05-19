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

	// Declare variable outside the loop to avoid optimizations
	var i uint64 = 0
	for b.Loop() {
		res, _ := wasmRun(ctx, add, i, 1)
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
	runtime := createRuntime(ctx)
	defer runtime.Close(ctx)
	mod, err := createWasmMod(ctx, runtime)
	if err != nil {
		panic(err)
	}
	iterations := mod.ExportedFunction("inner_loop_iterations")

	inner := uint64(iterationsInner)
	outer := uint64(iterationsOuter)

	for b.Loop() {
		_, err = wasmRun(ctx, iterations, outer, inner)
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
