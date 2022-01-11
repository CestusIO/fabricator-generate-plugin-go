// Generated with the code-generator
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
}

func (p *plugin) Root() string {
	return p.root
}
func (p *plugin) GoModule() string {
	return p.goModule
}

func (p *plugin) generationContexts(ctx context.Context, io fabricator.IOStreams) ([]interface{}, error) {
	var contexts []interface{}
	for _, component := range p.pluginConfig.Components {
		contexts = append(contexts, &struct {
			CodeGenerator   buildinfo.BuildInfo
			GoModule        string
			PluginComponent PluginComponent
		}{
			CodeGenerator:   buildinfo.ProvideBuildInfo(),
			GoModule:        p.GoModule(),
			PluginComponent: component,
		})
	}

	return contexts, nil
}

func (p *plugin) Generate(ctx context.Context, io fabricator.IOStreams, patterns ...string) (err error) {
	genCtxs, err := p.generationContexts(ctx, io)
	if err != nil {
		return fmt.Errorf("failed to generate template contexts for %s: %s", PluginName, err)
	}

	for _, genCtx := range genCtxs {
		templates, err := p.pack.LoadTemplates()
		if err != nil {
			return fmt.Errorf("failed to load template for plugin: %s", err)
		}
		generatedFiles, generationCommands, err := templating.Render(templates, p.Root(), genCtx, patterns...)
		if err != nil {
			return fmt.Errorf("failed to generate template for %s: %s", PluginName, err)
		}

		executor := helpers.NewExecutor(p.Root(), io)
		for _, generatedFile := range generatedFiles {

			generatedFile, _ = filepath.Rel(p.Root(), generatedFile)
			fmt.Fprintf(io.Out, "%s\n", generatedFile)
		}

		for _, generationCommand := range generationCommands {
			if err = executor.Run(ctx, generationCommand[0], generationCommand[1:]...); err != nil {
				return fmt.Errorf("failed to run template generation commands for project: %s", err)
			}
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

	return plugin, err
}
