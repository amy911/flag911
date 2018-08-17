package flag

import (
	golang_flag "flag"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/suite911/error911"
	"github.com/suite911/term911/vt"

	ogier_pflag "github.com/ogier/pflag"
	spf13_pflag "github.com/spf13/pflag"
)

type Hook func(fs *FlagSet, fn string, p interface{}, name string, shorthand string, value interface{}, usage string, user []interface{}) error

type hookWrap struct {
	Hook
	User []interface{}
}

type FlagSet struct {
	name       string
	impl       interface{}
	usage      interface{}
	hooks      map[string]hookWrap
	hooksMutex sync.RWMutex

	usageHeader, usageFooter string

	output io.Writer
}

func New(name string, impl interface{}) *FlagSet {
	return new(FlagSet).Init(name, impl)
}

func (fs *FlagSet) Init(name string, impl interface{}) *FlagSet {
	fs.name = name
	fs.impl = impl
	fs.hooks = make(map[string]hookWrap)
	fs.output = os.Stderr
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

// Special

func IsZeroValue(value string) bool {
	return len(value) < 1 || value == "0" || value == "false"
}

func (fs *FlagSet) PrettyPrintDefaults() {
	switch fs.impl.(type) {
	case *golang_flag.FlagSet:
		w := fs.output
		fs.SetOutput(os.Stderr)
		fs.impl.(*golang_flag.FlagSet).VisitAll(prettyPrintDefault_golang_flag)
		fs.SetOutput(w)
	case *ogier_pflag.FlagSet:
		w := fs.output
		fs.SetOutput(os.Stderr)
		fs.impl.(*ogier_pflag.FlagSet).VisitAll(prettyPrintDefault_ogier_pflag)
		fs.SetOutput(w)
	case *spf13_pflag.FlagSet:
		w := fs.output
		fs.SetOutput(os.Stderr)
		fs.impl.(*spf13_pflag.FlagSet).VisitAll(prettyPrintDefault_spf13_pflag)
		fs.SetOutput(w)
	default:
		fs.PrintDefaults()
	}
}

func prettyPrintDefault_golang_flag(flag *golang_flag.Flag) {
	typeName, usage := golang_flag.UnquoteUsage(flag)
	fmt.Fprintf(os.Stderr, prettyPrintDefault(flag.Name, flag.Shorthand, typeName, usage, flag.DefValue, flag.Value), vt.SafeNewLine)
}

func prettyPrintDefault_ogier_pflag(flag *ogier_pflag.Flag) {
	typeName, usage := ogier_pflag.UnquoteUsage(flag)
	fmt.Fprintf(os.Stderr, prettyPrintDefault(flag.Name, flag.Shorthand, typeName, usage, flag.DefValue, flag.Value), vt.SafeNewLine)
}

func prettyPrintDefault_spf13_pflag(flag *spf13_pflag.Flag) {
	typeName, usage := spf13_pflag.UnquoteUsage(flag)
	fmt.Fprintf(os.Stderr, prettyPrintDefault(flag.Name, flag.Shorthand, typeName, usage, flag.DefValue, flag.Value), vt.SafeNewLine)
}

func prettyPrintDefault(name, shorthand, typeName, usage, defValue string, value interface{}) string {
	var s string
	if len(shorthand) > 0 {
		s = "  " + vt.SafeB("-"+shorthand) + ", " + vt.SafeB("--"+name)
	} else {
		s = "  " + vt.SafeB("--"+name)
	}
	if len(typeName) > 0 {
		s += " " + vt.SafeU(typeName)
	}
	s += vt.SafeNewLine + "    \t"
	s += usage
	if !IsZeroValue(defValue) {
		if reflect.TypeOf(value).Kind() == reflect.String {
			s += fmt.Sprintf(" (default %q)", defValue)
		} else {
			s += fmt.Sprintf(" (default %v)", defValue)
		}
	}
	return s
}

// Standard

func (fs *FlagSet) Args() []string {
	if f, ok := fs.impl.(interface{ Args() []string }); ok {
		return f.Args()
	}
	panic(error911.NewNotSupported("Args"))
}

func (fs *FlagSet) Parse(arguments []string) error {
	if f, ok := fs.impl.(interface{ Parse([]string) error }); ok {
		return f.Parse(arguments)
	}
	panic(error911.NewNotSupported("Parse"))
}

func (fs *FlagSet) PrintDefaults() {
	if f, ok := fs.impl.(interface{ PrintDefaults() }); ok {
		f.PrintDefaults()
		return
	}
	panic(error911.NewNotSupported("PrintDefaults"))
}

func (fs *FlagSet) SetOutput(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	fs.output = w
	if f, ok := fs.impl.(interface{ SetOutput(io.Writer) }); ok {
		f.SetOutput(w)
		return
	}
	panic(error911.NewNotSupported("SetOutput"))
}

// Usage

func (fs *FlagSet) Usage(user ...interface{}) error {
	if text := fs.usageHeader; len(text) > 0 {
		if _, err := io.WriteString(fs.output, text); err != nil {
			return err
		}
	} else {
		if len(fs.name) > 0 {
			if _, err := io.WriteString(fs.output, "Usage of "+fs.name+":\n"); err != nil {
				return err
			}
		} else {
			if _, err := io.WriteString(fs.output, "Usage:\n"); err != nil {
				return err
			}
		}
	}
	if fs.usage != nil {
		if f, ok := fs.usage.(func(...interface{}) error); ok {
			if err := f(user...); err != nil {
				return err
			}
		} else if f, ok := fs.usage.(func() error); ok {
			if err := f(); err != nil {
				return err
			}
		} else if f, ok := fs.usage.(func()); ok {
			f()
		}
	} else {
		fs.PrettyPrintDefaults()
	}
	if text := fs.usageFooter; len(text) > 0 {
		if _, err := io.WriteString(fs.output, text); err != nil {
			return err
		}
	}
	return nil
}
func (fs *FlagSet) UsageFallback() {
	fs.Usage()
}
func (fs *FlagSet) UsageFallbackError() error {
	return fs.Usage()
}
func (fs *FlagSet) SetUsage(cb interface{}) {
	if _, ok := cb.(func(...interface{}) error); ok {
		fs.usage = cb
	} else if _, ok := cb.(func() error); ok {
		fs.usage = cb
	} else if _, ok := cb.(func()); ok {
		fs.usage = cb
	} else {
		panic("Argument to SetUsage must be `func(...interface{}) error`, `func() error`, or `func()`")
	}
}
func (fs *FlagSet) SetUsageHeader(text string) {
	fs.usageHeader = text
}
func (fs *FlagSet) SetUsageFooter(text string) {
	fs.usageFooter = text
}

// Feature tests

func (fs *FlagSet) CountWorks() bool {
	if f, ok := fs.impl.(interface{ CountWorks() bool }); ok {
		return f.CountWorks()
	}
	return false
}
func (fs *FlagSet) ShorthandWorks() bool {
	if f, ok := fs.impl.(interface{ PVersionsWork() bool }); ok {
		return f.PVersionsWork()
	}
	return false
}

// Bool

func (fs *FlagSet) Bool(name string, value bool, usage string) *bool {
	if fs.runHooks("Bool", nil, name, "", value, usage) {
		panic(error911.NewCancel("Bool")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Bool(string, bool, string) *bool
	}); ok {
		return f.Bool(name, value, usage)
	}
	panic(error911.NewNotSupported("Bool"))
}
func (fs *FlagSet) BoolP(name string, shorthand string, value bool, usage string) *bool {
	if fs.runHooks("BoolP", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("BoolP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		BoolP(string, string, bool, string) *bool
	}); ok {
		return f.BoolP(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Bool(string, bool, string) *bool
	}); ok {
		return f.Bool(name, value, usage)
	}
	panic(error911.NewNotSupported("BoolP"))
}
func (fs *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	if fs.runHooks("BoolVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		BoolVar(*bool, string, bool, string)
	}); ok {
		f.BoolVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("BoolVar"))
}
func (fs *FlagSet) BoolVarP(p *bool, name string, shorthand string, value bool, usage string) {
	if fs.runHooks("BoolVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		BoolVarP(*bool, string, string, bool, string)
	}); ok {
		f.BoolVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		BoolVar(*bool, string, bool, string)
	}); ok {
		f.BoolVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("BoolVarP"))
}

// Count

func (fs *FlagSet) Count(name string, value int, usage string) *int {
	if fs.runHooks("Count", nil, name, "", value, usage) {
		panic(error911.NewCancel("Count")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Count(string, int, string) *int
	}); ok {
		return f.Count(name, value, usage)
	}
	panic(error911.NewNotSupported("Count"))
}
func (fs *FlagSet) CountP(name string, shorthand string, value int, usage string) *int {
	if fs.runHooks("CountP", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("CountP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		CountP(string, string, int, string) *int
	}); ok {
		return f.CountP(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Count(string, int, string) *int
	}); ok {
		return f.Count(name, value, usage)
	}
	panic(error911.NewNotSupported("CountP"))
}
func (fs *FlagSet) CountVar(p *int, name string, value int, usage string) {
	if fs.runHooks("CountVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		CountVar(*int, string, int, string)
	}); ok {
		f.CountVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("CountVar"))
}
func (fs *FlagSet) CountVarP(p *int, name string, shorthand string, value int, usage string) {
	if fs.runHooks("CountVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		CountVarP(*int, string, string, int, string)
	}); ok {
		f.CountVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		CountVar(*int, string, int, string)
	}); ok {
		f.CountVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("CountVarP"))
}

// Duration

func (fs *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	if fs.runHooks("Duration", nil, name, "", value, usage) {
		panic(error911.NewCancel("Duration")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Duration(string, time.Duration, string) *time.Duration
	}); ok {
		return f.Duration(name, value, usage)
	}
	panic(error911.NewNotSupported("Duration"))
}
func (fs *FlagSet) DurationP(name string, shorthand string, value time.Duration, usage string) *time.Duration {
	if fs.runHooks("DurationP", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("DurationP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		DurationP(string, string, time.Duration, string) *time.Duration
	}); ok {
		return f.DurationP(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Duration(string, time.Duration, string) *time.Duration
	}); ok {
		return f.Duration(name, value, usage)
	}
	panic(error911.NewNotSupported("DurationP"))
}
func (fs *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	if fs.runHooks("DurationVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		DurationVar(*time.Duration, string, time.Duration, string)
	}); ok {
		f.DurationVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("DurationVar"))
}
func (fs *FlagSet) DurationVarP(p *time.Duration, name string, shorthand string, value time.Duration, usage string) {
	if fs.runHooks("DurationVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		DurationVarP(*time.Duration, string, string, time.Duration, string)
	}); ok {
		f.DurationVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		DurationVar(*time.Duration, string, time.Duration, string)
	}); ok {
		f.DurationVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("DurationVarP"))
}

// Float32

func (fs *FlagSet) Float32(name string, value float32, usage string) *float32 {
	if fs.runHooks("Float32", nil, name, "", value, usage) {
		panic(error911.NewCancel("Float32")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Float32(string, float32, string) *float32
	}); ok {
		return f.Float32(name, value, usage)
	}
	panic(error911.NewNotSupported("Float32"))
}
func (fs *FlagSet) Float32P(name string, shorthand string, value float32, usage string) *float32 {
	if fs.runHooks("Float32P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Float32P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Float32P(string, string, float32, string) *float32
	}); ok {
		return f.Float32P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Float32(string, float32, string) *float32
	}); ok {
		return f.Float32(name, value, usage)
	}
	panic(error911.NewNotSupported("Float32P"))
}
func (fs *FlagSet) Float32Var(p *float32, name string, value float32, usage string) {
	if fs.runHooks("Float32Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Float32Var(*float32, string, float32, string)
	}); ok {
		f.Float32Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Float32Var"))
}
func (fs *FlagSet) Float32VarP(p *float32, name string, shorthand string, value float32, usage string) {
	if fs.runHooks("Float32VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Float32VarP(*float32, string, string, float32, string)
	}); ok {
		f.Float32VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Float32Var(*float32, string, float32, string)
	}); ok {
		f.Float32Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Float32VarP"))
}

// Float64

func (fs *FlagSet) Float64(name string, value float64, usage string) *float64 {
	if fs.runHooks("Float64", nil, name, "", value, usage) {
		panic(error911.NewCancel("Float64")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Float64(string, float64, string) *float64
	}); ok {
		return f.Float64(name, value, usage)
	}
	panic(error911.NewNotSupported("Float64"))
}
func (fs *FlagSet) Float64P(name string, shorthand string, value float64, usage string) *float64 {
	if fs.runHooks("Float64P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Float64P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Float64P(string, string, float64, string) *float64
	}); ok {
		return f.Float64P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Float64(string, float64, string) *float64
	}); ok {
		return f.Float64(name, value, usage)
	}
	panic(error911.NewNotSupported("Float64P"))
}
func (fs *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
	if fs.runHooks("Float64Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Float64Var(*float64, string, float64, string)
	}); ok {
		f.Float64Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Float64Var"))
}
func (fs *FlagSet) Float64VarP(p *float64, name string, shorthand string, value float64, usage string) {
	if fs.runHooks("Float64VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Float64VarP(*float64, string, string, float64, string)
	}); ok {
		f.Float64VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Float64Var(*float64, string, float64, string)
	}); ok {
		f.Float64Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Float64VarP"))
}

// Int

func (fs *FlagSet) Int(name string, value int, usage string) *int {
	if fs.runHooks("Int", nil, name, "", value, usage) {
		panic(error911.NewCancel("Int")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Int(string, int, string) *int
	}); ok {
		return f.Int(name, value, usage)
	}
	panic(error911.NewNotSupported("Int"))
}
func (fs *FlagSet) IntP(name string, shorthand string, value int, usage string) *int {
	if fs.runHooks("IntP", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("IntP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		IntP(string, string, int, string) *int
	}); ok {
		return f.IntP(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Int(string, int, string) *int
	}); ok {
		return f.Int(name, value, usage)
	}
	panic(error911.NewNotSupported("IntP"))
}
func (fs *FlagSet) IntVar(p *int, name string, value int, usage string) {
	if fs.runHooks("IntVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		IntVar(*int, string, int, string)
	}); ok {
		f.IntVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("IntVar"))
}
func (fs *FlagSet) IntVarP(p *int, name string, shorthand string, value int, usage string) {
	if fs.runHooks("IntVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		IntVarP(*int, string, string, int, string)
	}); ok {
		f.IntVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		IntVar(*int, string, int, string)
	}); ok {
		f.IntVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("IntVarP"))
}

// Int8

func (fs *FlagSet) Int8(name string, value int8, usage string) *int8 {
	if fs.runHooks("Int8", nil, name, "", value, usage) {
		panic(error911.NewCancel("Int8")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Int8(string, int8, string) *int8
	}); ok {
		return f.Int8(name, value, usage)
	}
	panic(error911.NewNotSupported("Int8"))
}
func (fs *FlagSet) Int8P(name string, shorthand string, value int8, usage string) *int8 {
	if fs.runHooks("Int8P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Int8P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Int8P(string, string, int8, string) *int8
	}); ok {
		return f.Int8P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Int8(string, int8, string) *int8
	}); ok {
		return f.Int8(name, value, usage)
	}
	panic(error911.NewNotSupported("Int8P"))
}
func (fs *FlagSet) Int8Var(p *int8, name string, value int8, usage string) {
	if fs.runHooks("Int8Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Int8Var(*int8, string, int8, string)
	}); ok {
		f.Int8Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Int8Var"))
}
func (fs *FlagSet) Int8VarP(p *int8, name string, shorthand string, value int8, usage string) {
	if fs.runHooks("Int8VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Int8VarP(*int8, string, string, int8, string)
	}); ok {
		f.Int8VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Int8Var(*int8, string, int8, string)
	}); ok {
		f.Int8Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Int8VarP"))
}

// Int16

func (fs *FlagSet) Int16(name string, value int16, usage string) *int16 {
	if fs.runHooks("Int16", nil, name, "", value, usage) {
		panic(error911.NewCancel("Int16")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Int16(string, int16, string) *int16
	}); ok {
		return f.Int16(name, value, usage)
	}
	panic(error911.NewNotSupported("Int16"))
}
func (fs *FlagSet) Int16P(name string, shorthand string, value int16, usage string) *int16 {
	if fs.runHooks("Int16P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Int16P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Int16P(string, string, int16, string) *int16
	}); ok {
		return f.Int16P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Int16(string, int16, string) *int16
	}); ok {
		return f.Int16(name, value, usage)
	}
	panic(error911.NewNotSupported("Int16P"))
}
func (fs *FlagSet) Int16Var(p *int16, name string, value int16, usage string) {
	if fs.runHooks("Int16Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Int16Var(*int16, string, int16, string)
	}); ok {
		f.Int16Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Int16Var"))
}
func (fs *FlagSet) Int16VarP(p *int16, name string, shorthand string, value int16, usage string) {
	if fs.runHooks("Int16VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Int16VarP(*int16, string, string, int16, string)
	}); ok {
		f.Int16VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Int16Var(*int16, string, int16, string)
	}); ok {
		f.Int16Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Int16VarP"))
}

// Int32

func (fs *FlagSet) Int32(name string, value int32, usage string) *int32 {
	if fs.runHooks("Int32", nil, name, "", value, usage) {
		panic(error911.NewCancel("Int32")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Int32(string, int32, string) *int32
	}); ok {
		return f.Int32(name, value, usage)
	}
	panic(error911.NewNotSupported("Int32"))
}
func (fs *FlagSet) Int32P(name string, shorthand string, value int32, usage string) *int32 {
	if fs.runHooks("Int32P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Int32P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Int32P(string, string, int32, string) *int32
	}); ok {
		return f.Int32P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Int32(string, int32, string) *int32
	}); ok {
		return f.Int32(name, value, usage)
	}
	panic(error911.NewNotSupported("Int32P"))
}
func (fs *FlagSet) Int32Var(p *int32, name string, value int32, usage string) {
	if fs.runHooks("Int32Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Int32Var(*int32, string, int32, string)
	}); ok {
		f.Int32Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Int32Var"))
}
func (fs *FlagSet) Int32VarP(p *int32, name string, shorthand string, value int32, usage string) {
	if fs.runHooks("Int32VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Int32VarP(*int32, string, string, int32, string)
	}); ok {
		f.Int32VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Int32Var(*int32, string, int32, string)
	}); ok {
		f.Int32Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Int32VarP"))
}

// Int64

func (fs *FlagSet) Int64(name string, value int64, usage string) *int64 {
	if fs.runHooks("Int64", nil, name, "", value, usage) {
		panic(error911.NewCancel("Int64")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Int64(string, int64, string) *int64
	}); ok {
		return f.Int64(name, value, usage)
	}
	panic(error911.NewNotSupported("Int64"))
}
func (fs *FlagSet) Int64P(name string, shorthand string, value int64, usage string) *int64 {
	if fs.runHooks("Int64P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Int64P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Int64P(string, string, int64, string) *int64
	}); ok {
		return f.Int64P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Int64(string, int64, string) *int64
	}); ok {
		return f.Int64(name, value, usage)
	}
	panic(error911.NewNotSupported("Int64P"))
}
func (fs *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	if fs.runHooks("Int64Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Int64Var(*int64, string, int64, string)
	}); ok {
		f.Int64Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Int64Var"))
}
func (fs *FlagSet) Int64VarP(p *int64, name string, shorthand string, value int64, usage string) {
	if fs.runHooks("Int64VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Int64VarP(*int64, string, string, int64, string)
	}); ok {
		f.Int64VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Int64Var(*int64, string, int64, string)
	}); ok {
		f.Int64Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Int64VarP"))
}

// IP

func (fs *FlagSet) IP(name string, value net.IP, usage string) *net.IP {
	if fs.runHooks("IP", nil, name, "", value, usage) {
		panic(error911.NewCancel("IP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		IP(string, net.IP, string) *net.IP
	}); ok {
		return f.IP(name, value, usage)
	}
	panic(error911.NewNotSupported("IP"))
}
func (fs *FlagSet) IPP(name string, shorthand string, value net.IP, usage string) *net.IP {
	if fs.runHooks("IPP", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("IPP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		IPP(string, string, net.IP, string) *net.IP
	}); ok {
		return f.IPP(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		IP(string, net.IP, string) *net.IP
	}); ok {
		return f.IP(name, value, usage)
	}
	panic(error911.NewNotSupported("IPP"))
}
func (fs *FlagSet) IPVar(p *net.IP, name string, value net.IP, usage string) {
	if fs.runHooks("IPVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		IPVar(*net.IP, string, net.IP, string)
	}); ok {
		f.IPVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("IPVar"))
}
func (fs *FlagSet) IPVarP(p *net.IP, name string, shorthand string, value net.IP, usage string) {
	if fs.runHooks("IPVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		IPVarP(*net.IP, string, string, net.IP, string)
	}); ok {
		f.IPVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		IPVar(*net.IP, string, net.IP, string)
	}); ok {
		f.IPVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("IPVarP"))
}

// IPMask

func (fs *FlagSet) IPMask(name string, value net.IPMask, usage string) *net.IPMask {
	if fs.runHooks("IPMask", nil, name, "", value, usage) {
		panic(error911.NewCancel("IPMask")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		IPMask(string, net.IPMask, string) *net.IPMask
	}); ok {
		return f.IPMask(name, value, usage)
	}
	panic(error911.NewNotSupported("IPMask"))
}
func (fs *FlagSet) IPMaskP(name string, shorthand string, value net.IPMask, usage string) *net.IPMask {
	if fs.runHooks("IPMaskP", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("IPMaskP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		IPMaskP(string, string, net.IPMask, string) *net.IPMask
	}); ok {
		return f.IPMaskP(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		IPMask(string, net.IPMask, string) *net.IPMask
	}); ok {
		return f.IPMask(name, value, usage)
	}
	panic(error911.NewNotSupported("IPMaskP"))
}
func (fs *FlagSet) IPMaskVar(p *net.IPMask, name string, value net.IPMask, usage string) {
	if fs.runHooks("IPMaskVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		IPMaskVar(*net.IPMask, string, net.IPMask, string)
	}); ok {
		f.IPMaskVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("IPMaskVar"))
}
func (fs *FlagSet) IPMaskVarP(p *net.IPMask, name string, shorthand string, value net.IPMask, usage string) {
	if fs.runHooks("IPMaskVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		IPMaskVarP(*net.IPMask, string, string, net.IPMask, string)
	}); ok {
		f.IPMaskVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		IPMaskVar(*net.IPMask, string, net.IPMask, string)
	}); ok {
		f.IPMaskVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("IPMaskVarP"))
}

// String

func (fs *FlagSet) String(name string, value string, usage string) *string {
	if fs.runHooks("String", nil, name, "", value, usage) {
		panic(error911.NewCancel("String")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		String(string, string, string) *string
	}); ok {
		return f.String(name, value, usage)
	}
	panic(error911.NewNotSupported("String"))
}
func (fs *FlagSet) StringP(name string, shorthand string, value string, usage string) *string {
	if fs.runHooks("StringP", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("StringP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		StringP(string, string, string, string) *string
	}); ok {
		return f.StringP(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		String(string, string, string) *string
	}); ok {
		return f.String(name, value, usage)
	}
	panic(error911.NewNotSupported("StringP"))
}
func (fs *FlagSet) StringVar(p *string, name string, value string, usage string) {
	if fs.runHooks("StringVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		StringVar(*string, string, string, string)
	}); ok {
		f.StringVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("StringVar"))
}
func (fs *FlagSet) StringVarP(p *string, name string, shorthand string, value string, usage string) {
	if fs.runHooks("StringVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		StringVarP(*string, string, string, string, string)
	}); ok {
		f.StringVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		StringVar(*string, string, string, string)
	}); ok {
		f.StringVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("StringVarP"))
}

// Time

func (fs *FlagSet) Time(name string, value time.Time, usage string) *time.Time {
	if fs.runHooks("Time", nil, name, "", value, usage) {
		panic(error911.NewCancel("Time")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Time(string, time.Time, string) *time.Time
	}); ok {
		return f.Time(name, value, usage)
	}
	panic(error911.NewNotSupported("Time"))
}
func (fs *FlagSet) TimeP(name string, shorthand string, value time.Time, usage string) *time.Time {
	if fs.runHooks("TimeP", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("TimeP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		TimeP(string, string, time.Time, string) *time.Time
	}); ok {
		return f.TimeP(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Time(string, time.Time, string) *time.Time
	}); ok {
		return f.Time(name, value, usage)
	}
	panic(error911.NewNotSupported("TimeP"))
}
func (fs *FlagSet) TimeVar(p *time.Time, name string, value time.Time, usage string) {
	if fs.runHooks("TimeVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		TimeVar(*time.Time, string, time.Time, string)
	}); ok {
		f.TimeVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("TimeVar"))
}
func (fs *FlagSet) TimeVarP(p *time.Time, name string, shorthand string, value time.Time, usage string) {
	if fs.runHooks("TimeVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		TimeVarP(*time.Time, string, string, time.Time, string)
	}); ok {
		f.TimeVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		TimeVar(*time.Time, string, time.Time, string)
	}); ok {
		f.TimeVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("TimeVarP"))
}

// Uint

func (fs *FlagSet) Uint(name string, value uint, usage string) *uint {
	if fs.runHooks("Uint", nil, name, "", value, usage) {
		panic(error911.NewCancel("Uint")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Uint(string, uint, string) *uint
	}); ok {
		return f.Uint(name, value, usage)
	}
	panic(error911.NewNotSupported("Uint"))
}
func (fs *FlagSet) UintP(name string, shorthand string, value uint, usage string) *uint {
	if fs.runHooks("UintP", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("UintP")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		UintP(string, string, uint, string) *uint
	}); ok {
		return f.UintP(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Uint(string, uint, string) *uint
	}); ok {
		return f.Uint(name, value, usage)
	}
	panic(error911.NewNotSupported("UintP"))
}
func (fs *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
	if fs.runHooks("UintVar", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		UintVar(*uint, string, uint, string)
	}); ok {
		f.UintVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("UintVar"))
}
func (fs *FlagSet) UintVarP(p *uint, name string, shorthand string, value uint, usage string) {
	if fs.runHooks("UintVarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		UintVarP(*uint, string, string, uint, string)
	}); ok {
		f.UintVarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		UintVar(*uint, string, uint, string)
	}); ok {
		f.UintVar(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("UintVarP"))
}

// Uint8

func (fs *FlagSet) Uint8(name string, value uint8, usage string) *uint8 {
	if fs.runHooks("Uint8", nil, name, "", value, usage) {
		panic(error911.NewCancel("Uint8")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Uint8(string, uint8, string) *uint8
	}); ok {
		return f.Uint8(name, value, usage)
	}
	panic(error911.NewNotSupported("Uint8"))
}
func (fs *FlagSet) Uint8P(name string, shorthand string, value uint8, usage string) *uint8 {
	if fs.runHooks("Uint8P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Uint8P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Uint8P(string, string, uint8, string) *uint8
	}); ok {
		return f.Uint8P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Uint8(string, uint8, string) *uint8
	}); ok {
		return f.Uint8(name, value, usage)
	}
	panic(error911.NewNotSupported("Uint8P"))
}
func (fs *FlagSet) Uint8Var(p *uint8, name string, value uint8, usage string) {
	if fs.runHooks("Uint8Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Uint8Var(*uint8, string, uint8, string)
	}); ok {
		f.Uint8Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Uint8Var"))
}
func (fs *FlagSet) Uint8VarP(p *uint8, name string, shorthand string, value uint8, usage string) {
	if fs.runHooks("Uint8VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Uint8VarP(*uint8, string, string, uint8, string)
	}); ok {
		f.Uint8VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Uint8Var(*uint8, string, uint8, string)
	}); ok {
		f.Uint8Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Uint8VarP"))
}

// Uint16

func (fs *FlagSet) Uint16(name string, value uint16, usage string) *uint16 {
	if fs.runHooks("Uint16", nil, name, "", value, usage) {
		panic(error911.NewCancel("Uint16")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Uint16(string, uint16, string) *uint16
	}); ok {
		return f.Uint16(name, value, usage)
	}
	panic(error911.NewNotSupported("Uint16"))
}
func (fs *FlagSet) Uint16P(name string, shorthand string, value uint16, usage string) *uint16 {
	if fs.runHooks("Uint16P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Uint16P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Uint16P(string, string, uint16, string) *uint16
	}); ok {
		return f.Uint16P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Uint16(string, uint16, string) *uint16
	}); ok {
		return f.Uint16(name, value, usage)
	}
	panic(error911.NewNotSupported("Uint16P"))
}
func (fs *FlagSet) Uint16Var(p *uint16, name string, value uint16, usage string) {
	if fs.runHooks("Uint16Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Uint16Var(*uint16, string, uint16, string)
	}); ok {
		f.Uint16Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Uint16Var"))
}
func (fs *FlagSet) Uint16VarP(p *uint16, name string, shorthand string, value uint16, usage string) {
	if fs.runHooks("Uint16VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Uint16VarP(*uint16, string, string, uint16, string)
	}); ok {
		f.Uint16VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Uint16Var(*uint16, string, uint16, string)
	}); ok {
		f.Uint16Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Uint16VarP"))
}

// Uint32

func (fs *FlagSet) Uint32(name string, value uint32, usage string) *uint32 {
	if fs.runHooks("Uint32", nil, name, "", value, usage) {
		panic(error911.NewCancel("Uint32")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Uint32(string, uint32, string) *uint32
	}); ok {
		return f.Uint32(name, value, usage)
	}
	panic(error911.NewNotSupported("Uint32"))
}
func (fs *FlagSet) Uint32P(name string, shorthand string, value uint32, usage string) *uint32 {
	if fs.runHooks("Uint32P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Uint32P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Uint32P(string, string, uint32, string) *uint32
	}); ok {
		return f.Uint32P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Uint32(string, uint32, string) *uint32
	}); ok {
		return f.Uint32(name, value, usage)
	}
	panic(error911.NewNotSupported("Uint32P"))
}
func (fs *FlagSet) Uint32Var(p *uint32, name string, value uint32, usage string) {
	if fs.runHooks("Uint32Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Uint32Var(*uint32, string, uint32, string)
	}); ok {
		f.Uint32Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Uint32Var"))
}
func (fs *FlagSet) Uint32VarP(p *uint32, name string, shorthand string, value uint32, usage string) {
	if fs.runHooks("Uint32VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Uint32VarP(*uint32, string, string, uint32, string)
	}); ok {
		f.Uint32VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Uint32Var(*uint32, string, uint32, string)
	}); ok {
		f.Uint32Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Uint32VarP"))
}

// Uint64

func (fs *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	if fs.runHooks("Uint64", nil, name, "", value, usage) {
		panic(error911.NewCancel("Uint64")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Uint64(string, uint64, string) *uint64
	}); ok {
		return f.Uint64(name, value, usage)
	}
	panic(error911.NewNotSupported("Uint64"))
}
func (fs *FlagSet) Uint64P(name string, shorthand string, value uint64, usage string) *uint64 {
	if fs.runHooks("Uint64P", nil, name, shorthand, value, usage) {
		panic(error911.NewCancel("Uint64P")) // This operation cannot be safely canceled
	}
	if f, ok := fs.impl.(interface {
		Uint64P(string, string, uint64, string) *uint64
	}); ok {
		return f.Uint64P(name, shorthand, value, usage)
	} else if f, ok := fs.impl.(interface {
		Uint64(string, uint64, string) *uint64
	}); ok {
		return f.Uint64(name, value, usage)
	}
	panic(error911.NewNotSupported("Uint64P"))
}
func (fs *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
	if fs.runHooks("Uint64Var", p, name, "", value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Uint64Var(*uint64, string, uint64, string)
	}); ok {
		f.Uint64Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Uint64Var"))
}
func (fs *FlagSet) Uint64VarP(p *uint64, name string, shorthand string, value uint64, usage string) {
	if fs.runHooks("Uint64VarP", p, name, shorthand, value, usage) {
		return
	}
	if f, ok := fs.impl.(interface {
		Uint64VarP(*uint64, string, string, uint64, string)
	}); ok {
		f.Uint64VarP(p, name, shorthand, value, usage)
		return
	} else if f, ok := fs.impl.(interface {
		Uint64Var(*uint64, string, uint64, string)
	}); ok {
		f.Uint64Var(p, name, value, usage)
		return
	}
	panic(error911.NewNotSupported("Uint64VarP"))
}

// Internal

func (fs *FlagSet) runHooks(fn string, p interface{}, name string, shorthand string, value interface{}, usage string) (cancel bool) {
	fs.hooksMutex.RLock()
	defer fs.hooksMutex.RUnlock()
	for _, hook := range fs.hooks {
		if hook.Hook != nil {
			if err := hook.Hook(fs, fn, p, name, shorthand, value, usage, hook.User); err != nil {
				if _, ok := err.(error911.Cancel); ok {
					cancel = true
				} else {
					panic(err)
				}
			}
		}
	}
	return
}
