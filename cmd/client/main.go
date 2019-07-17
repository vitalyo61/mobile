package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	// надо
	_ "nanomsg.org/go/mangos/v2/transport/all"

	"gonative/mobile"
)

func main() {
	log.SetFlags(log.Lshortfile)

	reg, err := mobile.NewRegistration("/tmp", "socks")
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := reg.Close(); err != nil {
			log.Println(err)
		}

		// if err := reg.ClearSys(); err != nil {
		// 	log.Println(err)
		// }
	}()

	mt := newModuleTimer()
	log.Println(mt.generateModuleChannel())

	chanExit := make(chan struct{})
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	reader := mobile.NewReader(reg)

	reader.Run()

	err = reg.Run()
	if err != nil {
		log.Panic(err)
	}

	defer reader.Close()

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		chanExit <- struct{}{}
	}()

	<-chanExit

	log.Printf("Reg: %+v\n", reg)
}

type moduleTimer struct {
	typeModule string
	chanIn     chan []byte
	chanOut    chan []byte
}

func newModuleTimer() *moduleTimer {
	return &moduleTimer{
		typeModule: "timer",
		chanIn:     make(chan []byte),
		chanOut:    make(chan []byte),
	}
}

func (m *moduleTimer) generateModuleChannel() *mobile.ModuleChannel {
	return mobile.NewModuleChannel(m.typeModule, m.chanIn, m.chanOut)
}
