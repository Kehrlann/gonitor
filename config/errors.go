package config

import "fmt"

// NoDefaultConfigError represents an error thrown when the user hasn't specified a config, and that config wasn't found
type NoDefaultConfigError struct {
	HelpMessage string
}

// Error is the error message for NoDefautConfigError
func (n NoDefaultConfigError) Error() string {
	return fmt.Sprintf("No default config found at %v", DEFAULT_CONFIG_PATH)
}

func NewDefaultConfigError() *NoDefaultConfigError {
	return &NoDefaultConfigError{HelpMessage: fmt.Sprintf(`
It seems you didn't specify a config file. Gonitor attempted to load a default config file, located at  :
	%v

No such file was found. Please create one, or specify an existing config with the -config flag. If you wish to create
one, you can find a start template by running gonitor -example`, DEFAULT_CONFIG_PATH)}
}
