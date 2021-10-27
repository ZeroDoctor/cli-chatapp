package channel

var (
	GlobalShutdown = make(chan bool)

	HeaderChan = make(chan Data, 4)

	ScreenChan = make(chan Data, 10)

	TextBoxChan = make(chan Data, 4)

	MsgChan = make(chan string, 4)
)

type Data struct {
	Type   string
	Object interface{}
}
