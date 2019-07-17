package mobile

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gogo/protobuf/proto"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/rep"

	"gonative/mobile/data"
)

// Module ...
type Module struct {
	data.Module
	URI string
}

// Registration ...
type Registration struct {
	modules            map[string]*Module
	uriDir             string
	sock               mangos.Socket
	chanEventAddModule chan *data.Module
}

// NewRegistration ...
func NewRegistration(dir, nameSocket string) (reg *Registration, err error) {
	reg = new(Registration)

	reg.modules = make(map[string]*Module)

	if reg.sock, err = rep.NewSocket(); err != nil {
		return
	}

	switch {
	case dir == "" && nameSocket == "":
		dir = "/tmp"
		nameSocket = "socket"
	case nameSocket == "":
		err = errors.New("haven't nameSocket")
	}

	reg.uriDir = fmt.Sprintf("%s/%s", dir, nameSocket)

	if err = os.MkdirAll(reg.uriDir, os.ModeDir|0700); err != nil {
		return
	}

	reg.chanEventAddModule = make(chan *data.Module)

	return
}

// GetChanEventAddModule ...
func (r *Registration) GetChanEventAddModule() chan *data.Module {
	return r.chanEventAddModule
}

func (r *Registration) addModule(module *data.Module) (string, error) {
	if _, ok := r.modules[module.Name]; ok {
		return "", errors.New("name's module exists")
	}

	c := &Module{
		Module: *module,
		URI:    fmt.Sprintf("ipc://%s/%s.ipc", r.uriDir, module.Name),
	}
	r.modules[module.Name] = c

	log.Printf("Module (%T): %+v, %s\n", c, c, c.URI)
	r.chanEventAddModule <- module

	return c.URI, nil
	// return nil
}

// GetURI ...
func (r *Registration) GetURI() string {
	return fmt.Sprintf("ipc://%s/%s.ipc", r.uriDir, registrationNameSock)
}

// // GetURIModule ...
// func (r *Registration) GetURIModule(name string) (string, error) {
// 	m, ok := r.modules[name]
// 	if !ok {
// 		return "", errors.New("haven't module: " + name)
// 	}
// 	return fmt.Sprintf("ipc://%s/%s.ipc", r.uriDir, m.UriClient), nil
// }

func (r *Registration) getURIModule(name string) (string, error) {
	if _, ok := r.modules[name]; ok {
		return fmt.Sprintf("ipc://%s/%s.ipc", r.uriDir, name), nil
	}
	return "", errors.New("haven't module " + name)
}

// func (r *Registration) getSock() mangos.Socket {
//   return r.sock
// }

// Close ...
func (r *Registration) Close() error {
	for _, m := range r.modules {
		log.Printf("%T: %+v (%s)\n", m, m, m.URI)
	}
	return r.sock.Close()
}

// ClearSys ...
func (r *Registration) ClearSys() error {
	return os.RemoveAll(r.uriDir)
}

// Run ...
func (r *Registration) Run() error {
	if err := r.sock.Listen(r.GetURI()); err != nil {
		return err
	}

	go func() {
		var (
			msg     *mangos.Message
			err     error
			respMsg []byte
		)

		for {
			if msg, err = r.sock.RecvMsg(); err != nil {
				// log.Println(err)
				continue
			}

			m := new(data.Module)
			resp := new(data.ModuleResp)
			var uriClient string

			if err = proto.Unmarshal(msg.Body, m); err != nil {
				log.Println(err)
				resp.Error = err.Error()
			} else if uriClient, err = r.addModule(m); err != nil {
				// } else if err = r.addModule(m); err != nil {
				log.Println(err)
				resp.Error = err.Error()
			} else {
				resp.Uri = fmt.Sprintf("ipc://%s/%s.ipc", r.uriDir, busModulesSock)
				resp.UriClient = uriClient
			}
			msg.Free()

			log.Printf("Resp: %+v\n", resp)

			if respMsg, err = proto.Marshal(resp); err != nil {
				log.Println(err)
			}

			if err = sendMessage(r.sock, respMsg); err != nil {
				log.Println(err)
			}
		}
	}()

	return nil
}
