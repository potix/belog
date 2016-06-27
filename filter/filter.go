package filter

type Filter interface {
	Evaluate(log belog.Log) bool
}
