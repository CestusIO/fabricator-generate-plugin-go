// Generated with the code-generator
//
// Modifications in code regions will be lost during regeneration!

package fabricatorgenerateplugingo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"code.cestus.io/libs/codegenerator/pkg/templating"
	"code.cestus.io/tools/fabricator-generate-plugin-go/pkg/fabricator-generate-plugin-go/templates"
	"code.cestus.io/tools/fabricator/pkg/fabricator"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

//PluginName is the name of the plugin
const PluginName string = "fabricator-generate-plugin-go"

// region CODE_REGION(OPTIONS)
type options struct {
	fabricator.RootOptions
	fabricator.IOStreams
	// endregion
	Add bool // SampleFlag
}

// newOptions returns initialized options
func newOptions(ioStreams fabricator.IOStreams, flagset *pflag.FlagSet, flagparser fabricator.FlagParser) *options {
	o := options{
		IOStreams: ioStreams,
	}
	o.RootOptions.FlagParser = flagparser
	o.RootOptions.RegisterOptions(flagset)
	return &o
}

// region CODE_REGION(OPTION_COMPLETE)
func (o *options) complete(cmd *cobra.Command) error {
	err := o.FlagParser(cmd)
	if err != nil {
		return err
	}
	// endregion
	return nil
}

func (o *options) currentPath() (string, error) {
	path := o.RootOptions.FabricatorFile
	if !filepath.IsAbs(path) {
		path = filepath.Join(o.RootOptions.RootDirectory, path)
	}
	return filepath.Abs(path)
}

// region CODE_REGION(OPTION_RUN)
// Run executes command
func (o *options) run(ctx context.Context) error {
	path, err := o.currentPath()
	if err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	config, err := LoadPluginConfig(f)
	if err != nil {
		return err
	}
	// endregion
	fmt.Fprintf(o.Out, "Loading from %s\n", path)
	templating.CodeGeneratorName = PluginName
	packProvider := templating.NewPackProvider()
	packProvider.RegisterProvider(templating.NewEmbededPackProvider(templates.GetTemplates()))
	plugin, err := newPlugin(ctx, o.IOStreams, config, o.RootOptions.RootDirectory, packProvider)

	if err != nil {
		return err
	}

	err = plugin.Generate(ctx, o.IOStreams)
	if err != nil {
		return err
	}

	return nil
}

// region CODE_REGION(CREATEPLUGIN)
func NewFabricatorGeneratePluginGo(ioStreams fabricator.IOStreams, flagparser fabricator.FlagParser) *cobra.Command {
	//endregion
	cmd := &cobra.Command{
		Use:     PluginName,
		Short:   "ShortDescription",
		Long:    "LongDescription",
		Example: "",
	}
	o := newOptions(ioStreams, cmd.Flags(), flagparser)
	cmd.Flags().BoolVar(&o.Add, "add", o.Add, "Add is a sampleflag")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := o.complete(cmd); err != nil {
			return err
		}
		return o.run(cmd.Context())
	}
	return cmd
}

// Generator is an internal interface for a generator
type Generator interface {
	Root() string
	Generate(ctx context.Context, io fabricator.IOStreams, patterns ...string) (err error)
}
