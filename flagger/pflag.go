package flagger

import flag "github.com/ogier/pflag"

type FlagSet struct {
	NoCount
	NoCountP
	flag.FlagSet
}

// Usage

func (fs *FlagSet) Usage() interface{} {
	flag.Usage()
	return nil
}

func (fs *FlagSet) SetUsage(cb func(interface{}) interface{}, user interface{}, flbk func()) {
	flag.Usage = flbk
}
