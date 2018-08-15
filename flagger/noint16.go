package flagger

type NoInt16 struct {
}

// Int16

func (fs NoInt16) Int16(name string, value int16, usage string) *int16 {
	panic(NewErrNotSupported())
}

func (fs NoInt16) Int16Var(p *int16, name string, value int16, usage string) {
	panic(NewErrNotSupported())
}

// Uint16

func (fs NoInt16) Uint16(name string, value uint16, usage string) *uint16 {
	panic(NewErrNotSupported())
}

func (fs NoInt16) Uint16Var(p *uint16, name string, value uint16, usage string) {
	panic(NewErrNotSupported())
}
