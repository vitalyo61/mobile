package mobile

import (
	"errors"
	"fmt"
	"log"

	"github.com/gogo/protobuf/proto"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/bus"
	"nanomsg.org/go/mangos/v2/protocol/req"

	"gonative/mobile/data"
)

// ModuleRegistration ...
type ModuleRegistration struct {
	name            string
	typeModule      string
	sock            mangos.Socket
	uriRegistration string // registration's uri
	uri             string // module's uri
	uriClient       string // client's uri
}

// NewModuleRegistration ...
func NewModuleRegistration(name, typeModule, dir, nameSocket string) (module *ModuleRegistration, err error) {
	module = &ModuleRegistration{
		name:            name,
		typeModule:      typeModule,
		uriRegistration: fmt.Sprintf("ipc://%s/%s/%s.ipc", dir, nameSocket, registrationNameSock),
	}

	module.sock, err = req.NewSocket()
	return
}

// GetURI ...
func (m *ModuleRegistration) GetURI() string {
	return m.uri
}

// GetURIClient ...
func (m *ModuleRegistration) GetURIClient() string {
	return m.uriClient
}

// Run ...
func (m *ModuleRegistration) Run() error {
	var (
		err     error
		msg     []byte
		message *mangos.Message
	)

	// BEGIN: Registration
	if err = m.sock.Dial(m.uriRegistration); err != nil {
		return err
	}

	req := &data.Module{
		Name: m.name,
		Type: m.typeModule,
	}

	if msg, err = proto.Marshal(req); err != nil {
		return err
	}

	if err = sendMessage(m.sock, msg); err != nil {
		return err
	}

	if message, err = m.sock.RecvMsg(); err != nil {
		return err
	}

	resp := new(data.ModuleResp)
	if err = proto.Unmarshal(message.Body, resp); err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	m.uri = resp.Uri
	m.uriClient = resp.UriClient

	// END: Registration

	return nil
}

// WorkerModule ...
type WorkerModule struct {
	sock    mangos.Socket
	chanIn  chan []byte
	chanOut chan []byte
}

// NewWorkerModule ...
func NewWorkerModule(uri, uriClient string, chanIn, chanOut chan []byte) (*WorkerModule, error) {
	var (
		sock mangos.Socket
		err  error
	)

	if sock, err = bus.NewSocket(); err != nil {
		return nil, err
	}

	if err = sock.Listen(uri); err != nil {
		return nil, err
	}

	if err = sock.Dial(uriClient); err != nil {
		return nil, err
	}

	return &WorkerModule{
		sock: sock,
	}, nil
}

// Close ...
func (w *WorkerModule) Close() error {
	return w.sock.Close()
}

// Run ...
func (w *WorkerModule) Run() {
	go func() {
		var (
			msg *mangos.Message
			err error
		)
		for {
			if msg, err = w.sock.RecvMsg(); err != nil {
				log.Println(err)
				continue
			}
			w.chanOut <- msg.Body
			msg.Free()
		}
	}()
}

// Send ...
func (w *WorkerModule) Send(d []byte) error {
	return sendMessage(w.sock, d)
}
