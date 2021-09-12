package exhibiting

import (
	"strings"
)

type Layout int

const (
	Default Layout = iota
	Ultrawide
	Split
)

var lmM = map[Layout]string{
	Default:   "960x1080+0+0 1920x1080+960+0 960x1080+2880+0",
	Ultrawide: "640x1080+0+0 2560x1080+640+0 640x1080+3200+0",
	Split:     "1920x1080+0+0 1920x1080+1920+0",
}

var slM = map[string]Layout{
	"default":   Default,
	"ultrawide": Ultrawide,
	"split":     Split,
}

var lsM = map[Layout]string{
	Default:   "Default",
	Ultrawide: "Ultrawide",
	Split:     "Split",
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
	case Default,
		Ultrawide:
		return l + 1
	}

	return Default
}
