/*
	useful utility to encode/decode drawfcall.Msg instances into minimal JSON.

	How do I know if I've left out something important or not? I push everything through
	the encoder / decoder and see if things still work... 
*/

package main

import (
	"code.google.com/p/goplan9/draw/drawfcall"
	"code.google.com/p/goplan9/draw"
	"encoding/binary"
	"image"
	"log"
//	"code.google.com/p/goplan9/draw"
//	"encoding/json"
//	"fmt"
//	"io"
//	"log"
//	"os"
//	"strings"
//	"sync"
//	"syscall"
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

// Twrdraw data bundles for when we don't know what to send.
type TwrdrawDefault struct {
	Type byte
}

/*

	I don't want to create a separate type for each kind fo message. As you
	see below, there are lots of them. There needs to be a better way. A more
	sophisticated scheme.  Nope. I need a struct for each kind. But I want an
	array of cmds where each cmd is a struct type. How do I make a polymorphic
	array in go?

	Pretty easily. With an interface.

	Draw needs its own file. Probably needs its own namespace.

*/

type DrawType struct {
	CmdType byte
}

type DrawBlit struct {
	DrawType
	DstID uint32
	SrcID uint32
	MaskID uint32
	DstRect image.Rectangle
	SrcPt image.Point
	MaskPt image.Point
}

type DrawOp struct {
	DrawType
	Op string
}

type DrawCmd interface {
	Type() byte;
}

func (d *DrawType) Type() byte {
	return d.CmdType
}

type JsonedTwrdrawDraw struct {
	JsonedMsgCore
	Cmds []DrawCmd
}

var opnames [13]string
func init() {
	opnames[draw.Clear] = "Clear"
	opnames[draw.SinD] = "SinD"
	opnames[draw.DinS] = "DinS"
	opnames[draw.SoutD] = "SoutD"
	opnames[draw.DoutS] = "DoutS"
	opnames[draw.SinD | draw.SoutD] = "S"
	opnames[draw.SinD | draw.SoutD | draw.DoutS] = "SoverD"
	opnames[draw.SinD | draw.DoutS] = "SatopD"
	opnames[draw.SoutD | draw.DoutS] = "SxorD"
	opnames[draw.DinS | draw.DoutS] = "D"
	opnames[draw.DinS | draw.DoutS | draw.SoutD] = "DoverS"
	opnames[draw.DinS | draw.SoutD] = "DatopS"
	opnames[draw.DoutS | draw.SoutD] = "DxorS"
}

/* Extracts the SOP and prints it nicely. */
func extractSOP(a []byte) ([]byte, *DrawOp) {
	return a[2:], &DrawOp{DrawType{a[0]}, opnames[a[1]]}
}

func CreateDrawData(tag uint8, a []byte) (interface{}) {
	jd := &JsonedTwrdrawDraw{JsonedMsgCore{tag, "Twrdraw"}, []DrawCmd{}}
	
	for len(a) > 0 {
		t := a[0]
		log.Print("t: ", t, " len(a): ", len(a))
		oldalength := len(a)
		switch t {
		default:
			log.Fatal("lazy! unhandled message type")

		/* allocate screen: 'A' id[4] imageid[4] fillid[4] public[1] */
		case 'A':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+4+4+1:]

		/* allocate: 'b' id[4] screenid[4] refresh[1] chan[4] repl[1]
			R[4*4] clipR[4*4] rrggbbaa[4]
		 */
		case 'b':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+4+1+4+1+4*4+4*4+4:]

		/* set repl and clip: 'c' dstid[4] repl[1] clipR[4*4] */
		case 'c':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+1+4*4:]

		/* toggle debugging: 'D' val[1] */
		case 'D':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+1:]

		/* draw: 'd' dstid[4] srcid[4] maskid[4] R[4*4] P[2*4] P[2*4] */
		case 'd':
			jd.Cmds = append(jd.Cmds, &DrawBlit{
				DrawType{a[0]},
				binary.LittleEndian.Uint32(a[1:]),
				binary.LittleEndian.Uint32(a[5:]),
				binary.LittleEndian.Uint32(a[9:]),
				image.Rectangle{
					image.Point{
						int(binary.LittleEndian.Uint32(a[13:])),
						int(binary.LittleEndian.Uint32(a[17:]))},
					image.Point{
						int(binary.LittleEndian.Uint32(a[21:])),
						int(binary.LittleEndian.Uint32(a[25:]))}},
				image.Point{
					int(binary.LittleEndian.Uint32(a[29:])),
					int(binary.LittleEndian.Uint32(a[33:]))},
				image.Point{
					int(binary.LittleEndian.Uint32(a[37:])),
					int(binary.LittleEndian.Uint32(a[41:]))},
				})
			a = a[1 + 4 + 4 + 4 + 4*4 + 2*4 + 2*4:]

		/* ellipse: 'e' dstid[4] srcid[4] center[2*4] a[4] b[4] thick[4] sp[2*4] alpha[4] phi[4]*/
		case 'e', 'E':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+4+2*4+4+4+4+2*4+2*4:]

		/* free: 'f' id[4] */
		case 'f':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4:]

		/* free screen: 'F' id[4] */
		case 'F':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4:]

		/* initialize font: 'i' fontid[4] nchars[4] ascent[1] */
		case 'i':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+4+1:]

		/* set image 0 to screen image */
		case 'J':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1:]

		/* get image info: 'I' */
		case 'I':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1:]

		/* query: 'Q' n[1] queryspec[n] */
		case 'q':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+1+a[1]:]

		/* load character: 'l' fontid[4] srcid[4] index[2] R[4*4] P[2*4] left[1] width[1] */
		case 'l':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+4+2+4*4+2*4+1+1:]

		/* draw line: 'L' dstid[4] p0[2*4] p1[2*4] end0[4] end1[4] radius[4] srcid[4] sp[2*4] */
		case 'L':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+4+2+4*4+2*4+1+1:]

		/* attach to a named image: 'n' dstid[4] j[1] name[j] */
		case 'n':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+1:]

		/* name an image: 'N' dstid[4] in[1] j[1] name[j] */
		case 'N':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+1+1 + a[6]:]

		/* position window: 'o' id[4] r.min [2*4] screenr.min [2*4] */
		case 'o':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+2*4+2*4:]

		/* set compositing operator for next draw operation: 'O' op */
		case 'O':
			newa, op := extractSOP(a)
			a = newa
			jd.Cmds = append(jd.Cmds, op)

		/* filled polygon: 'P' dstid[4] n[2] wind[4] ignore[2*4] srcid[4] sp[2*4] p0[2*4] dp[2*2*n] */
		/* polygon: 'p' dstid[4] n[2] end0[4] end1[4] radius[4] srcid[4] sp[2*4] p0[2*4] dp[2*2*n] */
		case 'p', 'P':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[ 1+4+2+4+4+4+4+2*4:]

		/* read: 'r' id[4] R[4*4] */
		case 'r':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+4*4:]

		// Note: The way that I am drawing styled strings is not efficient.
		/* string: 's' dstid[4] srcid[4] fontid[4] P[2*4] clipr[4*4] sp[2*4] ni[2] ni*(index[2]) */
		/* stringbg: 'x' dstid[4] srcid[4] fontid[4] P[2*4] clipr[4*4] sp[2*4] ni[2] bgid[4] bgpt[2*4] ni*(index[2]) */
		case 's', 'x':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			m := 1+4+4+4+2*4+4*4+2*4+2
			if a[0] == 'x' {
				m += 4+2*4
			}
			ni := int(binary.LittleEndian.Uint16(a[45:]))
			m += ni * 2
		     a = a[m:]

		/* use public screen: 'S' id[4] chan[4] */
		case 'S':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+4:]

		/* top or bottom windows: 't' top[1] nw[2] n*id[4] */
		case 't':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+1+2:]

		/* visible: 'v' */
		case 'v':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1:]

		// TODO(rjkroege): Be more clever.
		/* write: 'y' id[4] R[4*4] data[x*1] */
		/* write from compressed data: 'Y' id[4] R[4*4] data[x*1] */
		case 'y', 'Y':
			jd.Cmds = append(jd.Cmds, &DrawType{ t })
			a = a[1+4+4*4:]
			a = a[len(a):]	// y uses up whole remaining buffer.
		}
		newalength := len(a)
		if newalength >= oldalength {
			log.Fatal("we are not decrementing?")
		}
	}
	return jd
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
		log.Print("Twrdraw: ", len(m.Data), m.Data)
		return CreateDrawData(m.Tag, m.Data)
		// return &JsonedMsgCore{m.Tag, "Twrdraw - expand data", }
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
