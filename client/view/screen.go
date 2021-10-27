package view

import (
	"errors"
	"fmt"
	"sync"

	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/channel"
)

type Screen struct {
	g *gocui.Gui

	msg string
}

func NewScreen(g *gocui.Gui, wg *sync.WaitGroup) *Screen {
	s := &Screen{g: g}

	wg.Add(1)
	go s.PrintView(wg)

	return s
}

func (s *Screen) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("screen", 0, (maxY/15)+1, maxX-2, (maxY-(maxY/15))-2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "screen"
		v.Wrap = true
	}

	return nil
}

func (s *Screen) PrintView(wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range channel.ScreenChan {
		switch data.Type {
		case "msg":
			s.Display(data.Object.(string) + "\n")
		}

	}
}

func (s *Screen) Display(msg string) {
	s.g.UpdateAsync(func(g *gocui.Gui) error {
		v, err := g.View("screen")
		if err != nil {
			return err
		}

		line := v.ViewLinesHeight()
		_, cols := v.Size()
		if line > cols {
			ox, oy := v.Origin()
			v.SetOrigin(ox, oy+1)
		}
		fmt.Fprintf(v, "%s", msg)
		return nil
	})
}
