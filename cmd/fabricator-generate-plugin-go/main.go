// Generated with the code-generator
//
// Modifications in code regions will be lost during regeneration!

package main

import (
	"context"
	"os"

	_ "code.cestus.io/tools/fabricator-generate-plugin-go"
	fabricatorgenerateplugingo "code.cestus.io/tools/fabricator-generate-plugin-go/pkg/fabricator-generate-plugin-go"
	"code.cestus.io/tools/fabricator/pkg/cmd/version"
	"code.cestus.io/tools/fabricator/pkg/fabricator"
	"code.cestus.io/tools/fabricator/pkg/helpers"
)

func main() {
	// region CODE_REGION(Main)
	ctx := context.Background()
	io := fabricator.NewStdIOStreams()
	ctx, cancel := helpers.WithCancelOnSignal(ctx, io, fabricator.TerminationSignals...)
	defer cancel()
	root := fabricatorgenerateplugingo.NewFabricatorGeneratePluginGo(io, helpers.DefaultFlagParser)
	root.AddCommand(version.NewCmdVersion(io))
	if err := root.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
	// endregion
}
