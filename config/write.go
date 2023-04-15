package config

import (
	"encoding/json"
	"os"
	"path"
)

func Write(conf *Config) error {
	ospath, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	content, err := json.MarshalIndent(conf, "", " ")

	if err != nil {
		return err
	}

	configPath := path.Join(ospath, "jacquard-ai", "config.json")

	os.MkdirAll(path.Dir(configPath), 0755)

	return os.WriteFile(configPath, content, 0644)
}
