package view

import (
	"errors"
	"fmt"
	"sync"

	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/channel"
)

type TextBox struct {
	view *gocui.View
	g    *gocui.Gui

	msg string
}

func NewTextBox(wg *sync.WaitGroup) *TextBox {
	t := &TextBox{}
	return t
}

func (t *TextBox) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("textbox", 0, (maxY-(maxY/15))-1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		v.Title = "textbox"
		v.Wrap = false
		v.Editable = true
		v.Editor = t
		t.view = v
		t.g = g

		_, err := g.SetCurrentView("textbox")
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (t *TextBox) PrintView(wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range channel.TextBoxChan {
		if t.view == nil {
			continue
		}

		switch data.Type {
		}

		t.Display()
	}
}

func (t *TextBox) Display() {
	t.g.UpdateAsync(func(g *gocui.Gui) error {
		t.view.Clear()
		fmt.Fprint(t.view, t.msg)
		return nil
	})
}

func (t *TextBox) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	// TODO: Ctrl-Backspace
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyTab:
		
	case key == gocui.KeyEnter:
		channel.MsgChan <- v.ViewBuffer() 
		v.SetCursor(0, 0)
		v.Clear()
	case key == gocui.KeyArrowDown:
	case key == gocui.KeyArrowUp:
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0)
	}
}


