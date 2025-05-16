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

// func BenchmarkSimpleFunctionCallSubtest(b *testing.B) {
// 	ctx := context.Background()
// 	runtime := createRuntime(ctx)
// 	defer runtime.Close(ctx)
// 	mod, err := createWasmMod(ctx, runtime)
// 	if err != nil {
// 		panic(err)
// 	}
// 	add := mod.ExportedFunction("add")

//		b.Run("adds", func(b *testing.B) {
//			for i := 0; i < b.N; i++ {
//				_, err := wasmRun(ctx, add, 1000, 1000)
//				if err != nil {
//					panic(err)
//				}
//			}
//		})
//	}
func BenchmarkSimpleFunctionCall(b *testing.B) {
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

func BenchmarkExpensiveFunctionCall(b *testing.B) {
	ctx := context.Background()
	runtime := createRuntime(ctx)
	defer runtime.Close(ctx)
	mod, err := createWasmMod(ctx, runtime)
	if err != nil {
		panic(err)
	}
	add := mod.ExportedFunction("expensiveAdd")

	for b.Loop() {
		_, err = wasmRun(ctx, add, 1000, 1000)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkSimpleAdd(b *testing.B) {
	for b.Loop() {
		_ = 1000 + 1000
	}
}