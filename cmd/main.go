package main

import (
	"log"
	"time"

	// надо
	_ "nanomsg.org/go/mangos/v2/transport/all"

	"gonative/mobile"
)

var sockName = "/reqrep.ipc"

func main() {
	log.SetFlags(log.Lshortfile)

	// dir, err := ioutil.TempDir("/tmp", "socket")
	// if err != nil {
	// 	log.Panicln(err)
	// }

	// sockName = dir + sockName

	// defer func() {
	// 	if err := os.RemoveAll(dir); err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	// var wg sync.WaitGroup

	// go  registration()

	reg, err := mobile.NewRegistration("/tmp", "qqq")
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := reg.ClearSys(); err != nil {
			log.Println(err)
		}
	}()

	err = reg.Run()
	if err != nil {
		log.Panic(err)
	}

	time.Sleep(time.Millisecond * 20)

	log.Printf("Reg: %+v\n", reg)

	module, err := mobile.NewModuleRegistration("module_1", "module", "/tmp", "qqq")
	if err != nil {
		log.Panic(err)
	}

	err = module.Run()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Reg: %+v\n", reg)
	log.Printf("Cl: %+v\n", module)

	// wg.Add(1)
	// go func() {
	// 	module("module_1")
	// 	wg.Done()
	// }()

	// wg.Add(1)
	// go func() {
	// 	module("module_2")
	// 	wg.Done()
	// }()

	// wg.Wait()
}

// func registration() (reg *mobile.Registration, err error) {

// 	var (
// 		sock mangos.Socket
// 		err  error
// 		msg  *mangos.Message
// 	)

// 	if sock, err = rep.NewSocket(); err != nil {
// 		log.Panicf("%s%v\n", prefix, err)
// 	}

// 	if err = sock.Listen("ipc://" + sockName); err != nil {
// 		log.Panicf("%s%v\n", prefix, err)
// 	}

// 	for {
// 		if msg, err = sock.RecvMsg(); err != nil {
// 			log.Printf("%s%v\n", prefix, err)
// 			continue
// 		}

// 		log.Printf("%sRecive: %s\n", prefix, msg.Body)

// 		if err = sendMessage(sock, []byte(fmt.Sprintf("Replay on '%s'", msg.Body))); err != nil {
// 			log.Printf("%s%v\n", prefix, err)
// 		}
// 	}
// }

// func module(name string) {
// 	prefix := fmt.Sprintf("Module (%s): ", name)

// 	var (
// 		sock mangos.Socket
// 		err  error
// 		msg  *mangos.Message
// 	)

// 	if sock, err = req.NewSocket(); err != nil {
// 		log.Panicf("%s%v\n", prefix, err)
// 	}

// 	if err = sock.Dial("ipc://" + sockName); err != nil {
// 		log.Panicf("%s%v\n", prefix, err)
// 	}

// 	counter := 0
// 	for {
// 		if err = sendMessage(sock, []byte(fmt.Sprintf("%s send %d", name, counter))); err != nil {
// 			log.Printf("%s%v\n", prefix, err)
// 		}

// 		if msg, err = sock.RecvMsg(); err != nil {
// 			log.Printf("%s%v\n", prefix, err)
// 			continue
// 		}

// 		log.Printf("%sRecive: %s\n", prefix, msg.Body)

// 		// time.Sleep(time.Millisecond * 10)

// 		if counter > 9 {
// 			break
// 		}
// 		counter++
// 	}

// }

// func sendMessage(sock mangos.Socket, text []byte) error {
// 	msg := mangos.NewMessage(len(text))
// 	msg.Body = text
// 	return sock.SendMsg(msg)
// }
