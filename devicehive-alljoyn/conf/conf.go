package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// All parameters are optional. See defaults in the fix()
type Conf struct {
	AJServiceConnectionTimeoutMilliseconds   uint32 `yaml:"AJServiceConnectionTimeoutMilliseconds,omitempty"`
	AJServiceMsgUnmarshalTimeoutMilliseconds uint32 `yaml:"AJServiceMsgUnmarshalTimeoutMilliseconds,omitempty"`

	AJServiceConnectionPort uint16 `yaml:"AJServiceConnectionPort,omitempty"`
}

func (c *Conf) fix() {
	if c.AJServiceConnectionTimeoutMilliseconds == 0 {
		c.AJServiceConnectionTimeoutMilliseconds = 60 * 1000
	}

	if c.AJServiceMsgUnmarshalTimeoutMilliseconds == 0 {
		c.AJServiceMsgUnmarshalTimeoutMilliseconds = 5 * 1000
	}

	if c.AJServiceConnectionPort == 0 {
		c.AJServiceConnectionPort = 25
	}
}

func TestConf() *Conf {
	c := new(Conf)
	c.fix()
	return c
}

func FromArgs() (filepath string, c *Conf, err error) {
	parseArgs()
	filepath = confArgValue
	if len(filepath) == 0 {
		c = TestConf()
		return
	}
	c, err = readConf(confArgValue)
	return
}

func readConf(filepath string) (c *Conf, err error) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, c)

	if err == nil {
		c.fix()
	}

	return
}
