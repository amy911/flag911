package flagger

import flag "github.com/spf13/pflag"

type SPF13FlagSet struct {
	Common
	NoCount
	NoCountP
	NoInt16
	NoInt16P
	*flag.FlagSet
}

// Usage

func (fs SPF13FlagSet) Usage() interface{} {
	flag.Usage()
	return nil
}

func (fs SPF13FlagSet) SetUsage(cb func(interface{}) interface{}, user interface{}, flbk func()) {
	flag.Usage = flbk
}

func (fs SPF13FlagSet) PVersionsWork() bool {
	return true
}
