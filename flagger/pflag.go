package flagger

import flag "github.com/ogier/pflag"

type PFlagSet struct {
	Common
	NoCount
	NoCountP
	NoInt16
	NoInt16P
	*flag.FlagSet
}

// Usage

func (fs PFlagSet) Usage() interface{} {
	flag.Usage()
	return nil
}

func (fs PFlagSet) SetUsage(cb func(interface{}) interface{}, user interface{}, flbk func()) {
	flag.Usage = flbk
}

func (fs PFlagSet) PVersionsWork() bool {
	return true
}
