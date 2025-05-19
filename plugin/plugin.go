package plugin

import (
	_ "embed"
)

// //go:embed wasm/plugin.wasm
// var PluginSrc []byte

// //go:embed wasm/language.wasm
//
//go:embed gazelle_go/language.wasm
var LanguageSrc []byte
