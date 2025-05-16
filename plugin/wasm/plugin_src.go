package main

// import "time"

//go:wasmexport add
func add(x, y uint32) uint32 {
	return x + y
}

//go:wasmexport sub
func sub(x, y uint32) uint32 {
	return x - y
}

//go:wasmexport inner_loop_iterations
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

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}
