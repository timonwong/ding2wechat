// Copyright (c) Timon Wong
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.UnmarshalStrict(content, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Config for the ding2wechat program.
type Config struct {
	Receivers []Receiver `yaml:"receivers"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	trackReceivers := make(map[string]struct{})
	for _, receiver := range c.Receivers {
		if _, ok := trackReceivers[receiver.Name]; ok {
			return fmt.Errorf("duplicate receiver: %s", receiver.Name)
		}
		trackReceivers[receiver.Name] = struct{}{}
	}

	return nil
}

type Receiver struct {
	Name    string   `yaml:"name"`
	Targets []Target `yaml:"targets"`
}

type Target struct {
	MentionedList       []string `yaml:"mentioned_list"`
	MentionedMobileList []string `yaml:"mentioned_mobile_list"`
	URL                 string   `yaml:"url"`
}
