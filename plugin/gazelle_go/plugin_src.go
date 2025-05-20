package main

import (
	"fmt"
	"unsafe"

	"github.com/bazelbuild/bazel-gazelle/language"
	gazellego "github.com/bazelbuild/bazel-gazelle/language/go"
)

//go:wasmexport add
func add(x, y uint32) uint32 {
	return x + y
}

//go:wasmexport sub
func sub(x, y uint32) uint32 {
	return x - y
}

//go:wasmexportg inner_loop_iterations
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
func main() {
	fmt.Println("Hello World")
}

const trivialName = "trivial"

type trivialLang struct {
	delegate language.Language
}

var delegate = gazellego.NewLanguage()
var delegateName string = delegate.Name()

//go:wasmexport Name
func Name() unsafe.Pointer { return unsafe.Pointer(&delegateName) }

// //go:wasmexport CheckFlags
// func (p *trivialLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
// 	return p.delegate.CheckFlags(fs, c)
// }

// //go:wasmexport Configure
// func (p *trivialLang) Configure(c *config.Config, rel string, f *rule.File) {
// 	p.delegate.Configure(c, rel, f)
// }

// // Embeds implements language.Language.
// //
// //go:wasmexport Embeds
// func (p *trivialLang) Embeds(r *rule.Rule, from label.Label) []label.Label {
// 	panic("unimplemented")
// }

// // Fix implements language.Language.
// //
// //go:wasmexport Fix
// func (p *trivialLang) Fix(c *config.Config, f *rule.File) {
// 	panic("unimplemented")
// }

// // GenerateRules implements language.Language.
// //
// //go:wasmexport GenerateRules
// func (p *trivialLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
// 	panic("unimplemented")
// }

// // Imports implements language.Language.
// //
// //go:wasmexport Imports
// func (p *trivialLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
// 	panic("unimplemented")
// }

// // Kinds implements language.Language.
// //
// //go:wasmexport Kinds
// func (p *trivialLang) Kinds() map[string]rule.KindInfo {
// 	panic("unimplemented")
// }

// // KnownDirectives implements language.Language.
// //
// //go:wasmexport KnownDirectives
// func (p *trivialLang) KnownDirectives() []string {
// 	panic("unimplemented")
// }

// // Loads implements language.Language.
// //
// //go:wasmexport Loads
// func (p *trivialLang) Loads() []rule.LoadInfo {
// 	panic("unimplemented")
// }

// // RegisterFlags implements language.Language.
// //
// //go:wasmexport RegisterFlags
// func (p *trivialLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
// 	panic("unimplemented")
// }

// // Resolve implements language.Language.
// //
// //go:wasmexport Resolve
// func (p *trivialLang) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
// 	panic("unimplemented")
// }

// func NewLanguage() language.Language {
// 	return &trivialLang{
// 		delegate: gazellego.NewLanguage(),
// 	}
// }
