package flagger

import (
	"flag"
	"github.com/ogier/pflag"
	spf13 "github.com/spf13/pflag"
)

type FlagSet struct {
	Flagger
}

func New(args ...interface{}) *FlagSet {
	return new(FlagSet).Init(args...)
}

func (fs *FlagSet) Init(args ...interface{}) *FlagSet {
	switch len(flagger) {
	case 0:
		fs.Flagger = PFlagSet{FlagSet: pflag.CommandLine}
	case 1:
		switch args[0].(type) {
		case Flagger:
			fs.Flagger = args[0].(Flagger)
		case *flag.FlagSet:
			fs.Flagger = FlagSet{FlagSet: args[0].(*flag.FlagSet)}
		case *pflag.FlagSet:
			fs.Flagger = PFlagSet{FlagSet: args[0].(*pflag.FlagSet)}
		case *spf13.FlagSet:
			fs.Flagger = SPF13FlagSet{FlagSet: args[0].(*spf13.FlagSet)}
		default:
			panic("Unsupported initialization argument")
		}
	default:
		panic("Provide up to one Flagger instance")
	}
	return fs
}
