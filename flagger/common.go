package flagger

type Common struct {
}

// Usage

func (fs Common) SetUsageFooter(text string) {
	fs.SetUsage(nil, nil, UsageWithHeaderAndFooter)
	UsageFooter = text
}

func (fs Common) SetUsageHeader(text string) {
	fs.SetUsage(nil, nil, UsageWithHeaderAndFooter)
	UsageHeader = text
}
