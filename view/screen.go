package view

import (
	"errors"
	"fmt"
	"sync"

	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/channel"
)

type Screen struct {
	view *gocui.View
	g    *gocui.Gui

	msg string
}

func NewScreen(wg *sync.WaitGroup) *Screen {
	s := &Screen{}

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
		s.view = v
		s.g = g
	}

	return nil
}

func (s *Screen) PrintView(wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range channel.ScreenChan {
		if s.view == nil {
			continue
		}

		switch data.Type {

		}

		s.Display()

	}
}

func (s *Screen) Display() {
	s.g.UpdateAsync(func(g *gocui.Gui) error {
		s.view.Clear()
		fmt.Fprint(s.view, s.msg)
		return nil
	})
}
