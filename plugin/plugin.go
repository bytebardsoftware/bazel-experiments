package plugin

import (
	_ "embed"
)

//go:embed wasm/plugin.wasm
var PluginSrc []byte