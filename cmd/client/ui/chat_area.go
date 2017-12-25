package ui

import (
	tui "github.com/marcusolsson/tui-go"
)

type ChatArea struct {
	chatBox    *tui.Box
	messages   *tui.Box
	scrollArea *tui.ScrollArea

	msgNr int
}

func NewChatArea() *ChatArea {
	messages := tui.NewVBox()
	scrollArea := tui.NewScrollArea(messages)
	messages.SetSizePolicy(tui.Expanding, tui.Minimum)

	chat := tui.NewVBox(
		scrollArea,
	)
	chat.SetSizePolicy(tui.Expanding, tui.Minimum)
	chat.SetBorder(true)

	return &ChatArea{
		chatBox:    chat,
		messages:   messages,
		scrollArea: scrollArea,
	}
}

func (ca *ChatArea) AddNewMessage(msg string) {
	entry := tui.NewLabel(msg)
	ca.messages.Append(entry)

	ca.msgNr++
	if ca.msgNr > ca.chatBox.Size().Y-10 {
		ca.scrollArea.Scroll(0, 1)
	}
}

func (ca *ChatArea) ToWidget() tui.Widget {
	return ca.chatBox
}
