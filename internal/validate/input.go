package validate

import (
	"fmt"
	"strings"
)

var validFeatures = []string{"stress", "ramp", "spike", "pulse", "quick", "explain"}
var validModes    = []string{"constant", "ramp", "spike", "pulse"}
var validOutputs  = []string{"simple", "detailed", "json"}

func Feature(f string) error {
	f = strings.ToLower(strings.TrimSpace(f))
	for _, v := range validFeatures {
		if f == v {
			return nil
		}
	}
	return fmt.Errorf("unknown feature %q — choose from: %s", f, strings.Join(validFeatures, ", "))
}

func Mode(m string) error {
	if m == "" {
		return nil // optional, no validation needed
	}
	m = strings.ToLower(strings.TrimSpace(m))
	for _, v := range validModes {
		if m == v {
			return nil
		}
	}
	return fmt.Errorf("unknown mode %q — choose from: %s", m, strings.Join(validModes, ", "))
}

func Output(o string) error {
	o = strings.ToLower(strings.TrimSpace(o))
	for _, v := range validOutputs {
		if o == v {
			return nil
		}
	}
	return fmt.Errorf("unknown output format %q — choose from: %s", o, strings.Join(validOutputs, ", "))
}

func Rate(r int) error {
	if r <= 0 {
		return fmt.Errorf("--rate must be greater than 0 (got %d)", r)
	}
	return nil
}

func Threads(t int) error {
	if t <= 0 {
		return fmt.Errorf("--threads must be greater than 0 (got %d)", t)
	}
	return nil
}

func Duration(d int) error {
	if d <= 0 {
		return fmt.Errorf("--duration must be greater than 0 (got %d)", d)
	}
	return nil
}
