package view

import (
"errors"
	"fmt"
	"sync"
	"time"

	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/channel"
)

type Header struct {
	view *gocui.View
	g    *gocui.Gui

	msg string
}

func NewHeader(wg *sync.WaitGroup) *Header {
	h := &Header{}

	wg.Add(1)
	go h.PrintView(wg)

	return h
}

func (h *Header) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("header", 0, 0, maxX-1, (maxY / 15), 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "header"
		v.Wrap = false
		h.view = v
	}

	return nil
}

func (h *Header) PrintView(wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range channel.HeaderChan {
		if h.view == nil {
			continue
		}

		switch data.Type {
		case "clock":
			nowStr := time.Now().Format("02/01/2006 15:04:05")
			h.msg = nowStr + " | " + h.msg
		}
		h.Display()
	}
}

func (h *Header) Display() {
	h.g.UpdateAsync(func(g *gocui.Gui) error {
		h.view.Clear()
		fmt.Fprint(h.view, h.msg)
		return nil
	})
}
