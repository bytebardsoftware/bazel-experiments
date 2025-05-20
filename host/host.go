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

type exportedFunctions struct {
	Name api.Function
}

type LanguageHost struct {
	runtime           wazero.Runtime
	module            api.Module
	ctx               context.Context
	exportedFunctions exportedFunctions
}

// CheckFlags implements language.Language.
func (p *LanguageHost) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	panic("unimplemented")
}

// Configure implements language.Language.
func (p *LanguageHost) Configure(c *config.Config, rel string, f *rule.File) {
	panic("unimplemented")
}

// Embeds implements language.Language.
func (p *LanguageHost) Embeds(r *rule.Rule, from label.Label) []label.Label {
	panic("unimplemented")
}

// Fix implements language.Language.
func (p *LanguageHost) Fix(c *config.Config, f *rule.File) {
	panic("unimplemented")
}

// GenerateRules implements language.Language.
func (p *LanguageHost) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	panic("unimplemented")
}

// Imports implements language.Language.
func (p *LanguageHost) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	panic("unimplemented")
}

// Kinds implements language.Language.
func (p *LanguageHost) Kinds() map[string]rule.KindInfo {
	panic("unimplemented")
}

// KnownDirectives implements language.Language.
func (p *LanguageHost) KnownDirectives() []string {
	panic("unimplemented")
}

// Loads implements language.Language.
func (p *LanguageHost) Loads() []rule.LoadInfo {
	panic("unimplemented")
}

// RegisterFlags implements language.Language.
func (p *LanguageHost) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	panic("unimplemented")
}

// Resolve implements language.Language.
func (p *LanguageHost) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	panic("unimplemented")
}

func (p *LanguageHost) Name() string {
	return p.callString("Name")
}

func (p *LanguageHost) callString(funName string) string {
	outs, err := wazero_wrapper.WasmRun(p.ctx, p.exportedFunctions.Name)
	if err != nil {
		log.Panicf("Could not run function %s from wasm plugin: %s", funName, err)
	}
	ptrAndSize := outs[0]
	resPtr := uint32(ptrAndSize >> 32)
	resSize := uint32(ptrAndSize)

	// TODO Call free on the guest to communicate that we no longer need its memory

	bytes, err := ReadMemory(p.module.Memory(), resPtr, resSize)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func (p *LanguageHost) Close() {
	p.runtime.Close(p.ctx)
}

func NewHost() *LanguageHost {
	ctx := context.Background()
	runtime := wazero_wrapper.CreateRuntime(ctx)

	mod, err := wazero_wrapper.CreateWasmMod(ctx, runtime, wasm_plugin.LanguageSrc)
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}

	functions := exportedFunctions{
		Name: mod.ExportedFunction("Name"),
	}
	return &LanguageHost{
		ctx:               ctx,
		runtime:           runtime,
		module:            mod,
		exportedFunctions: functions,
	}
}
