package observing

import "context"

type ctxKey int

const (
	ctxTagCount ctxKey = iota
)

func TagCount(c context.Context) int {
	return c.Value(ctxTagCount).(int)
}
