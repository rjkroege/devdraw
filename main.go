
/*

	A devdraw listener. Maybe I've named this wrong.

*/

package main

import (
//	"fmt"
	"log"
//	"code.google.com/p/goplan9/draw"
//	"image"
	"syscall"
	"code.google.com/p/goplan9/draw/drawfcall"
	"os"
	"strings"
	"sync"
)



func modifyEnvironment() {
	envs := os.Environ()
	os.Clearenv()

	for _, v := range envs {
		// log.Print("env contents", v)
		if !strings.HasPrefix(v, "DEVDRAW") {
			ss := strings.Split(v, "=")
			// log.Print("env chunks", ss)
			err := os.Setenv(ss[0], ss[1])
			if err != nil {
				log.Fatal("setting env")
			}
		} else {
			log.Print("clearing DEVDRAW from environment")
		}
	}
}

type AppConn struct {
	w  sync.Mutex
	o *os.File
}

func marshalsxtx(inbuffy []byte, appconn *AppConn, devdraw *drawfcall.Conn) {
	tx := new(drawfcall.Msg)
	rx := new(drawfcall.Msg)

	tag := inbuffy[4]
	log.Print("set inbuffy to something, tag ", tag)

	log.Print("bar")
	err := tx.Unmarshal(inbuffy)
	log.Print("parsed into a tx: ", tx)
	if err != nil {
		log.Fatal("build a msg: ", err)
	}

	log.Print("foo")

	// Write message to real devdraw and get response.
	log.Print("sending tx to real devdraw, getting rx back")
	err = devdraw.RPC(tx, rx)
	if err != nil {
		log.Fatal("send/receive to real devdraw: ", err)
	}
	log.Print("got rx back: ", rx)

	// write to cout
	outbuffy := rx.Marshal()
	log.Print("returned tag ", outbuffy[4])
	log.Print("changing tag")
	outbuffy[4] = tag
	
	appconn.w.Lock()
	_, err = appconn.o.Write(outbuffy)
	appconn.w.Unlock()
	if err != nil {
		log.Fatal("write to app: ", err)
	}		
}

func main() {
	// I assume that in is 0
	in2, err := syscall.Dup(0)
	if err != nil {
		log.Fatal("dupping 0", err)
	}

	out2, err  := syscall.Dup(1)
	if err != nil {
		log.Fatal("dupping 1", err)
	}

	os.Stdin.Close()
	os.Stdout.Close()
	os.Stdout = os.NewFile(uintptr(2), "/dev/stdout")

	cin := os.NewFile(uintptr(in2), "fromapp")
	cout := os.NewFile(uintptr(out2), "toapp")

	// TODO(rjkroege): Modify environment here.
	modifyEnvironment();

	// Fire up a new devdraw here
	devdraw, err  := drawfcall.New()
	if err != nil {
		log.Fatal("making a Conn", err)
	}

	log.Print("hello")
	var appconn AppConn
	appconn.o = cout

	for {
		// read crap from cin
		inbuffy, err := drawfcall.ReadMsg(cin)
		if err != nil {
			log.Fatal("read from app: ", err)
			break;
		}

		go marshalsxtx(inbuffy, &appconn, devdraw)
	}
}
