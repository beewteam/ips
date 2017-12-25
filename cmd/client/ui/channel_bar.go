package ui

import (
	"github.com/marcusolsson/tui-go"
)

const (
	defaultWidth = 20
)

type ChannelBar struct {
	channelBox *tui.Box
	header     *tui.Entry
	channels   *tui.Box
	scrollArea *tui.ScrollArea

	channelsNr int
}

func NewChannelBar() *ChannelBar {
	// Gonna be changes as channels list can be modified by user
	header := tui.NewEntry()
	header.SetText("Channels:")

	channels := tui.NewVBox()
	scrollArea := tui.NewScrollArea(channels)

	channelBox := tui.NewVBox(
		header,
		channels,
	)
	header.MinSizeHint()
	channelBox.SetSizePolicy(tui.Maximum, tui.Maximum)
	channelBox.SetBorder(true)

	return &ChannelBar{
		channelBox: channelBox,
		header:     header,
		channels:   channels,
		scrollArea: scrollArea,
	}
}

func (cb *ChannelBar) AddChannel(channel string) {
	entry := tui.NewEntry()
	entry.SetText(channel[0:defaultWidth])
	cb.channels.Append(entry)

	cb.channelsNr++
	if cb.channelsNr > cb.channelBox.Size().Y/cb.header.Size().Y {
		cb.scrollArea.Scroll(0, 1)
	}
}

func (cb *ChannelBar) ToWidget() tui.Widget {
	return cb.channelBox
}
