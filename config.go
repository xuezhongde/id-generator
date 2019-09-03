package conf

import (
    "github.com/BurntSushi/toml"
    "github.com/juju/errors"
    "io/ioutil"
)

type Config struct {
    Port         int    `toml:"port"`
    Router       string `toml:"router"`
    DateCenterId int64  `toml:"date_center_id"`
    WorkerId     int64  `toml:"worker_id"`
}

func LoadConfig(filename string) (*Config, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, errors.Trace(err)
    }

    var c Config
    if _, err := toml.Decode(string(data), &c); err != nil {
        return nil, errors.Trace(err)
    }

    return &c, nil
}
