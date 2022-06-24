package nexrad

import "fmt"

type RangeCheck struct {
	Name  string
	Value float64
	Min   float64
	Max   float64
}

func validateRanges(rangeChecks []*RangeCheck) error {
	for _, check := range rangeChecks {
		if check.Value < check.Min || check.Value > check.Max {
			return fmt.Errorf("Invalid %s | Expected %v to %v | Got %v", check.Name, check.Min, check.Max, check.Value)
		}
	}
	return nil
}
