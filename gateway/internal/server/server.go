package server

import (
	"io"
	"net/http"
	"sync"
)

var (
	once sync.Once
	s *Server
)

type Server struct {
	tracerCloser io.Closer
	handler http.Handler
	wg sync.WaitGroup
}

func Get() *Server  {
	once.Do(func() {
		s = &Server{}
	})
	return s
}



func (s *Server) Init() (err error) {
	return nil
}

func (s *Server) Serve() (err error) {
	return nil
}

func (s *Server) Close() (err error) {
	s.wg.Done()
	return nil
}


