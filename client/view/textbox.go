package view

import (
	"errors"
	"fmt"
	"sync"

	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/channel"
)

type TextBox struct {
	g    *gocui.Gui

	msg string
}

func NewTextBox(g *gocui.Gui, wg *sync.WaitGroup) *TextBox {
	t := &TextBox{g: g}
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
		switch data.Type {
		}

		t.Display(t.msg)
	}
}

func (t *TextBox) Display(msg string) {
	t.g.UpdateAsync(func(g *gocui.Gui) error {
		v, err := g.View("textbox")
		if err != nil {
			return err
		}
		v.Clear()
		fmt.Fprint(v, msg)
		return nil
	})
}

func (t *TextBox) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	// TODO: Ctrl-Backspace
        // TODO: Ctrl-Arrow_Keys
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
