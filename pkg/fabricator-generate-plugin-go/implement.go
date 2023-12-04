// Code generated by fabricator-generate-plugin-go
//
// Modifications in code regions will be lost during regeneration!

package fabricatorgenerateplugingo

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"code.cestus.io/libs/codegenerator/pkg/templating"

	"code.cestus.io/libs/buildinfo"
	"code.cestus.io/tools/fabricator/pkg/fabricator"
	"code.cestus.io/tools/fabricator/pkg/helpers"
)

// ensure it implements Generator
var _ Generator = (*plugin)(nil)

type plugin struct {
	pluginConfig PluginConfig
	root         string
	goModule     string

	pack templating.Pack
	packprovider templating.PackProvider
}

func (p *plugin) Root() string {
	return p.root
}
func (p *plugin) GoModule() string {
	return p.goModule
}

// region CODE_REGION(GENERATION_CONTEXT)
type GenerationContext struct {
	CodeGenerator       buildinfo.BuildInfo
	GoModule            string
	PluginComponent     PluginComponent
	PinDependencies     PinDependencies
	ReplaceDependencies ReplaceDependencies
	ToolDependencies    ToolDependencies
	// endregion
}

func (p *plugin) generationContexts(ctx context.Context, io fabricator.IOStreams) ([]GenerationContext, error) {
	var contexts []GenerationContext
	for _, component := range p.pluginConfig.Components {
		gencontext := GenerationContext{
			CodeGenerator:       buildinfo.ProvideBuildInfo(),
			GoModule:            p.GoModule(),
			PluginComponent:     component,
			PinDependencies:     DefaultPins,
			ReplaceDependencies: DefaultReplacements,
			ToolDependencies:    DefaultToolDependencies,
		}
		// override defaults if necessary
		for k, o := range component.Spec.PinDependency {
			gencontext.PinDependencies[k] = o
		}
		for k, o := range component.Spec.ReplaceDependency {
			gencontext.ReplaceDependencies[k] = o
		}
		for k, o := range component.Spec.ToolDependency {
			gencontext.ToolDependencies[k] = o
		}
		contexts = append(contexts, gencontext)
	}

	return contexts, nil
}

func (p *plugin) Generate(ctx context.Context, io fabricator.IOStreams, patterns ...string) (err error) {
	genCtxs, err := p.generationContexts(ctx, io)
	if err != nil {
		return fmt.Errorf("failed to generate template contexts for %s: %s", PluginName, err)
	}
	var extGen [][]string
	var genCmds [][]string
	executor := helpers.NewExecutor(p.Root(), io)
	for _, genCtx := range genCtxs {
		templates, err := p.pack.LoadTemplates()
		if err != nil {
			return fmt.Errorf("failed to load template for plugin: %s", err)
		}
		generatedFiles, generationCommands, err := templating.Render(templates, p.Root(), genCtx, patterns...)
		if err != nil {
			return fmt.Errorf("failed to generate template for %s: %s", PluginName, err)
		}

		for _, generatedFile := range generatedFiles {
			generatedFile, _ = filepath.Rel(p.Root(), generatedFile)
			fmt.Fprintf(io.Out, "%s\n", generatedFile)
		}

		for _, v := range genCtx.PinDependencies {
			extGen = append(extGen, []string{"go", "mod", "edit", "--require", fmt.Sprintf("%s@%s", v.Name, v.Version)})
		}
		for _, v := range genCtx.ReplaceDependencies {
			extGen = append(extGen, []string{"go", "mod", "edit", "--replace", fmt.Sprintf("%s=%s", v.Name, v.With)})
		}

		genCmds = append(genCmds, generationCommands...)
	}
	for _, generationCommand := range extGen {
		if err = executor.Run(ctx, generationCommand[0], generationCommand[1:]...); err != nil {
			return fmt.Errorf("failed to run template generation commands for project: %s", err)
		}
	}
	if err = executor.Run(ctx, "go", "mod", "tidy" /*, "-compat=1.17"*/); err != nil {
		fmt.Fprintf(io.Out, "go mod tidy failed: %s\n", err.Error())
	}
	// format files first so we dont run into generation failures
	for _, generationCommand := range genCmds {
		if generationCommand[0] == "goimports" {
			if err = executor.Run(ctx, generationCommand[0], generationCommand[1:]...); err != nil {
				return fmt.Errorf("failed to run template generation commands for project: %s", err)
			}
		}
	}
	for _, generationCommand := range genCmds {
		if err = executor.Run(ctx, generationCommand[0], generationCommand[1:]...); err != nil {
			return fmt.Errorf("failed to run template generation commands for project: %s", err)
		}
	}

	return nil
}

func newPlugin(ctx context.Context, io fabricator.IOStreams, config PluginConfig, root string, packprovider templating.PackProvider) (plugin, error) {

	goModule, err := helpers.GetGoModule(ctx, io, root)
	if err != nil {
		return plugin{}, errors.New("could not deduce gomodule")
	}
	//templating.SetCodeGeneratorName("fabricator")
	plugin := plugin{
		pluginConfig: config,
		root:         root,
		goModule:     goModule.Path,
	}

	pack, err := packprovider.Provide(PluginName, "go")
	if err != nil {
		return plugin, fmt.Errorf("failed to load pack for plugin: %s", err)
	}
	// test that they can be loaded
	_, err = pack.LoadTemplates()
	if err != nil {
		return plugin, fmt.Errorf("failed to load template for plugin: %s", err)
	}

	plugin.pack = pack
	plugin.packprovider = packprovider
	return plugin, err
}
