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

//go:wasmexport iterations_squared
func iterations_squared(iterations uint32) uint32 {
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

// main is required for the `wasi` target, even if it isn't used.
// See https://wazero.io/languages/tinygo/#why-do-i-have-to-define-main
func main() {}
