// Generated with the code-generator
//
// Modifications in code regions will be lost during regeneration!

package fabricatorgenerateplugingo_test

import (
	"os"

	fabricatorgenerateplugingo "code.cestus.io/tools/fabricator-generate-plugin-go/pkg/fabricator-generate-plugin-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	It("Deserializes a fabricator config", func() {
		file, err := os.Open("./testdata/deserialize.yml")
		Expect(err).ToNot(HaveOccurred())
		config, err := fabricatorgenerateplugingo.LoadPluginConfig(file)
		Expect(err).ToNot(HaveOccurred())
		expected := fabricatorgenerateplugingo.PluginConfig{
			ApiVersion: "fabricator.cestus.io/v1alpha1",
			Kind:       "Config",
			Components: []fabricatorgenerateplugingo.PluginComponent{
				{
					Name:      "fabricator-generate-plugin-go",
					Generator: "fabricator-generate-plugin-go",
					Spec: fabricatorgenerateplugingo.Spec{
						PluginName:         "fabricator-generate-plugin-go",
						IsGenerationPlugin: true,
					},
				},
			},
		}
		Expect(config).To(Equal(expected))
	})
})
