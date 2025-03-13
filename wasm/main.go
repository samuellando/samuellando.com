//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/htmlcanvas"
)

const FPS = 60

func main() {
	fmt.Println("Hello WASM")
	cvs := js.Global().Get("document").Call("getElementById", "canvas")
	w := float64(1000)
	h := float64(500)
	c := htmlcanvas.New(cvs, w, h, 1)
	ctx := canvas.NewContext(c)

	rect := canvas.Rectangle(w, h)

	frame_delay := time.Second / FPS
	planets := []*planet{
		ceratePlanet(200, 200, 100, 0.50, 10, canvas.Black),
		ceratePlanet(50, 50, 50, 1, 5, canvas.Blue),
	}
	ctx.SetFill(canvas.White)
	for {
		start := time.Now()
		for _, p := range planets {
			p.draw(ctx)
		}
		sleep := frame_delay - time.Since(start)
		if sleep > 0 {
			time.Sleep(frame_delay)
		} else {
			fmt.Println("Dropped frames")
		}
		ctx.DrawPath(0, 0, rect)
	}
}
