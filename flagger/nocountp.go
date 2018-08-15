package flagger

type NoCountP struct {
}

// Count

func (fs *NoCountP) CountP(name string, shorthand string, value int, usage string) *int {
	panic(NewErrNotSupported())
}

func (fs *NoCountP) CountVarP(p *int, name string, shorthand string, value int, usage string) {
	panic(NewErrNotSupported())
}
