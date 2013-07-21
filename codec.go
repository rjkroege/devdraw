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

type JsonedMsg struct {
	Type string				`json:",omitempty"`
	Tag uint8				`json:",omitempty"`
	Mouse drawfcall.Mouse	`json:",omitempty"`
	Resized bool				`json:",omitempty"`
	Cursor  drawfcall.Cursor	`json:",omitempty"`
	Arrow   bool				`json:",omitempty"`
	Rune    rune				`json:",omitempty"`
	Winsize string			`json:",omitempty"`
	Label   string				`json:",omitempty"`
	Snarf   string				`json:",omitempty"`
	Error   string				`json:",omitempty"`
	// TODO(rjkroege): do I need to care about this?
	// Data    []byte
	Count   int				`json:",omitempty"`
	Rect    image.Rectangle		`json:",omitempty"`
}

func PrettyJsonOutput(m* drawfcall.Msg) (*JsonedMsg) {
	var j *JsonedMsg
	switch m.Type {
	default:
		j = &JsonedMsg{Tag: m.Tag, Type: "unknown"}
	case drawfcall.Rerror:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rerror"}
	case drawfcall.Trdmouse:
		j = &JsonedMsg{Tag: m.Tag, Type: "Trdmouse"}
	case drawfcall.Rrdmouse:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rrdmouse", Mouse: m.Mouse}
	case drawfcall.Tbouncemouse:
		j = &JsonedMsg{Tag: m.Tag, Type: "Tbouncemouse", Mouse: m.Mouse}
	case drawfcall.Rbouncemouse:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rbouncemouse"}
	case drawfcall.Tmoveto:
		j = &JsonedMsg{Tag: m.Tag, Type: "Tmoveto", Mouse: m.Mouse}
	case drawfcall.Rmoveto:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rmoveto"}
	case drawfcall.Tcursor:
		j = &JsonedMsg{Tag: m.Tag, Type: "Tcursor", Arrow: m.Arrow}
	case drawfcall.Rcursor:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rcursor"}
	case drawfcall.Trdkbd:
		j = &JsonedMsg{Tag: m.Tag, Type: "Trdkbd"}
	case drawfcall.Rrdkbd:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rrdkbd", Rune: m.Rune}
	case drawfcall.Tlabel:
		j = &JsonedMsg{Tag: m.Tag, Type: "Tlabel", Label: m.Label}
	case drawfcall.Rlabel:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rlabel"}
	case drawfcall.Tinit:
		j = &JsonedMsg{Tag: m.Tag, Type: "Tinit", Label: m.Label, Winsize: m.Winsize}
	case drawfcall.Rinit:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rinit"}
	case drawfcall.Trdsnarf:
		j = &JsonedMsg{Tag: m.Tag, Type: "Trdsnarf"}
	case drawfcall.Rrdsnarf:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rrdsnarf", Snarf: m.Snarf}
	case drawfcall.Twrsnarf:
		j = &JsonedMsg{Tag: m.Tag, Type: "Twrsnarf", Snarf: m.Snarf}
	case drawfcall.Rwrsnarf:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rwrsnarf"}
	case drawfcall.Trddraw:
		j = &JsonedMsg{Tag: m.Tag, Type: "Trddraw", Count: m.Count}
	case drawfcall.Rrddraw:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rrddraw - need data", }
	case drawfcall.Twrdraw:
		j = &JsonedMsg{Tag: m.Tag, Type: "Twrdraw - need data", }
	case drawfcall.Rwrdraw:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rwrdraw", Count: m.Count}
	case drawfcall.Ttop:
		j = &JsonedMsg{Tag: m.Tag, Type: "Ttop"}
	case drawfcall.Rtop:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rtop"}
	case drawfcall.Tresize:
		j = &JsonedMsg{Tag: m.Tag, Type: "Tresize", Rect: m.Rect}
	case drawfcall.Rresize:
		j = &JsonedMsg{Tag: m.Tag, Type: "Rresize"}
	}
	return j
}

// TODO(rjkroege):
// write a matching command to convert the pretty json
// to drawfcall.Msg.
