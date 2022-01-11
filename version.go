package fabricatorgenerateplugingo

import (
	_ "embed"

	"code.cestus.io/libs/buildinfo"
)

//go:embed version.yml
var version string

func init() {
	buildinfo.GenerateVersionFromVersionYaml(GetVersionYaml(), "fabricator-generate-plugin-go")
}

func GetVersionYaml() []byte {
	return []byte(version)
}
