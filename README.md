
## Proposal
To create a system to register Gazelle languages at startup, so that gazelle extensions don't have to be compiled with the binary. Further, we would strongly prefer if the extensions were language agnostic, so that users could write Gazelle extensions in Zig, for instance.

We will define a Gazelle extension that will act as an intermediary between the WebAssembly world and Gazelle proper. We'll call this extension the "host extension". The host extension will take care of instantiating a WebAssembly runtime (in this case, [wazero](https://wazero.io/)), as well as loading into it the modules from the compiled WASM sources of any plugins the user wants to use (we'll call these "guest extensions").

The host extension will register itself as just another Gazelle `language.Language`. However, all the calls to functions of the `Language` interface, such as `Fix` and `GenerateRules` will be delegated to the guest extensions loaded into the WASM runtime. An example of such a call can be found [here](https://github.com/bytebardsoftware/bazel-experiments/blob/eb813298e92722120cac64d5de920ff69f35316f/host/host.go#L89-L109).
## Tradeoffs & Concerns

### The `Language` Interface

Communicating with WebAssembly requires C-style argument passing and memory management. If we have to pass non-trivial data (e.g. a `string`), we have to allocate it ourselves and pass a pointer to it. If we want to read non-trivial data, we have to get a raw memory address from it and read it into a `[]byte` before unmarshalling it.

This means that the interface between the host extension and the guest extensions will not be quite 1:1 with the `language.Language` interface. For instance, we will have to re-define data structures for `config.Config` so that it can be passed via raw pointer and parsed by the guest extension.

I believe this is a minor inconvenience, and one that we would have to deal with anyway when defining a cross-language API.
### Performance

WebAssembly as a runtime is less performant than native languages like Rust, or close-to-native languages like Go.

Here are the results of running some naive benchmarks on different parts of the WASM pipeline ([Source Code](https://github.com/bytebardsoftware/bazel-experiments/blob/main/main_test.go)):

```shell
$   bazel run //:bazel-experiments_test --  -test.count 1 -test.run=^# -test.bench=.  -test.benchtime=10s

goos: darwin
goarch: arm64
cpu: Apple M1 Max
BenchmarkRuntimeCreation-10                        52356            226627 ns/op
BenchmarkModuleLoading-10                            141          81433219 ns/op
BenchmarkSimpleFunctionCallWasm-10              289165588               41.13 ns/op
BenchmarkSimpleFunctionCallGo-10                1000000000              10.00 ns/op
BenchmarkExpensiveFunctionCallWasm-10                 90         132338913 ns/op
BenchmarkExpensiveFunctionCallGo-10                  151          79093680 ns/op
BenchmarkHostLangName-10                        282122169               42.58 ns/op
BenchmarkJustCallingLangName-10                 393749365               30.43 ns/op
BenchmarkReadLangNameFromMemory-10              1000000000              10.00 ns/op
```

From here, we can surmise that:
- On trivial function calls (`BenchmarkSimpleFunctionCallWasm-10`), WASM performs 4-10 times worse than pure Go. The longer we run the benchmarks for, the smaller the gap.
- On expensive function calls without I/0 (`BenchmarkExpensiveFunctionCall`, WASM performs about 2 times worse than Go.
- Reading from WASM memory takes a trivial amount of time.

I believe these benchmarks mean that, while WASM is almost universally slower than pure Go, it's still fast enough to handle Gazelle language extensions at a reasonable speed. 

Note that wazero compiles and optimizes the WASM binary when the module is loaded. This means that WASM might actually be more performant than certain source languages, such as Python and Ruby. There could be a future where users can write Gazelle extensions in Python, and have them perform just as well as one written in Go.

### Performance (cont'd)

A big concern is that, currently, WebAssembly does not have a concurrency model. There are [proposals](https://github.com/webassembly/threads) trying to bring threads to it, but they're not here yet.

I believe the Gazelle `Language` API is well-designed to work around this, since the execution is coordinated by Gazelle-core already, and thus the concurrency already comes "for free". In other words, a WebAssembly guest extension doesn't need concurrency for the same reason that the `proto` extension doesn't need concurrency.
## Alternatives Considered

### An RPC-like server (e.g. Cap'n Proto) with shared-memory IPC 

The obvious alternative is to define the Gazelle Language API in a performant cross-language data interchange format, such as [Cap'n Proto](https://capnproto.org/), and use some sort of shared-memory IPC mechanism to perform procedure calls on the guest extensions.

For the common use-case of writing extensions in a compiled language such as Rust or Zig, this solution would almost certainly be more performant, since WASM will (for our purposes) always be slower than a compiled binary.

However, I believe that the ergonomics of compiling to WASM are better than the ergonomics of generating code from a Cap'n Proto or Protobuf specification. The state of the art of code generation is variable across languages. Furthermore, all plugin authors would have to find libraries to serialize their data to the right format, as well as run RPC servers using whichever IPC protocol we deem adequate.

One option we may want to consider is to define this API in Cap'n Proto anyway, and then have the WebAssembly host extension just be a layer on top of it:
### A pre-made Go plugin API (e.g. go-plugin)

There is a project that offers a WebAssembly-based Go plugin system out of the box: https://github.com/knqyf263/go-plugin. It uses gRPC and Protocol Buffers as API definition and interchange formats, and it includes code generation facilities to perform RPCs and marshal/unmarshal the results.

While this project offers a lot of utilities when it comes to talking to WebAssembly over the wire (e.g. reading from raw pointers), it's only tailored writing guest extensions in Go. I don't think we want to use gRPC as an API definition format for the reasons outlined in [[#An RPC-like server (e.g. Cap'n Proto) with shared-memory IPC]], and the memory management facilities are quite easy to replicate.