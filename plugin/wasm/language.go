package main

import (
	"flag"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const trivialName = "trivial"

type trivialLang struct{}

// CheckFlags implements language.Language.
func (p *trivialLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	panic("unimplemented")
}

// Configure implements language.Language.
func (p *trivialLang) Configure(c *config.Config, rel string, f *rule.File) {
	panic("unimplemented")
}

// Embeds implements language.Language.
func (p *trivialLang) Embeds(r *rule.Rule, from label.Label) []label.Label {
	panic("unimplemented")
}

// Fix implements language.Language.
func (p *trivialLang) Fix(c *config.Config, f *rule.File) {
	panic("unimplemented")
}

// GenerateRules implements language.Language.
func (p *trivialLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	panic("unimplemented")
}

// Imports implements language.Language.
func (p *trivialLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	panic("unimplemented")
}

// Kinds implements language.Language.
func (p *trivialLang) Kinds() map[string]rule.KindInfo {
	panic("unimplemented")
}

// KnownDirectives implements language.Language.
func (p *trivialLang) KnownDirectives() []string {
	panic("unimplemented")
}

// Loads implements language.Language.
func (p *trivialLang) Loads() []rule.LoadInfo {
	panic("unimplemented")
}

// RegisterFlags implements language.Language.
func (p *trivialLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	panic("unimplemented")
}

// Resolve implements language.Language.
func (p *trivialLang) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	panic("unimplemented")
}

func (*trivialLang) Name() string { return trivialName }

func NewLanguage() language.Language {
	return &trivialLang{}
}
