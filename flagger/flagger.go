package flagger

type Flagger interface {
	Usage() interface{}
	SetUsage(func(interface{}) interface{}, interface{}, func())
	SetUsageFooter(string)
	SetUsageHeader(string)

	CountWorks() bool
	PVersionsWork() bool

	Bool(string, bool, string) *bool
	BoolP(string, string, bool, string) *bool
	BoolVar(*bool, string, bool, string)
	BoolVarP(*bool, string, string, bool, string)

	Count(string, int, string) *int
	CountP(string, string, int, string) *int
	CountVar(*int, string, int, string)
	CountVarP(*int, string, string, int, string)

	Float32(string, float32, string) *float32
	Float32P(string, string, float32, string) *float32
	Float32Var(*float32, string, float32, string)
	Float32VarP(*float32, string, string, float32, string)

	Float64(string, float64, string) *float64
	Float64P(string, string, float64, string) *float64
	Float64Var(*float64, string, float64, string)
	Float64VarP(*float64, string, string, float64, string)

	Int(string, int, string) *int
	IntP(string, string, int, string) *int
	IntVar(*int, string, int, string)
	IntVarP(*int, string, string, int, string)

	Int8(string, int8, string) *int8
	Int8P(string, string, int8, string) *int8
	Int8Var(*int8, string, int8, string)
	Int8VarP(*int8, string, string, int8, string)

	Int16(string, int16, string) *int16
	Int16P(string, string, int16, string) *int16
	Int16Var(*int16, string, int16, string)
	Int16VarP(*int16, string, string, int16, string)

	Int32(string, int32, string) *int32
	Int32P(string, string, int32, string) *int32
	Int32Var(*int32, string, int32, string)
	Int32VarP(*int32, string, string, int32, string)

	Int64(string, int64, string) *int64
	Int64P(string, string, int64, string) *int64
	Int64Var(*int64, string, int64, string)
	Int64VarP(*int64, string, string, int64, string)

	String(string, string, string) *string
	StringP(string, string, string, string) *string
	StringVar(*string, string, string, string)
	StringVarP(*string, string, string, string, string)

	Uint(string, uint, string) *uint
	UintP(string, string, uint, string) *uint
	UintVar(*uint, string, uint, string)
	UintVarP(*uint, string, string, uint, string)

	Uint8(string, uint8, string) *uint8
	Uint8P(string, string, uint8, string) *uint8
	Uint8Var(*uint8, string, uint8, string)
	Uint8VarP(*uint8, string, string, uint8, string)

	Uint16(string, uint16, string) *uint16
	Uint16P(string, string, uint16, string) *uint16
	Uint16Var(*uint16, string, uint16, string)
	Uint16VarP(*uint16, string, string, uint16, string)

	Uint32(string, uint32, string) *uint32
	Uint32P(string, string, uint32, string) *uint32
	Uint32Var(*uint32, string, uint32, string)
	Uint32VarP(*uint32, string, string, uint32, string)

	Uint64(string, uint64, string) *uint64
	Uint64P(string, string, uint64, string) *uint64
	Uint64Var(*uint64, string, uint64, string)
	Uint64VarP(*uint64, string, string, uint64, string)
}
