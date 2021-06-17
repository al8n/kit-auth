package server

import (
	"sync"
)

var (
	once sync.Once
	s *server
)

type server struct {
	onceInit, onceServe, onceClose sync.Once
}

func Get() *server  {
	once.Do(func() {
		s = &server{}
	})
	return s
}



func (s *server) Init() (err error) {
	s.onceInit.Do(func() {
		
	})
	return nil
}

func (s *server) Serve() (err error) {
	s.onceServe.Do(func() {
		
	})
	return nil
}

func (s *server) Close() {
	s.onceClose.Do(func() {
		
	})
	return
}


