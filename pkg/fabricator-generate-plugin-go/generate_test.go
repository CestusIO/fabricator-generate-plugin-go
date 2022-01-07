package fabricatorgenerateplugingo

import (
	"context"
	"os"

	"code.cestus.io/libs/codegenerator/pkg/templating"
	"code.cestus.io/tools/fabricator-generate-plugin-go/pkg/fabricator-generate-plugin-go/templates"
	"code.cestus.io/tools/fabricator/pkg/fabricator"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Generation", func() {
	ginkgo.It("Generates", func() {
		ctx := context.Background()
		packProvider := templating.NewPackProvider()
		packProvider.RegisterProvider(templating.NewEmbededPackProvider(templates.GetTemplates()))
		io := fabricator.NewTestIOStreamsDiscard()
		file, err := os.Open("./testdata/deserialize.yml")
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		config, err := LoadPluginConfig(file)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		root, err := os.MkdirTemp("./testdata", "generation")
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		plugin, err := newPlugin(ctx, io, config, root, packProvider)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		err = plugin.Generate(ctx, io)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		err = os.RemoveAll(root)
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})
})
