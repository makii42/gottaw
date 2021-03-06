package output

import (
	"fmt"

	n "github.com/0xAX/notificator"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	c "github.com/makii42/gottaw/config"
)

type Level int

var (
	Trace, Quiet      bool
	spinnerWorkChars  = spinner.CharSets[11]
	spinnerWorkSuffix = "   executing pipeline"
	spinnerWaitChars  = spinner.CharSets[38]
	spinnerWaitSuffix = "   waiting for changes"
)

const (
	L_QUIET Level = iota
	L_NOTICE
	L_TRACE
)

type (
	log struct {
		cfg                                        *c.Config
		errors, notices, triggers, success, normal *color.Color
		n                                          *n.Notificator
		level                                      Level
		spin                                       *spinner.Spinner
	}

	// Logger is the main interface for outputing things.
	Logger interface {
		Errorf(format string, a ...interface{})
		Noticef(format string, a ...interface{})
		Triggerf(format string, a ...interface{})
		Successf(format string, a ...interface{})
		Tracef(format string, a ...interface{})
	}
)

func (o *log) growl(title, msg, icon, urgency string) {
	if o.n != nil {
		o.n.Push(title, msg, icon, urgency)
	}
}

// NewLog creates a new logger that handles output.
func NewLog(cfg *c.Config) (Logger, error) {
	var lvl Level

	if Trace && Quiet {
		return nil, fmt.Errorf("please decide whether you want me to trace or be quiet ;)")
	} else if Trace {
		lvl = L_TRACE
	} else if Quiet {
		lvl = L_QUIET
	} else {
		lvl = L_NOTICE
	}

	var notificator *n.Notificator
	if cfg.Growl {
		if notificator == nil {
			notificator = makeNotificator()
		}
	}

	l := log{
		cfg:      cfg,
		errors:   color.New(color.FgHiRed),
		notices:  color.New(color.FgBlue),
		triggers: color.New(color.FgYellow),
		success:  color.New(color.FgGreen),
		normal:   color.New(color.FgHiWhite),
		n:        notificator,
		level:    lvl,
	}
	return &l, nil
}

func (l *log) GetLog() Logger {
	return l
}

/*

		spin = spinner.New(spinnerWorkChars, 200*time.Millisecond)
		spin.Suffix = spinnerWorkSuffix
		spin.Start()

				spin.UpdateCharSet(spinnerWaitChars)
		spin.Suffix = spinnerWaitSuffix
		spin.Restart()

	spin.UpdateCharSet(spinnerWorkChars)
	spin.Suffix = spinnerWorkSuffix
	spin.Restart()

	spin.UpdateCharSet(spinnerWaitChars)
	spin.Suffix = spinnerWaitSuffix
	spin.Restart()

*/

func (l *log) Errorf(format string, a ...interface{}) {
	l.errors.Printf(format, a...)
	l.growl("Error", fmt.Sprintf(format, a...), "", n.UR_CRITICAL)
}

func (l *log) Noticef(format string, a ...interface{}) {
	l.notices.Printf(format, a...)
}

func (l *log) Triggerf(format string, a ...interface{}) {
	l.triggers.Printf(format, a...)
}

func (l *log) Successf(format string, a ...interface{}) {
	l.success.Printf(format, a...)
	l.growl("Pipeline Success", fmt.Sprintf(format, a...), "", n.UR_NORMAL)
}

func (l *log) Tracef(format string, a ...interface{}) {
	if l.level >= L_TRACE {
		l.normal.Printf(format, a...)
	}
}

func makeNotificator() *n.Notificator {
	return n.New(n.Options{
		AppName: "Gotta Watch!",
	})
}
