package observing

type Layout int

const (
	Default Layout = iota
	Ultrawide
	Split
	Explode
)

var mlM = map[string]Layout{
	"960x1080+0+0 1920x1080+960+0 960x1080+2880+0":                                                    Default,
	"640x1080+0+0 2560x1080+640+0 640x1080+3200+0":                                                    Ultrawide,
	"1920x1080+0+0 1920x1080+1920+0":                                                                  Split,
	"1280x540+0+0 1280x540+1280+0 1280x540+2560+0 1280x540+0+540 1280x540+1280+540 1280x540+2560+540": Explode,
}

var lsM = map[Layout]string{
	Default:   "Default",
	Ultrawide: "Ultrawide",
	Split:     "Split",
	Explode:   "Explode",
}

func (l Layout) String() string {
	return lsM[l]
}

func layout(s string) Layout {
	return mlM[s]
}
