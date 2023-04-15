package config

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

type Config struct {
	OpenAIAPIKey         string `json:"openai_api_key"`
	GoogleSearchEngineID string `json:"google_search_engine_id"`
	GoogleAPIKey         string `json:"google_api_key"`
}

func Read() (*Config, error) {
	ospath, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := path.Join(ospath, "jacquard-ai", "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	configFile, err := os.Open(configPath)

	if err != nil {
		return nil, err
	}

	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)

	var config Config

	err = json.Unmarshal(byteValue, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil

}
