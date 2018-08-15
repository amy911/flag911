package flagger

import flag "github.com/ogier/pflag"

type PFlagSet struct {
	NoCount
	NoCountP
	*flag.FlagSet
}

// Usage

func (fs *PFlagSet) Usage() interface{} {
	flag.Usage()
	return nil
}

func (fs *PFlagSet) SetUsage(cb func(interface{}) interface{}, user interface{}, flbk func()) {
	flag.Usage = flbk
}
