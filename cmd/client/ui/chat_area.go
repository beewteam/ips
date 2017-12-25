package ui

import (
	tui "github.com/marcusolsson/tui-go"
)

type ChatArea struct {
	chatBox    *tui.Box
	header     *tui.Entry
	messages   *tui.Box
	scrollArea *tui.ScrollArea

	msgNr int
}

func NewChatArea() *ChatArea {
	header := tui.NewEntry()
	header.SetText("Chat:")

	messages := tui.NewVBox()
	scrollArea := tui.NewScrollArea(messages)

	chat := tui.NewVBox(
		header,
		scrollArea,
	)
	chat.SetSizePolicy(tui.Maximum, tui.Maximum)
	chat.SetBorder(true)

	return &ChatArea{
		chatBox:    chat,
		header:     header,
		messages:   messages,
		scrollArea: scrollArea,
	}
}

func (ca *ChatArea) AddNewMessage(msg string) {
	entry := tui.NewEntry()
	entry.SetText(msg)
	ca.messages.Append(entry)

	ca.msgNr++
	if ca.msgNr > ca.chatBox.Size().Y/ca.header.Size().Y {
		ca.scrollArea.Scroll(0, 1)
	}
}

func (ca *ChatArea) ToWidget() tui.Widget {
	return ca.chatBox
}
