package flagger

type NoCount struct {
}

func (fs NoCount) CountWorks() bool {
	return false
}

// Count

func (fs NoCount) Count(name string, value int, usage string) *int {
	panic(NewErrNotSupported())
}

func (fs NoCount) CountVar(p *int, name string, value int, usage string) {
	panic(NewErrNotSupported())
}
