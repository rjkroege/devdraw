/*
	useful utility to encode/decode drawfcall.Msg instances into minimal JSON.

	How do I know if I've left out something important or not? I push everything through
	the encoder / decoder and see if things still work... 
*/

package main

import (
//	"fmt"
//	"log"
//	"code.google.com/p/goplan9/draw"
	"image"
//	"syscall"
	"code.google.com/p/goplan9/draw/drawfcall"
//	"os"
//	"strings"
//	"sync"
//	"io"
//	"encoding/json"
)

type JsonedMsgCore struct {
	Tag uint8
	Type string
}

type JsonedMsgMouse struct {
	JsonedMsgCore
	Mouse drawfcall.Mouse
}

type JsonedMsgCursor struct {
	JsonedMsgCore
	Cursor  drawfcall.Cursor
	Arrow   bool
}

type JsonedMsgLabel struct {
	JsonedMsgCore
	Label   string
}

type JsonTinit struct {
	JsonedMsgCore
	Label   string
	Winsize string
}

type JsonRrdkbd struct {
	JsonedMsgCore
	Rune    rune
}

type JsonSnarf struct {
	JsonedMsgCore
	Snarf   string
}

type JsonDrawCount struct {
	JsonedMsgCore
	Count   int
}

type JsonResize struct {
	JsonedMsgCore
	Rect    image.Rectangle
}

func PrettyJsonOutput(m* drawfcall.Msg) (interface{}) {
	switch m.Type {
	default:
		return &JsonedMsgCore{m.Tag, "unknown"}
	case drawfcall.Rerror:
		return &JsonedMsgCore{m.Tag, "Rerror"}
	case drawfcall.Trdmouse:
		return&JsonedMsgCore{m.Tag, "Trdmouse"}
	case drawfcall.Rrdmouse:
		return &JsonedMsgMouse{JsonedMsgCore{m.Tag,  "Rrdmouse"},  m.Mouse}
	case drawfcall.Tbouncemouse:
		return &JsonedMsgMouse{JsonedMsgCore{m.Tag, "Tbouncemouse"}, m.Mouse}
	case drawfcall.Rbouncemouse:
		return &JsonedMsgCore{m.Tag, "Rbouncemouse"}
	case drawfcall.Tmoveto:
		return &JsonedMsgMouse{JsonedMsgCore{m.Tag, "Tmoveto"}, m.Mouse}
	case drawfcall.Rmoveto:
		return &JsonedMsgCore{m.Tag, "Rmoveto"}
	case drawfcall.Tcursor:
		return &JsonedMsgCursor{JsonedMsgCore{m.Tag, "Tcursor"}, m.Cursor, m.Arrow}
	case drawfcall.Rcursor:
		return &JsonedMsgCore{m.Tag, "Rcursor"}
	case drawfcall.Trdkbd:
		return &JsonedMsgCore{m.Tag, "Trdkbd"}
	case drawfcall.Rrdkbd:
		return &JsonRrdkbd{JsonedMsgCore{m.Tag, "Rrdkbd"}, m.Rune}
	case drawfcall.Tlabel:
		return &JsonedMsgLabel{JsonedMsgCore{m.Tag, "Tlabel"}, m.Label}
	case drawfcall.Rlabel:
		return &JsonedMsgCore{m.Tag, "Rlabel"}
	case drawfcall.Tinit:
		return &JsonTinit{JsonedMsgCore{m.Tag, "Tinit"}, m.Label, m.Winsize}
	case drawfcall.Rinit:
		return &JsonedMsgCore{m.Tag, "Rinit"}
	case drawfcall.Trdsnarf:
		return &JsonedMsgCore{m.Tag, "Trdsnarf"}
	case drawfcall.Rrdsnarf:
		return &JsonSnarf{JsonedMsgCore{m.Tag, "Rrdsnarf"}, m.Snarf}
	case drawfcall.Twrsnarf:
		return &JsonSnarf{JsonedMsgCore{m.Tag, "Twrsnarf"}, m.Snarf}
	case drawfcall.Rwrsnarf:
		return &JsonedMsgCore{m.Tag, "Rwrsnarf"}
	case drawfcall.Trddraw:
		return &JsonDrawCount{JsonedMsgCore{m.Tag, "Trddraw"}, m.Count}
	case drawfcall.Rrddraw:
		return &JsonedMsgCore{m.Tag, "Rrddraw - expand data", }
	case drawfcall.Twrdraw:
		return &JsonedMsgCore{m.Tag, "Twrdraw - expand data", }
	case drawfcall.Rwrdraw:
		return &JsonDrawCount{JsonedMsgCore{m.Tag, "Rwrdraw"}, m.Count}
	case drawfcall.Ttop:
		return &JsonedMsgCore{m.Tag, "Ttop"}
	case drawfcall.Rtop:
		return &JsonedMsgCore{m.Tag, "Rtop"}
	case drawfcall.Tresize:
		return &JsonResize{JsonedMsgCore{m.Tag, "Tresize"}, m.Rect}
	case drawfcall.Rresize:
		return &JsonedMsgCore{m.Tag, "Rresize"}
	}
}

// TODO(rjkroege):
// write a matching command to convert the pretty json
// to drawfcall.Msg.
