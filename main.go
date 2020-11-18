package main

import (
	"image"
	"image/draw"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"

	"github.com/fogleman/nes/nes"
)

func main() {
	c, err := nes.NewConsole("roms/test.nes")
	if err != nil {
		log.Fatal(err)
	}
	driver.Main(func(s screen.Screen) {
		cbufb := c.Buffer().Bounds()
		buf, err := s.NewBuffer(image.Point{X: cbufb.Dx(), Y: cbufb.Dy()})
		if err != nil {
			log.Fatal(err)
		}
		defer buf.Release()

		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  "nes",
			Width:  buf.Size().X,
			Height: buf.Size().Y,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		var buttons [8]bool
	mainloop:
		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					break mainloop
				}
			case key.Event:
				var nextState bool
				if e.Direction == key.DirPress {
					nextState = true
				}
				switch e.Code {
				case key.CodeX:
					buttons[0] = nextState
				case key.CodeZ:
					buttons[1] = nextState
				case key.CodeS:
					buttons[2] = nextState
				case key.CodeA:
					buttons[3] = nextState
				case key.CodeUpArrow:
					buttons[4] = nextState
				case key.CodeDownArrow:
					buttons[5] = nextState
				case key.CodeLeftArrow:
					buttons[6] = nextState
				case key.CodeRightArrow:
					buttons[7] = nextState
				}
			case paint.Event:
				if !e.External {
					draw.Draw(buf.RGBA(), buf.Bounds(), c.Buffer(), image.ZP, draw.Over)
					w.Upload(image.ZP, buf, buf.Bounds())
					w.Publish()
				}
			case size.Event:
				// TODO

			case error:
				log.Println(e)
			}
			c.SetButtons1(buttons)
			c.StepFrame()
			w.Send(paint.Event{})
		}
	})
}
