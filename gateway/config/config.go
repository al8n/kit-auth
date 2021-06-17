package config

import "sync"

var (
	once, onceInit sync.Once
	c *config
)

type config struct {
	
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
