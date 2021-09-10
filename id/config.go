// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package id

import (
	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
	"io/ioutil"
)

type Config struct {
	AppName       string `toml:"appName"`
	Profile       string `toml:"profile"`
	Port          int    `toml:"port"`
	Router        string `toml:"router"`
	DateCenterId  int64  `toml:"date_center_id"`
	WorkerId      int64  `toml:"worker_id"`
	ConnectString string `toml:"connectString"`
	NodePath      string `toml:"nodePath"`
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
