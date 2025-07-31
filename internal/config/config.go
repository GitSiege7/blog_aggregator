package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

const configFileName string = "/.gatorconfig.json"

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home_dir + configFileName), nil
}

func Read() (Config, error) {
	config_path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(config_path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var conf Config

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}

func write(c Config) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	config_path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(config_path, bytes, os.FileMode(0600))
	if err != nil {
		return err
	}

	return nil
}

func (c Config) SetUser(username string) error {
	c.Current_user_name = username

	err := write(c)
	if err != nil {
		return err
	}

	return nil
}
