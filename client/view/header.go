package view

import (
	"errors"
	"fmt"
	"sync"

	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/channel"
)

type Header struct {
	g *gocui.Gui

	msg  string
	cmsg string
}

func NewHeader(g *gocui.Gui, wg *sync.WaitGroup) *Header {
	h := &Header{g: g}

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
	}

	return nil
}

func (h *Header) PrintView(wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range channel.HeaderChan {
		switch data.Type {
		case "clock":
			h.cmsg = data.Object.(string) + "|"
		case "msg":
			h.msg = "[" + data.Object.(string) + "]"
		}
		h.Display(h.cmsg + h.msg)
	}
}

func (h *Header) Display(msg string) {
	h.g.UpdateAsync(func(g *gocui.Gui) error {
		v, err := g.View("header")
		if err != nil {
			return err
		}

		v.Clear()
		fmt.Fprint(v, msg)
		return nil
	})
}
