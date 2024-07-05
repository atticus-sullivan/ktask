package main

import (
	"ktask/ktask"
	"ktask/ktask/parser"
	"time"

	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
)

// TextSerialiser is a specialised parser.Serialiser implementation for the terminal.
type TextSerialiser struct {
	DecimalDuration bool
	Styler          tf.Styler
}

func NewSerialiser(styler tf.Styler, decimal bool) TextSerialiser {
	return TextSerialiser{
		DecimalDuration: decimal,
		Styler:          styler,
	}
}

func (cs TextSerialiser) Date(d time.Time) string {
	return cs.Styler.Props(tf.StyleProps{Color: tf.TEXT, IsUnderlined: true}).Format(d.Format("2006-02-01"))
}

func (cs TextSerialiser) Stage(d ktask.Stage) string {
	return cs.Styler.Props(tf.StyleProps{Color: tf.PURPLE}).Format(string(d))
}

func (cs TextSerialiser) Name(s parser.NameText) string {
	txt := s.ToString()
	summaryStyler := cs.Styler.Props(tf.StyleProps{Color: tf.SUBDUED})
	txt = ktask.HashTagPattern.ReplaceAllStringFunc(txt, func(h string) string {
		return cs.Styler.Props(tf.StyleProps{Color: tf.SUBDUED, IsBold: true}).FormatAndRestore(
			h, summaryStyler,
		)
	})
	return summaryStyler.Format(txt)
}
