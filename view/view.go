package view

import (
	"sync"

	"github.com/awesome-gocui/gocui"
	"gitlab.com/smallwood/sw-chat/channel"
)

var instance *View
var once sync.Once

func Handler() *View {
	once.Do(func() {
		instance = new()
	})

	return instance
}

type View struct {
	list        []string
	currentView int
	wg          sync.WaitGroup
}

func new() *View {
	v := &View{}
	return v
}

func (v *View) Layout(g *gocui.Gui) error {
	header := NewHeader(&v.wg)
	if err := header.Layout(g); err != nil {
		panic(err)
	}

	screen := NewScreen(&v.wg)
	if err := screen.Layout(g); err != nil {
		panic(err)
	}

	textbox := NewTextBox(&v.wg)
	if err := textbox.Layout(g); err != nil {
		panic(err)
	}

	return nil
}

func (v *View) SetupViews(views []string) {
	v.list = views
}

func (v *View) Wait() {
	v.wg.Wait()
}

func (v *View) SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}

	return g.SetViewOnTop(name)
}

func (v *View) NextView(g *gocui.Gui, view *gocui.View) error {
	nextIndex := (v.currentView + 1) % (len(v.list) - 1) // note: ignore header view

	name := v.list[nextIndex]

	if _, err := v.SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	v.currentView = nextIndex

	return nil
}

func (v *View) Quit(g *gocui.Gui, view *gocui.View) error {
	select {
	case channel.GlobalShutdown <- true:
	default:
	}

	close(channel.HeaderChan)
	close(channel.ScreenChan)
	close(channel.TextBoxChan)

	return gocui.ErrQuit
}
