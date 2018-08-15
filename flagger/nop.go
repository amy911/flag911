package flagger

type NoP struct {
}

func (fs *NoP) PVersionsWork() bool {
	return false
}

// Bool

func (fs *NoP) BoolP(name string, shorthand string, value bool, usage string) *bool {
	panic(NewErrNotSupported())
}

func (fs *NoP) BoolVarP(p *bool, name string, shorthand string, value bool, usage string) {
	panic(NewErrNotSupported())
}

// Float32

func (fs *NoP) Float32P(name string, shorthand string, value float32, usage string) *float32 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Float32VarP(p *float32, name string, shorthand string, value float32, usage string) {
	panic(NewErrNotSupported())
}

// Float64

func (fs *NoP) Float64P(name string, shorthand string, value float64, usage string) *float64 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Float64VarP(p *float64, name string, shorthand string, value float64, usage string) {
	panic(NewErrNotSupported())
}

// Int

func (fs *NoP) IntP(name string, shorthand string, value int, usage string) *int {
	panic(NewErrNotSupported())
}

func (fs *NoP) IntVarP(p *int, name string, shorthand string, value int, usage string) {
	panic(NewErrNotSupported())
}

// Int8

func (fs *NoP) Int8P(name string, shorthand string, value int8, usage string) *int8 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Int8VarP(p *int8, name string, shorthand string, value int8, usage string) {
	panic(NewErrNotSupported())
}

// Int16

func (fs *NoP) Int16P(name string, shorthand string, value int16, usage string) *int16 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Int16VarP(p *int16, name string, shorthand string, value int16, usage string) {
	panic(NewErrNotSupported())
}

// Int32

func (fs *NoP) Int32P(name string, shorthand string, value int32, usage string) *int32 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Int32VarP(p *int32, name string, shorthand string, value int32, usage string) {
	panic(NewErrNotSupported())
}

// Int64

func (fs *NoP) Int64P(name string, shorthand string, value int64, usage string) *int64 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Int64VarP(p *int64, name string, shorthand string, value int64, usage string) {
	panic(NewErrNotSupported())
}

// String

func (fs *NoP) StringP(name string, shorthand string, value string, usage string) *string {
	panic(NewErrNotSupported())
}

func (fs *NoP) StringVarP(p *string, name string, shorthand string, value string, usage string) {
	panic(NewErrNotSupported())
}

// Uint

func (fs *NoP) UintP(name string, shorthand string, value uint, usage string) *uint {
	panic(NewErrNotSupported())
}

func (fs *NoP) UintVarP(p *uint, name string, shorthand string, value uint, usage string) {
	panic(NewErrNotSupported())
}

// Uint8

func (fs *NoP) Uint8P(name string, shorthand string, value uint8, usage string) *uint8 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Uint8VarP(p *uint8, name string, shorthand string, value uint8, usage string) {
	panic(NewErrNotSupported())
}

// Uint16

func (fs *NoP) Uint16P(name string, shorthand string, value uint16, usage string) *uint16 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Uint16VarP(p *uint16, name string, shorthand string, value uint16, usage string) {
	panic(NewErrNotSupported())
}

// Uint32

func (fs *NoP) Uint32P(name string, shorthand string, value uint32, usage string) *uint32 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Uint32VarP(p *uint32, name string, shorthand string, value uint32, usage string) {
	panic(NewErrNotSupported())
}

// Uint64

func (fs *NoP) Uint64P(name string, shorthand string, value uint64, usage string) *uint64 {
	panic(NewErrNotSupported())
}

func (fs *NoP) Uint64VarP(p *uint64, name string, shorthand string, value uint64, usage string) {
	panic(NewErrNotSupported())
}

