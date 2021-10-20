package exhibiting

import (
	"strings"
)

type Layout int

const (
	Default Layout = iota
	Split
	Column
	Ultrawide
	Fullscreen
	Explode
)

var lmM = map[Layout]string{
	Default:    "960x1080+0+0 1920x1080+960+0 960x1080+2880+0",
	Split:      "1920x1080+0+0 1920x1080+1920+0",
	Column:     "1280x1080+0+0 1280x1080+1280+0 1280x1080+2560+0",
	Ultrawide:  "640x1080+0+0 2560x1080+640+0 640x1080+3200+0",
	Fullscreen: "3840x1080+0+0",
	Explode:    "1280x540+0+0 1280x540+1280+0 1280x540+2560+0 1280x540+0+540 1280x540+1280+540 1280x540+2560+540",
}

var slM = map[string]Layout{
	"default":    Default,
	"split":      Split,
	"column":     Column,
	"ultrawide":  Ultrawide,
	"fullscreen": Fullscreen,
	"explode":    Explode,
}

var lsM = map[Layout]string{
	Default:    "Default",
	Split:      "Split",
	Column:     "Column",
	Ultrawide:  "Ultrawide",
	Fullscreen: "Fullscreen",
	Explode:    "Explode",
}

func (l Layout) String() string {
	return lsM[l]
}

func (l Layout) ToMonitors() string {
	return lmM[l]
}

func layout(s string) Layout {
	return slM[strings.ToLower(s)]
}

func nextLayout(l Layout) Layout {
	switch l {
	case Explode:
		return Default
	}

	return l + 1
}
