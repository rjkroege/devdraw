/*

	A devdraw listener. Maybe I've named this wrong.
	It's the devdraw interceptor

*/

package main

import (
	//	"fmt"
	"log"
	//	"code.google.com/p/goplan9/draw"
	//	"image"
	"9fans.net/go/draw/drawfcall"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
)

/*
	Strip DEVDRAW from environment so that the plan9 go code
	Doesn't invoke the interceptor program recursively.
	TODO(rjkroege): Note that I might explicitly want to specify which
	DEVDRAW is to be invoked for real.
*/
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

type App struct {
	w sync.Mutex
	o *os.File
	i *os.File
}

func checkedClose(f *os.File, msg string) {
	err := f.Close()
	if err != nil {
		log.Fatal(msg, err)
	}
}

func marshalsxtx(inbuffy []byte, app *App, devdraw *drawfcall.Conn, json *JsonRecorder) {
	tx := new(drawfcall.Msg)
	rx := new(drawfcall.Msg)

	tag := inbuffy[4]
	log.Print("set inbuffy to something, tag ", tag)

	log.Print("bar")
	err := tx.Unmarshal(inbuffy)
	json.Record(tx, tag)
	if err != nil {
		log.Fatal("build a msg: ", err)
	}

	// Write message to real devdraw and get response.
	log.Print("sending tx to real devdraw, getting rx back")
	err = devdraw.RPC(tx, rx)
	if err != nil {
		if err != io.EOF {
			log.Print("send/receive to real devdraw had error: ", err)
		}
		app.w.Lock()
		checkedClose(app.o, "couldn't close channel to host: ")
		app.w.Unlock()
		checkedClose(app.i, "Couldn't close channel from host: ")
		return
	}

	// TODO(rjkroege): Time-stamp the records.
	// TODO(rjkroege): I want the actual original tag.
	json.Record(rx, tag)

	// write to cout
	outbuffy := rx.Marshal()
	// log.Print("returned tag ", outbuffy[4])
	// log.Print("changing tag")
	outbuffy[4] = tag

	app.w.Lock()
	_, err = app.o.Write(outbuffy)
	app.w.Unlock()
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

	out2, err := syscall.Dup(1)
	if err != nil {
		log.Fatal("dupping 1", err)
	}

	os.Stdin.Close()
	os.Stdout.Close()
	os.Stdout = os.NewFile(uintptr(2), "/dev/stdout")

	/* Connections to the application. */
	cin := os.NewFile(uintptr(in2), "fromapp")
	cout := os.NewFile(uintptr(out2), "toapp")

	modifyEnvironment()

	// Fire up a new devdraw here
	devdraw, err := drawfcall.New()
	if err != nil {
		log.Fatal("making a Conn", err)
	}

	// There is probably a nicer way to do this.
	// TODO(rjkroege): do it the nicer way.
	var app App
	app.o = cout
	app.i = cin

	json := NewJsonRecorder()

	for {
		// read crap from cin
		log.Print("about to read from host")
		inbuffy, err := drawfcall.ReadMsg(cin)
		log.Print("read from host")
		if err != nil {
			devdraw.Close()
			break
		}
		go marshalsxtx(inbuffy, &app, devdraw, json)
	}

	log.Print("waiting on completion")
	json.WaitToComplete()
}
