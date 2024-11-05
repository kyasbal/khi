package parameters

import (
	"flag"
	"os"
)

var Help = &HelpParameters{}

type HelpParameters struct {
	// Help
	// If this parameter is set, KHI exits with printing the usage.
	Help *bool
}

func (h *HelpParameters) PostProcess() error {
	if *h.Help {
		flag.PrintDefaults()
		os.Exit(0)
		return nil
	}
	return nil
}

func (h *HelpParameters) Prepare() error {
	h.Help = flag.Bool("help", false, "If this flag is set, KHI exits with printing the usage.")
	return nil
}

var _ ParameterStore = (*HelpParameters)(nil)
