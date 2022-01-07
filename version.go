package fabricatorgenerateplugingo

import (
	_ "embed"

	"code.cestus.io/tools/fabricator/pkg/genericclioptions"
)

//go:embed version.yml
var version string

func init() {
	genericclioptions.SetupVersion(GetVersionYaml(), "fabricator-generate-plugin-go")
}

func GetVersionYaml() []byte {
	return []byte(version)
}
