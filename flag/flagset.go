// import "github.com/amy911/flag911/flag"
package flag

import (
	"net"
	"sync"
	"time"
)

type Hook func(fs *FlagSet, fn string, p interface{}, name string, shorthand string, value interface{}, usage string, user []interface{}) error

type hookWrap struct {
	Hook
	User []interface{}
}

type FlagSet struct {
	Impl       interface{}
	User       []interface{}
	hooks      map[string]hookWrap
	hooksMutex sync.RWMutex
}

func New(impl interface{}, user ...interface{}) *FlagSet {
	return new(FlagSet).Init(impl, user...)
}

func (fs *FlagSet) Init(impl interface{}, user ...interface{}) *FlagSet {
	fs.Impl = impl
	fs.User = user
	fs.hooks = make(map[string]hookWrap)
	return fs
}

// Hooks

func (fs *FlagSet) DeleteAllHooks_() {
	fs.hooksMutex.Lock()
	defer fs.hooksMutex.Unlock()
	fs.hooks = make(map[string]hookWrap)
}

func (fs *FlagSet) DeleteHook(key string) (hookWrap, bool) {
	fs.hooksMutex.Lock()
	defer fs.hooksMutex.Unlock()
	hook, ok := fs.hooks[key]
	delete(fs.hooks, key)
	if !ok || hook.Hook == nil {
		hook.Hook = nil
		ok = false
	}
	return hook, ok
}

func (fs *FlagSet) GetHook(key string) (hookWrap, bool) {
	fs.hooksMutex.RLock()
	defer fs.hooksMutex.RUnlock()
	hook, ok := fs.hooks[key]
	if !ok || hook.Hook == nil {
		hook.Hook = nil
		ok = false
	}
	return hook, ok
}

func (fs *FlagSet) SetHook(key string, hook Hook, user ...interface{}) {
	fs.hooksMutex.Lock()
	defer fs.hooksMutex.Unlock()
	fs.hooks[key] = hookWrap{Hook: hook, User: user}
}

// Standard

func (fs *FlagSet) Parse(arguments []string) error {
	if f, ok := fs.Impl.(interface{ Parse([]string) error }); ok {
		return f.Parse(arguments)
	}
	panic(NewErrNotSupported("Parse"))
}

// Usage

func (fs *FlagSet) Usage() {
	if f, ok := fs.Impl.(interface{ Usage() }); ok {
		f.Usage()
		return
	}
	panic(NewErrNotSupported("Usage"))
}
func (fs *FlagSet) SetUsageHeader(text string) {
	if f, ok := fs.Impl.(interface{ SetUsageHeader(string) }); ok {
		f.SetUsageHeader(text)
		return
	}
	panic(NewErrNotSupported("SetUsageHeader"))
}
func (fs *FlagSet) SetUsageFooter(text string) {
	if f, ok := fs.Impl.(interface{ SetUsageFooter(string) }); ok {
		f.SetUsageFooter(text)
		return
	}
	panic(NewErrNotSupported("SetUsageFooter"))
}

// Feature tests

func (fs *FlagSet) CountWorks() bool {
	if f, ok := fs.Impl.(interface{ CountWorks() bool }); ok {
		return f.CountWorks()
	}
	return false
}
func (fs *FlagSet) ShorthandWorks() bool {
	if f, ok := fs.Impl.(interface{ PVersionsWork() bool }); ok {
		return f.PVersionsWork()
	}
	return false
}

// Bool

func (fs *FlagSet) Bool(name string, value bool, usage string) *bool {
	if fs.runHooks("Bool", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Bool(string, bool, string) *bool
	}); ok {
		return f.Bool(name, value, usage)
	}
	panic(NewErrNotSupported("Bool"))
}
func (fs *FlagSet) BoolP(name string, shorthand string, value bool, usage string) *bool {
	if fs.runHooks("BoolP", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		BoolP(string, string, bool, string) *bool
	}); ok {
		return f.BoolP(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Bool(string, bool, string) *bool
	}); ok {
		return f.Bool(name, value, usage)
	}
	panic(NewErrNotSupported("BoolP"))
}
func (fs *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	if fs.runHooks("BoolVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		BoolVar(*bool, string, bool, string)
	}); ok {
		f.BoolVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("BoolVar"))
}
func (fs *FlagSet) BoolVarP(p *bool, name string, shorthand string, value bool, usage string) {
	if fs.runHooks("BoolVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		BoolVarP(*bool, string, string, bool, string)
	}); ok {
		f.BoolVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		BoolVar(*bool, string, bool, string)
	}); ok {
		f.BoolVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("BoolVarP"))
}

// Count

func (fs *FlagSet) Count(name string, value int, usage string) *int {
	if fs.runHooks("Count", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Count(string, int, string) *int
	}); ok {
		return f.Count(name, value, usage)
	}
	panic(NewErrNotSupported("Count"))
}
func (fs *FlagSet) CountP(name string, shorthand string, value int, usage string) *int {
	if fs.runHooks("CountP", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		CountP(string, string, int, string) *int
	}); ok {
		return f.CountP(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Count(string, int, string) *int
	}); ok {
		return f.Count(name, value, usage)
	}
	panic(NewErrNotSupported("CountP"))
}
func (fs *FlagSet) CountVar(p *int, name string, value int, usage string) {
	if fs.runHooks("CountVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		CountVar(*int, string, int, string)
	}); ok {
		f.CountVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("CountVar"))
}
func (fs *FlagSet) CountVarP(p *int, name string, shorthand string, value int, usage string) {
	if fs.runHooks("CountVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		CountVarP(*int, string, string, int, string)
	}); ok {
		f.CountVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		CountVar(*int, string, int, string)
	}); ok {
		f.CountVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("CountVarP"))
}

// Duration

func (fs *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	if fs.runHooks("Duration", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Duration(string, time.Duration, string) *time.Duration
	}); ok {
		return f.Duration(name, value, usage)
	}
	panic(NewErrNotSupported("Duration"))
}
func (fs *FlagSet) DurationP(name string, shorthand string, value time.Duration, usage string) *time.Duration {
	if fs.runHooks("DurationP", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		DurationP(string, string, time.Duration, string) *time.Duration
	}); ok {
		return f.DurationP(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Duration(string, time.Duration, string) *time.Duration
	}); ok {
		return f.Duration(name, value, usage)
	}
	panic(NewErrNotSupported("DurationP"))
}
func (fs *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	if fs.runHooks("DurationVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		DurationVar(*time.Duration, string, time.Duration, string)
	}); ok {
		f.DurationVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("DurationVar"))
}
func (fs *FlagSet) DurationVarP(p *time.Duration, name string, shorthand string, value time.Duration, usage string) {
	if fs.runHooks("DurationVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		DurationVarP(*time.Duration, string, string, time.Duration, string)
	}); ok {
		f.DurationVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		DurationVar(*time.Duration, string, time.Duration, string)
	}); ok {
		f.DurationVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("DurationVarP"))
}

// Float32

func (fs *FlagSet) Float32(name string, value float32, usage string) *float32 {
	if fs.runHooks("Float32", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Float32(string, float32, string) *float32
	}); ok {
		return f.Float32(name, value, usage)
	}
	panic(NewErrNotSupported("Float32"))
}
func (fs *FlagSet) Float32P(name string, shorthand string, value float32, usage string) *float32 {
	if fs.runHooks("Float32P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Float32P(string, string, float32, string) *float32
	}); ok {
		return f.Float32P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Float32(string, float32, string) *float32
	}); ok {
		return f.Float32(name, value, usage)
	}
	panic(NewErrNotSupported("Float32P"))
}
func (fs *FlagSet) Float32Var(p *float32, name string, value float32, usage string) {
	if fs.runHooks("Float32Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Float32Var(*float32, string, float32, string)
	}); ok {
		f.Float32Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Float32Var"))
}
func (fs *FlagSet) Float32VarP(p *float32, name string, shorthand string, value float32, usage string) {
	if fs.runHooks("Float32VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Float32VarP(*float32, string, string, float32, string)
	}); ok {
		f.Float32VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Float32Var(*float32, string, float32, string)
	}); ok {
		f.Float32Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Float32VarP"))
}

// Float64

func (fs *FlagSet) Float64(name string, value float64, usage string) *float64 {
	if fs.runHooks("Float64", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Float64(string, float64, string) *float64
	}); ok {
		return f.Float64(name, value, usage)
	}
	panic(NewErrNotSupported("Float64"))
}
func (fs *FlagSet) Float64P(name string, shorthand string, value float64, usage string) *float64 {
	if fs.runHooks("Float64P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Float64P(string, string, float64, string) *float64
	}); ok {
		return f.Float64P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Float64(string, float64, string) *float64
	}); ok {
		return f.Float64(name, value, usage)
	}
	panic(NewErrNotSupported("Float64P"))
}
func (fs *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
	if fs.runHooks("Float64Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Float64Var(*float64, string, float64, string)
	}); ok {
		f.Float64Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Float64Var"))
}
func (fs *FlagSet) Float64VarP(p *float64, name string, shorthand string, value float64, usage string) {
	if fs.runHooks("Float64VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Float64VarP(*float64, string, string, float64, string)
	}); ok {
		f.Float64VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Float64Var(*float64, string, float64, string)
	}); ok {
		f.Float64Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Float64VarP"))
}

// Int

func (fs *FlagSet) Int(name string, value int, usage string) *int {
	if fs.runHooks("Int", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Int(string, int, string) *int
	}); ok {
		return f.Int(name, value, usage)
	}
	panic(NewErrNotSupported("Int"))
}
func (fs *FlagSet) IntP(name string, shorthand string, value int, usage string) *int {
	if fs.runHooks("IntP", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		IntP(string, string, int, string) *int
	}); ok {
		return f.IntP(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Int(string, int, string) *int
	}); ok {
		return f.Int(name, value, usage)
	}
	panic(NewErrNotSupported("IntP"))
}
func (fs *FlagSet) IntVar(p *int, name string, value int, usage string) {
	if fs.runHooks("IntVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		IntVar(*int, string, int, string)
	}); ok {
		f.IntVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("IntVar"))
}
func (fs *FlagSet) IntVarP(p *int, name string, shorthand string, value int, usage string) {
	if fs.runHooks("IntVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		IntVarP(*int, string, string, int, string)
	}); ok {
		f.IntVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		IntVar(*int, string, int, string)
	}); ok {
		f.IntVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("IntVarP"))
}

// Int8

func (fs *FlagSet) Int8(name string, value int8, usage string) *int8 {
	if fs.runHooks("Int8", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Int8(string, int8, string) *int8
	}); ok {
		return f.Int8(name, value, usage)
	}
	panic(NewErrNotSupported("Int8"))
}
func (fs *FlagSet) Int8P(name string, shorthand string, value int8, usage string) *int8 {
	if fs.runHooks("Int8P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Int8P(string, string, int8, string) *int8
	}); ok {
		return f.Int8P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Int8(string, int8, string) *int8
	}); ok {
		return f.Int8(name, value, usage)
	}
	panic(NewErrNotSupported("Int8P"))
}
func (fs *FlagSet) Int8Var(p *int8, name string, value int8, usage string) {
	if fs.runHooks("Int8Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Int8Var(*int8, string, int8, string)
	}); ok {
		f.Int8Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Int8Var"))
}
func (fs *FlagSet) Int8VarP(p *int8, name string, shorthand string, value int8, usage string) {
	if fs.runHooks("Int8VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Int8VarP(*int8, string, string, int8, string)
	}); ok {
		f.Int8VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Int8Var(*int8, string, int8, string)
	}); ok {
		f.Int8Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Int8VarP"))
}

// Int16

func (fs *FlagSet) Int16(name string, value int16, usage string) *int16 {
	if fs.runHooks("Int16", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Int16(string, int16, string) *int16
	}); ok {
		return f.Int16(name, value, usage)
	}
	panic(NewErrNotSupported("Int16"))
}
func (fs *FlagSet) Int16P(name string, shorthand string, value int16, usage string) *int16 {
	if fs.runHooks("Int16P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Int16P(string, string, int16, string) *int16
	}); ok {
		return f.Int16P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Int16(string, int16, string) *int16
	}); ok {
		return f.Int16(name, value, usage)
	}
	panic(NewErrNotSupported("Int16P"))
}
func (fs *FlagSet) Int16Var(p *int16, name string, value int16, usage string) {
	if fs.runHooks("Int16Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Int16Var(*int16, string, int16, string)
	}); ok {
		f.Int16Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Int16Var"))
}
func (fs *FlagSet) Int16VarP(p *int16, name string, shorthand string, value int16, usage string) {
	if fs.runHooks("Int16VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Int16VarP(*int16, string, string, int16, string)
	}); ok {
		f.Int16VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Int16Var(*int16, string, int16, string)
	}); ok {
		f.Int16Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Int16VarP"))
}

// Int32

func (fs *FlagSet) Int32(name string, value int32, usage string) *int32 {
	if fs.runHooks("Int32", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Int32(string, int32, string) *int32
	}); ok {
		return f.Int32(name, value, usage)
	}
	panic(NewErrNotSupported("Int32"))
}
func (fs *FlagSet) Int32P(name string, shorthand string, value int32, usage string) *int32 {
	if fs.runHooks("Int32P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Int32P(string, string, int32, string) *int32
	}); ok {
		return f.Int32P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Int32(string, int32, string) *int32
	}); ok {
		return f.Int32(name, value, usage)
	}
	panic(NewErrNotSupported("Int32P"))
}
func (fs *FlagSet) Int32Var(p *int32, name string, value int32, usage string) {
	if fs.runHooks("Int32Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Int32Var(*int32, string, int32, string)
	}); ok {
		f.Int32Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Int32Var"))
}
func (fs *FlagSet) Int32VarP(p *int32, name string, shorthand string, value int32, usage string) {
	if fs.runHooks("Int32VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Int32VarP(*int32, string, string, int32, string)
	}); ok {
		f.Int32VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Int32Var(*int32, string, int32, string)
	}); ok {
		f.Int32Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Int32VarP"))
}

// Int64

func (fs *FlagSet) Int64(name string, value int64, usage string) *int64 {
	if fs.runHooks("Int64", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Int64(string, int64, string) *int64
	}); ok {
		return f.Int64(name, value, usage)
	}
	panic(NewErrNotSupported("Int64"))
}
func (fs *FlagSet) Int64P(name string, shorthand string, value int64, usage string) *int64 {
	if fs.runHooks("Int64P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Int64P(string, string, int64, string) *int64
	}); ok {
		return f.Int64P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Int64(string, int64, string) *int64
	}); ok {
		return f.Int64(name, value, usage)
	}
	panic(NewErrNotSupported("Int64P"))
}
func (fs *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	if fs.runHooks("Int64Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Int64Var(*int64, string, int64, string)
	}); ok {
		f.Int64Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Int64Var"))
}
func (fs *FlagSet) Int64VarP(p *int64, name string, shorthand string, value int64, usage string) {
	if fs.runHooks("Int64VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Int64VarP(*int64, string, string, int64, string)
	}); ok {
		f.Int64VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Int64Var(*int64, string, int64, string)
	}); ok {
		f.Int64Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Int64VarP"))
}

// IP

func (fs *FlagSet) IP(name string, value net.IP, usage string) *net.IP {
	if fs.runHooks("IP", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		IP(string, net.IP, string) *net.IP
	}); ok {
		return f.IP(name, value, usage)
	}
	panic(NewErrNotSupported("IP"))
}
func (fs *FlagSet) IPP(name string, shorthand string, value net.IP, usage string) *net.IP {
	if fs.runHooks("IPP", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		IPP(string, string, net.IP, string) *net.IP
	}); ok {
		return f.IPP(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		IP(string, net.IP, string) *net.IP
	}); ok {
		return f.IP(name, value, usage)
	}
	panic(NewErrNotSupported("IPP"))
}
func (fs *FlagSet) IPVar(p *net.IP, name string, value net.IP, usage string) {
	if fs.runHooks("IPVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		IPVar(*net.IP, string, net.IP, string)
	}); ok {
		f.IPVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("IPVar"))
}
func (fs *FlagSet) IPVarP(p *net.IP, name string, shorthand string, value net.IP, usage string) {
	if fs.runHooks("IPVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		IPVarP(*net.IP, string, string, net.IP, string)
	}); ok {
		f.IPVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		IPVar(*net.IP, string, net.IP, string)
	}); ok {
		f.IPVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("IPVarP"))
}

// IPMask

func (fs *FlagSet) IPMask(name string, value net.IPMask, usage string) *net.IPMask {
	if fs.runHooks("IPMask", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		IPMask(string, net.IPMask, string) *net.IPMask
	}); ok {
		return f.IPMask(name, value, usage)
	}
	panic(NewErrNotSupported("IPMask"))
}
func (fs *FlagSet) IPMaskP(name string, shorthand string, value net.IPMask, usage string) *net.IPMask {
	if fs.runHooks("IPMaskP", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		IPMaskP(string, string, net.IPMask, string) *net.IPMask
	}); ok {
		return f.IPMaskP(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		IPMask(string, net.IPMask, string) *net.IPMask
	}); ok {
		return f.IPMask(name, value, usage)
	}
	panic(NewErrNotSupported("IPMaskP"))
}
func (fs *FlagSet) IPMaskVar(p *net.IPMask, name string, value net.IPMask, usage string) {
	if fs.runHooks("IPMaskVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		IPMaskVar(*net.IPMask, string, net.IPMask, string)
	}); ok {
		f.IPMaskVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("IPMaskVar"))
}
func (fs *FlagSet) IPMaskVarP(p *net.IPMask, name string, shorthand string, value net.IPMask, usage string) {
	if fs.runHooks("IPMaskVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		IPMaskVarP(*net.IPMask, string, string, net.IPMask, string)
	}); ok {
		f.IPMaskVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		IPMaskVar(*net.IPMask, string, net.IPMask, string)
	}); ok {
		f.IPMaskVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("IPMaskVarP"))
}

// String

func (fs *FlagSet) String(name string, value string, usage string) *string {
	if fs.runHooks("String", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		String(string, string, string) *string
	}); ok {
		return f.String(name, value, usage)
	}
	panic(NewErrNotSupported("String"))
}
func (fs *FlagSet) StringP(name string, shorthand string, value string, usage string) *string {
	if fs.runHooks("StringP", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		StringP(string, string, string, string) *string
	}); ok {
		return f.StringP(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		String(string, string, string) *string
	}); ok {
		return f.String(name, value, usage)
	}
	panic(NewErrNotSupported("StringP"))
}
func (fs *FlagSet) StringVar(p *string, name string, value string, usage string) {
	if fs.runHooks("StringVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		StringVar(*string, string, string, string)
	}); ok {
		f.StringVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("StringVar"))
}
func (fs *FlagSet) StringVarP(p *string, name string, shorthand string, value string, usage string) {
	if fs.runHooks("StringVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		StringVarP(*string, string, string, string, string)
	}); ok {
		f.StringVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		StringVar(*string, string, string, string)
	}); ok {
		f.StringVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("StringVarP"))
}

// Time

func (fs *FlagSet) Time(name string, value time.Time, usage string) *time.Time {
	if fs.runHooks("Time", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Time(string, time.Time, string) *time.Time
	}); ok {
		return f.Time(name, value, usage)
	}
	panic(NewErrNotSupported("Time"))
}
func (fs *FlagSet) TimeP(name string, shorthand string, value time.Time, usage string) *time.Time {
	if fs.runHooks("TimeP", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		TimeP(string, string, time.Time, string) *time.Time
	}); ok {
		return f.TimeP(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Time(string, time.Time, string) *time.Time
	}); ok {
		return f.Time(name, value, usage)
	}
	panic(NewErrNotSupported("TimeP"))
}
func (fs *FlagSet) TimeVar(p *time.Time, name string, value time.Time, usage string) {
	if fs.runHooks("TimeVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		TimeVar(*time.Time, string, time.Time, string)
	}); ok {
		f.TimeVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("TimeVar"))
}
func (fs *FlagSet) TimeVarP(p *time.Time, name string, shorthand string, value time.Time, usage string) {
	if fs.runHooks("TimeVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		TimeVarP(*time.Time, string, string, time.Time, string)
	}); ok {
		f.TimeVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		TimeVar(*time.Time, string, time.Time, string)
	}); ok {
		f.TimeVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("TimeVarP"))
}

// Uint

func (fs *FlagSet) Uint(name string, value uint, usage string) *uint {
	if fs.runHooks("Uint", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Uint(string, uint, string) *uint
	}); ok {
		return f.Uint(name, value, usage)
	}
	panic(NewErrNotSupported("Uint"))
}
func (fs *FlagSet) UintP(name string, shorthand string, value uint, usage string) *uint {
	if fs.runHooks("UintP", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		UintP(string, string, uint, string) *uint
	}); ok {
		return f.UintP(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Uint(string, uint, string) *uint
	}); ok {
		return f.Uint(name, value, usage)
	}
	panic(NewErrNotSupported("UintP"))
}
func (fs *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
	if fs.runHooks("UintVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		UintVar(*uint, string, uint, string)
	}); ok {
		f.UintVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("UintVar"))
}
func (fs *FlagSet) UintVarP(p *uint, name string, shorthand string, value uint, usage string) {
	if fs.runHooks("UintVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		UintVarP(*uint, string, string, uint, string)
	}); ok {
		f.UintVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		UintVar(*uint, string, uint, string)
	}); ok {
		f.UintVar(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("UintVarP"))
}

// Uint8

func (fs *FlagSet) Uint8(name string, value uint8, usage string) *uint8 {
	if fs.runHooks("Uint8", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Uint8(string, uint8, string) *uint8
	}); ok {
		return f.Uint8(name, value, usage)
	}
	panic(NewErrNotSupported("Uint8"))
}
func (fs *FlagSet) Uint8P(name string, shorthand string, value uint8, usage string) *uint8 {
	if fs.runHooks("Uint8P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Uint8P(string, string, uint8, string) *uint8
	}); ok {
		return f.Uint8P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Uint8(string, uint8, string) *uint8
	}); ok {
		return f.Uint8(name, value, usage)
	}
	panic(NewErrNotSupported("Uint8P"))
}
func (fs *FlagSet) Uint8Var(p *uint8, name string, value uint8, usage string) {
	if fs.runHooks("Uint8Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Uint8Var(*uint8, string, uint8, string)
	}); ok {
		f.Uint8Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Uint8Var"))
}
func (fs *FlagSet) Uint8VarP(p *uint8, name string, shorthand string, value uint8, usage string) {
	if fs.runHooks("Uint8VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Uint8VarP(*uint8, string, string, uint8, string)
	}); ok {
		f.Uint8VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Uint8Var(*uint8, string, uint8, string)
	}); ok {
		f.Uint8Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Uint8VarP"))
}

// Uint16

func (fs *FlagSet) Uint16(name string, value uint16, usage string) *uint16 {
	if fs.runHooks("Uint16", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Uint16(string, uint16, string) *uint16
	}); ok {
		return f.Uint16(name, value, usage)
	}
	panic(NewErrNotSupported("Uint16"))
}
func (fs *FlagSet) Uint16P(name string, shorthand string, value uint16, usage string) *uint16 {
	if fs.runHooks("Uint16P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Uint16P(string, string, uint16, string) *uint16
	}); ok {
		return f.Uint16P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Uint16(string, uint16, string) *uint16
	}); ok {
		return f.Uint16(name, value, usage)
	}
	panic(NewErrNotSupported("Uint16P"))
}
func (fs *FlagSet) Uint16Var(p *uint16, name string, value uint16, usage string) {
	if fs.runHooks("Uint16Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Uint16Var(*uint16, string, uint16, string)
	}); ok {
		f.Uint16Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Uint16Var"))
}
func (fs *FlagSet) Uint16VarP(p *uint16, name string, shorthand string, value uint16, usage string) {
	if fs.runHooks("Uint16VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Uint16VarP(*uint16, string, string, uint16, string)
	}); ok {
		f.Uint16VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Uint16Var(*uint16, string, uint16, string)
	}); ok {
		f.Uint16Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Uint16VarP"))
}

// Uint32

func (fs *FlagSet) Uint32(name string, value uint32, usage string) *uint32 {
	if fs.runHooks("Uint32", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Uint32(string, uint32, string) *uint32
	}); ok {
		return f.Uint32(name, value, usage)
	}
	panic(NewErrNotSupported("Uint32"))
}
func (fs *FlagSet) Uint32P(name string, shorthand string, value uint32, usage string) *uint32 {
	if fs.runHooks("Uint32P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Uint32P(string, string, uint32, string) *uint32
	}); ok {
		return f.Uint32P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Uint32(string, uint32, string) *uint32
	}); ok {
		return f.Uint32(name, value, usage)
	}
	panic(NewErrNotSupported("Uint32P"))
}
func (fs *FlagSet) Uint32Var(p *uint32, name string, value uint32, usage string) {
	if fs.runHooks("Uint32Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Uint32Var(*uint32, string, uint32, string)
	}); ok {
		f.Uint32Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Uint32Var"))
}
func (fs *FlagSet) Uint32VarP(p *uint32, name string, shorthand string, value uint32, usage string) {
	if fs.runHooks("Uint32VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Uint32VarP(*uint32, string, string, uint32, string)
	}); ok {
		f.Uint32VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Uint32Var(*uint32, string, uint32, string)
	}); ok {
		f.Uint32Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Uint32VarP"))
}

// Uint64

func (fs *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	if fs.runHooks("Uint64", nil, name, "", value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Uint64(string, uint64, string) *uint64
	}); ok {
		return f.Uint64(name, value, usage)
	}
	panic(NewErrNotSupported("Uint64"))
}
func (fs *FlagSet) Uint64P(name string, shorthand string, value uint64, usage string) *uint64 {
	if fs.runHooks("Uint64P", nil, name, shorthand, value, usage) {
		panic(Cancel) // This operation cannot be safely canceled
	}
	if f, ok := fs.Impl.(interface {
		Uint64P(string, string, uint64, string) *uint64
	}); ok {
		return f.Uint64P(name, shorthand, value, usage)
	} else if f, ok := fs.Impl.(interface {
		Uint64(string, uint64, string) *uint64
	}); ok {
		return f.Uint64(name, value, usage)
	}
	panic(NewErrNotSupported("Uint64P"))
}
func (fs *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
	if fs.runHooks("Uint64Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Uint64Var(*uint64, string, uint64, string)
	}); ok {
		f.Uint64Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Uint64Var"))
}
func (fs *FlagSet) Uint64VarP(p *uint64, name string, shorthand string, value uint64, usage string) {
	if fs.runHooks("Uint64VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.Impl.(interface {
		Uint64VarP(*uint64, string, string, uint64, string)
	}); ok {
		f.Uint64VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.Impl.(interface {
		Uint64Var(*uint64, string, uint64, string)
	}); ok {
		f.Uint64Var(p, name, value, usage)
		return
	}
	panic(NewErrNotSupported("Uint64VarP"))
}

// Internal

func (fs *FlagSet) runHooks(fn string, p interface{}, name string, shorthand string, value interface{}, usage string) (cancel bool) {
	fs.hooksMutex.RLock()
	defer fs.hooksMutex.RUnlock()
	for _, hook := range fs.hooks {
		if hook.Hook != nil {
			if err := hook.Hook(fs, fn, p, name, shorthand, value, usage, hook.User); err != nil {
				if _, ok := err.(ErrCancel); ok {
					cancel = true
				} else {
					panic(err)
				}
			}
		}
	}
	return
}
