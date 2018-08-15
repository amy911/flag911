package flagger

type FlaggerCommon struct {
}

// Usage

func (fs *FlagSet) SetUsageFooter(text string) {
	fs.SetUsage(nil, nil, UsageWithHeaderAndFooter)
	UsageFooter = text
}

func (fs *FlagSet) SetUsageHeader(text string) {
	fs.SetUsage(nil, nil, UsageWithHeaderAndFooter)
	UsageHeader = text
}

// Feature tests

func (fs *FlagSet) CountWorks() bool {
	return false
}

func (fs *FlagSet) PVersionsWork() bool {
	return false
}

// Bool

func (fs *FlagSet) Bool(name string, value bool, usage string) *bool {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) BoolP(name string, shorthand string, value bool, usage string) *bool {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) BoolVarP(p *bool, name string, shorthand string, value bool, usage string) {
	panic(NewErrNotSupported())
}

// Count

func (fs *FlagSet) Count(name string, value int, usage string) *int {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) CountP(name string, shorthand string, value int, usage string) *int {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) CountVar(p *int, name string, value int, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) CountVarP(p *int, name string, shorthand string, value int, usage string) {
	panic(NewErrNotSupported())
}

// Float32

func (fs *FlagSet) Float32(name string, value float32, usage string) *float32 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Float32P(name string, shorthand string, value float32, usage string) *float32 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Float32Var(p *float32, name string, value float32, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Float32VarP(p *float32, name string, shorthand string, value float32, usage string) {
	panic(NewErrNotSupported())
}

// Float64

func (fs *FlagSet) Float64(name string, value float64, usage string) *float64 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Float64P(name string, shorthand string, value float64, usage string) *float64 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Float64VarP(p *float64, name string, shorthand string, value float64, usage string) {
	panic(NewErrNotSupported())
}

// Int

func (fs *FlagSet) Int(name string, value int, usage string) *int {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) IntP(name string, shorthand string, value int, usage string) *int {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) IntVar(p *int, name string, value int, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) IntVarP(p *int, name string, shorthand string, value int, usage string) {
	panic(NewErrNotSupported())
}

// Int8

func (fs *FlagSet) Int8(name string, value int8, usage string) *int8 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int8P(name string, shorthand string, value int8, usage string) *int8 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int8Var(p *int8, name string, value int8, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int8VarP(p *int8, name string, shorthand string, value int8, usage string) {
	panic(NewErrNotSupported())
}

// Int16

func (fs *FlagSet) Int16(name string, value int16, usage string) *int16 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int16P(name string, shorthand string, value int16, usage string) *int16 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int16Var(p *int16, name string, value int16, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int16VarP(p *int16, name string, shorthand string, value int16, usage string) {
	panic(NewErrNotSupported())
}

// Int32

func (fs *FlagSet) Int32(name string, value int32, usage string) *int32 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int32P(name string, shorthand string, value int32, usage string) *int32 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int32Var(p *int32, name string, value int32, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int32VarP(p *int32, name string, shorthand string, value int32, usage string) {
	panic(NewErrNotSupported())
}

// Int64

func (fs *FlagSet) Int64(name string, value int64, usage string) *int64 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int64P(name string, shorthand string, value int64, usage string) *int64 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Int64VarP(p *int64, name string, shorthand string, value int64, usage string) {
	panic(NewErrNotSupported())
}

// String

func (fs *FlagSet) String(name string, value string, usage string) *string {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) StringP(name string, shorthand string, value string, usage string) *string {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) StringVar(p *string, name string, value string, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) StringVarP(p *string, name string, shorthand string, value string, usage string) {
	panic(NewErrNotSupported())
}

// Uint

func (fs *FlagSet) Uint(name string, value uint, usage string) *uint {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) UintP(name string, shorthand string, value uint, usage string) *uint {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) UintVarP(p *uint, name string, shorthand string, value uint, usage string) {
	panic(NewErrNotSupported())
}

// Uint8

func (fs *FlagSet) Uint8(name string, value uint8, usage string) *uint8 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint8P(name string, shorthand string, value uint8, usage string) *uint8 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint8Var(p *uint8, name string, value uint8, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint8VarP(p *uint8, name string, shorthand string, value uint8, usage string) {
	panic(NewErrNotSupported())
}

// Uint16

func (fs *FlagSet) Uint16(name string, value uint16, usage string) *uint16 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint16P(name string, shorthand string, value uint16, usage string) *uint16 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint16Var(p *uint16, name string, value uint16, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint16VarP(p *uint16, name string, shorthand string, value uint16, usage string) {
	panic(NewErrNotSupported())
}

// Uint32

func (fs *FlagSet) Uint32(name string, value uint32, usage string) *uint32 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint32P(name string, shorthand string, value uint32, usage string) *uint32 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint32Var(p *uint32, name string, value uint32, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint32VarP(p *uint32, name string, shorthand string, value uint32, usage string) {
	panic(NewErrNotSupported())
}

// Uint64

func (fs *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint64P(name string, shorthand string, value uint64, usage string) *uint64 {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
	panic(NewErrNotSupported())
}

func (fs *FlagSet) Uint64VarP(p *uint64, name string, shorthand string, value uint64, usage string) {
	panic(NewErrNotSupported())
}
