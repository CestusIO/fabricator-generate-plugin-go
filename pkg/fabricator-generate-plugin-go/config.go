// Generated with the code-generator
//
// Modifications in code regions will be lost during regeneration!

package fabricatorgenerateplugingo

import (
	"fmt"
	"io"

	"code.cestus.io/tools/fabricator/pkg/fabricator"
	"gopkg.in/yaml.v3"
)

// Spec contains the specification for the plugin
type Spec struct {
	PluginName         string `yaml:"pluginName" json:"pluginName"`
	IsGenerationPlugin bool   `yaml:"isGenerationPlugin" json:"isGenerationPlugin"`
}

// region CODE_REGION(PLUGIN_COMPONENT)
type PluginComponent struct {
	Name      string `yaml:"name" json:"name"`
	Generator string `yaml:"generator" json:"generator"`
	Spec      Spec   `yaml:"spec" json:"spec"`
}

type PluginComponents []PluginComponent

type PluginConfig struct {
	ApiVersion string           `yaml:"apiVersion" json:"apiVersion"`
	Kind       string           `yaml:"kind" json:"kind"`
	Components PluginComponents `yaml:"components" json:"components"`
}

func LoadPluginConfig(r io.Reader) (PluginConfig, error) {
	pluginConfig := PluginConfig{}
	var err error
	if err = yaml.NewDecoder(r).Decode(&pluginConfig); err != nil {
		return pluginConfig, fmt.Errorf("failed to decode pluginConfig %s", err)
	}

	return pluginConfig, err
}

// endregion

func (s *PluginConfig) UnmarshalYAML(value *yaml.Node) error {

	var config fabricator.FabricatorConfig
	if err := value.Decode(&config); err != nil {
		return err
	}
	s.ApiVersion = config.ApiVersion
	s.Kind = config.Kind

	for _, component := range config.Components {

		if component.Generator != PluginName {
			continue
		}
		spec := Spec{}
		err := component.Spec.Decode(&spec)
		if err != nil {
			return err
		}
		plComponent := PluginComponent{
			Name:      component.Name,
			Generator: component.Generator,
			Spec:      spec,
		}
		s.Components = append(s.Components, plComponent)
	}
	return nil
}
