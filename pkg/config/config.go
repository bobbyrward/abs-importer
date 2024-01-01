package config

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConfigConfig struct {
	ConfigFilename string
}

type ConfigOptions func(*ConfigConfig)

func WithConfigFilename(filename string) ConfigOptions {
	return func(cc *ConfigConfig) {
		cc.ConfigFilename = filename
	}
}

func changeMapstructureTags(tagName string) viper.DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.TagName = tagName
	}
}

func decoderHookFunc() mapstructure.DecodeHookFunc {
	return mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)
}

func NewConfig(opts ...ConfigOptions) (*Config, error) {
	cc := ConfigConfig{}
	config := &Config{}

	for _, opt := range opts {
		opt(&cc)
	}

	if cc.ConfigFilename != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cc.ConfigFilename)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		// TODO: This should be in ~/.config
		viper.SetConfigName(".abs-importer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(config, changeMapstructureTags("json"), viper.DecodeHook(decoderHookFunc()))
	if err != nil {
		return nil, err
	}

	return config, nil
}
