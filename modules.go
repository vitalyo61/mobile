package mobile

import "errors"

// ModuleChannel ...
type ModuleChannel struct {
	typeModule string
	chanIn     chan []byte
	chanOut    chan []byte
}

// NewModuleChannel ...
func NewModuleChannel(typeModule string, in, out chan []byte) *ModuleChannel {
	return &ModuleChannel{
		typeModule: typeModule,
		chanIn:     in,
		chanOut:    out,
	}
}

// ModuleChannels ...
type ModuleChannels struct {
	modules []*ModuleChannel
}

// NewModuleChannels ...
func NewModuleChannels() *ModuleChannels {
	return &ModuleChannels{
		modules: make([]*ModuleChannel, 0),
	}
}

// Add ...
func (m *ModuleChannels) Add(f *ModuleChannel) error {
	for _, mf := range m.modules {
		if mf.typeModule == f.typeModule {
			return errors.New("module exist: " + f.typeModule)
		}
	}
	m.modules = append(m.modules, f)
	return nil
}

// Del ...
func (m *ModuleChannels) Del(f *ModuleChannel) {
	for i, mf := range m.modules {
		if mf.typeModule == f.typeModule {
			m.modules = append(m.modules[0:i], m.modules[i+1:len(m.modules)]...)
			return
		}
	}
}

// Get ...
func (m *ModuleChannels) Get(typeModule string) *ModuleChannel {
	for _, mf := range m.modules {
		if mf.typeModule == typeModule {
			return mf
		}
	}
	return nil
}
