
/*

	A devdraw listener. Maybe I've named this wrong.

*/

package main

import (
	"fmt"
	"log"
	"code.google.com/p/goplan9/draw"
	"image"
)

func  watcher() {
	
}

/*
 *	Redraws the world. (What world we have.)
 *	This is the "view" code. 
 */
func redraw(d *draw.Display, resized bool) {
	if resized {
		if err := d.Attach(draw.Refmesg); err != nil {
			log.Fatalf("can't reattach to window: %v", err)
		}
	}

	// draw coloured rects at mouse positions
	// first param is the clip rectangle. which can be 0. meaning no clip?
	var clipr image.Rectangle
	fmt.Printf("empty clip? %v\n", clipr)
	d.ScreenImage.Draw(clipr, d.White, nil, image.ZP)
	d.Flush(true)
}

/*
 *	Reads the mouse channel and does stuff. Like redrawing the screen.
 */
func mouse() {
	fmt.Printf("called mouse\n");
}

func main() {
	fmt.Print("hello from devdraw\n");

	// Make the window.	
	d, err := draw.Init(nil, "", "experiment1", "")
	if err != nil {
		log.Fatal(err)
	}

	
	// make some colors
	back, _ := d.AllocImage(image.Rect(0,0,1,1), d.ScreenImage.Pix, true, 0xDADBDAff);

	fmt.Printf("background colour: %v\n ", back);

	// get mouse positions
	mousectl := d.InitMouse()
	redraw(d, false);

	for {
		select {
		case <-mousectl.Resize:
			redraw(d, true)
		case m := <-mousectl.C:
			fmt.Printf("mouse field %v buttons %d\n", m, m.Buttons)
			// TODO(rjkroege): insert code here to do some drawing and stuff.
			d.ScreenImage.Draw(image.Rect(m.X, m.Y, m.X + 10, m.Y + 10), back, nil, image.ZP)
			d.Flush(true)
		}
	}

	fmt.Print("bye\n")
}
