package flagger

type NoInt16P struct {
}

// Int16

func (c *NoInt16P) Int16P(name string, shorthand string, value int16, usage string) *int16 {
	panic(NewErrNotSupported())
}

func (c *NoInt16P) Int16VarP(p *int16, name string, shorthand string, value int16, usage string) {
	panic(NewErrNotSupported())
}

// Uint16

func (c *NoInt16P) Uint16P(name string, shorthand string, value uint16, usage string) *uint16 {
	panic(NewErrNotSupported())
}

func (c *NoInt16P) Uint16VarP(p *uint16, name string, shorthand string, value uint16, usage string) {
	panic(NewErrNotSupported())
}
