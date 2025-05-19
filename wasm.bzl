def compile_to_wasm(name, srcs, out, **kwargs):
    """
    Compile one source file to wasm using tinygo.
    It requires tinygo to be installed: https://tinygo.org/getting-started/install/,
    and of course updating the absolute path.
    """
    native.genrule(
        name = name,
        srcs = srcs,
        outs = [out],
        cmd = "HOME='/Users/blorente' /opt/homebrew/bin/tinygo build -buildmode=c-shared -target=wasip1 -o $@ $(SRCS)",
        **kwargs
    )