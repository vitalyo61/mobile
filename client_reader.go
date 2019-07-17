package mobile

import (
	"log"

	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/bus"
)

type moduleReader struct {
	socket mangos.Socket
}

// Reader ...
type Reader struct {
	registration *Registration
	modules      map[string]*moduleReader
}

// NewReader ...
func NewReader(reg *Registration) *Reader {
	return &Reader{
		registration: reg,
		modules:      make(map[string]*moduleReader),
	}
}

// Run ...
func (r *Reader) Run() {
	go func() {
		chanEventAddModule := r.registration.GetChanEventAddModule()
		select {
		case m := <-chanEventAddModule:
			log.Printf("Module add: %T\n", m)
			uriClient, err := r.registration.getURIModule(m.Name)
			if err != nil {
				log.Println(err)
				break
			}

			mod := new(moduleReader)
			mod.socket, err = bus.NewSocket()
			if err != nil {
				log.Println(err)
				break
			}

			if err = mod.socket.Listen(uriClient); err != nil {
				log.Println(err)
				break
			}

			r.modules[m.Name] = mod
		}
	}()
}

// Close ...
func (r *Reader) Close() {
	for _, m := range r.modules {
		m.socket.Close()
	}
}
