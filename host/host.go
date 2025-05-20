package host

import (
	"context"
	"flag"
	"log"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"

	wasm_plugin "github.com/bytebard.software/bazel-experiments/plugin"
	"github.com/bytebard.software/bazel-experiments/wazero_wrapper"
)

const hostName = "host"

type hostLang struct {
	runtime wazero.Runtime
	module  api.Module
	ctx     context.Context
}

// CheckFlags implements language.Language.
func (p *hostLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	panic("unimplemented")
}

// Configure implements language.Language.
func (p *hostLang) Configure(c *config.Config, rel string, f *rule.File) {
	panic("unimplemented")
}

// Embeds implements language.Language.
func (p *hostLang) Embeds(r *rule.Rule, from label.Label) []label.Label {
	panic("unimplemented")
}

// Fix implements language.Language.
func (p *hostLang) Fix(c *config.Config, f *rule.File) {
	panic("unimplemented")
}

// GenerateRules implements language.Language.
func (p *hostLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	panic("unimplemented")
}

// Imports implements language.Language.
func (p *hostLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	panic("unimplemented")
}

// Kinds implements language.Language.
func (p *hostLang) Kinds() map[string]rule.KindInfo {
	panic("unimplemented")
}

// KnownDirectives implements language.Language.
func (p *hostLang) KnownDirectives() []string {
	panic("unimplemented")
}

// Loads implements language.Language.
func (p *hostLang) Loads() []rule.LoadInfo {
	panic("unimplemented")
}

// RegisterFlags implements language.Language.
func (p *hostLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	panic("unimplemented")
}

// Resolve implements language.Language.
func (p *hostLang) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	panic("unimplemented")
}

func (p *hostLang) Name() string {
	return p.callString("Name")
}

func (p *hostLang) callString(funName string) string {
	fun := p.module.ExportedFunction(funName)
	outs, err := wazero_wrapper.WasmRun(p.ctx, fun)
	if err != nil {
		log.Panicf("Could not run function %s from wasm plugin: %w", funName, err)
	}
	return string(outs[0])
}

func NewLanguage() language.Language {
	ctx := context.Background()
	runtime := wazero_wrapper.CreateRuntime(ctx)
	defer runtime.Close(ctx)

	mod, err := wazero_wrapper.CreateWasmMod(ctx, runtime, wasm_plugin.LanguageSrc)
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}
	return &hostLang{
		ctx:     ctx,
		runtime: runtime,
		module:  mod,
	}
}
