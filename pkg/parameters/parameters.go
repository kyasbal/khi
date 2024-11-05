package parameters

import (
	"errors"
	"flag"
)

// ParameterStore parses subset of parameters given to KHI program.
type ParameterStore interface {
	// Prepare initialize settings for reading parameters with `flag` package.
	Prepare() error

	// PostProcess override parsed parameters depending on the other parsed parameters.
	PostProcess() error
}

// Parse initializes the given parameter stores.
func Parse(stores ...ParameterStore) error {
	if flag.Parsed() {
		return errors.New("parameter flags are already parsed")
	}
	for _, store := range stores {
		err := store.Prepare()
		if err != nil {
			return err
		}
	}
	flag.Parse()
	for _, store := range stores {
		err := store.PostProcess()
		if err != nil {
			return err
		}
	}
	return nil
}
