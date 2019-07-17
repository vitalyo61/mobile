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

	module, err := mobile.NewModuleRegistration("module_1", "timer", "/tmp", "socks")
	if err != nil {
		log.Panic(err)
	}

	err = module.Run()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Mod: %+v\n", module)

	chanIn := make(chan []byte)
	chanOut := make(chan []byte)

	worker, err := mobile.NewWorkerModule(module.GetURI(), module.GetURIClient(), chanIn, chanOut)
	if err != nil {
		log.Panic(err)
	}

	chanExit := make(chan struct{})
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	worker.Run()

	defer func() {
		if err := worker.Close(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		chanExit <- struct{}{}
	}()

	<-chanExit

}
