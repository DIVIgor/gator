package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)


const cfgFilename string = ".gatorconfig.json"


type Config struct {
    DbUrl string `json:"db_url"`
    User string `json:"current_user_name"`
}


func (c *Config) SetUser(username string) (err error) {
    c.User = username

    return write(*c)
}

// get a full path to the config file
func getConfigPath() (path string, err error) {
    path, err = os.UserHomeDir()
    if err != nil {
        return
    }
    return filepath.Join(path, cfgFilename), err
}

// read the config file from home dir
func Read() (cfg Config, err error) {
    path, err := getConfigPath()
    if err != nil {return}

    file, err := os.ReadFile(path)
    if err != nil {return}

    err = json.Unmarshal(file, &cfg)

    return cfg, err
}

// write the config to the home dir
func write(cfg Config) (err error) {
    data, err := json.Marshal(cfg)
    if err != nil {return}

    fpath , err := getConfigPath()
    if err != nil {return}

    err = os.WriteFile(fpath, data, 0666)
    return err
}