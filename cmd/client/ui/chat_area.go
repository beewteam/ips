package ui

import (
	"image"

	tui "github.com/marcusolsson/tui-go"
)

type ChatArea struct {
	chatBox    *tui.Box
	header     *tui.Entry
	messages   *tui.Box
	scrollArea *tui.ScrollArea
}

func NewChatArea() *ChatArea {
	header := tui.NewEntry()
	header.SetText("Chat:")

	messages := tui.NewVBox()
	msgArea := tui.NewScrollArea(messages)

	chat := tui.NewVBox(
		header,
		msgArea,
	)

	return &ChatArea{
		chatBox:    chat,
		header:     header,
		messages:   messages,
		scrollArea: msgArea,
	}
}

func (ca *ChatArea) AddNewMessage(msg string) {
	ca.scrollArea.Scroll(0, -1)
	entry := tui.NewEntry()
	entry.SetText(msg)
	ca.messages.Append(entry)
}

func (ca *ChatArea) MinSizeHint() image.Point {
	return ca.chatBox.MinSizeHint()
}

// SizeHint returns the size hint of the underlying widget.
func (ca *ChatArea) SizeHint() image.Point {
	return ca.chatBox.SizeHint()
}

// SizePolicy returns the default layout behavior.
func (ca *ChatArea) SizePolicy() (tui.SizePolicy, tui.SizePolicy) {
	return ca.chatBox.SizePolicy()
}

// Draw draws the scroll area.
func (ca *ChatArea) Draw(p *tui.Painter) {
	ca.chatBox.Draw(p)
}

// Resize resizes the scroll area and the underlying widget.
func (ca *ChatArea) Resize(size image.Point) {
	ca.chatBox.Resize(size)
}

func (ca *ChatArea) Size() image.Point {
	return ca.chatBox.Size()
}
