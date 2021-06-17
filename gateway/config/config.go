package config

import (
	bootflag "github.com/al8n/micro-boot/flag"
	"sync"
)

var (
	once, onceInit sync.Once
	c *config
)

type config struct {
	
}

func (c *config) Initialize(name string) (err error) {
	panic("implement me")
}

func (c *config) BindFlags(fs *bootflag.FlagSet) {
	panic("implement me")
}

func (c *config) Parse() (err error) {
	panic("implement me")
}

func Get() *config {
	once.Do(func() {
		c = &config{}
	})
	return c
}


func (c *config) Init(filename string) (err error)  {
	onceInit.Do(func() {
		
	})
	return nil
}
