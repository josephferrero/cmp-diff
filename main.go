package main

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

type Driver struct {
	Name string
	Age  int8
	Vehicle
}

type Vehicle struct {
	Make  string
	Model string
	Color string
	Year  string
	LicensePlate
}

type LicensePlate struct {
	Number   string
	ExpMonth string
}

var driver1 = Driver{
	Name: "John Smith",
	Age:  29,
	Vehicle: Vehicle{
		Make:  "Ford",
		Model: "Ranger",
		Color: "White",
		Year:  "2006",
		LicensePlate: LicensePlate{
			Number:   "1AAA001",
			ExpMonth: "March",
		},
	},
}

var driver2 = Driver{
	Name: "Jane Smith",
	Age:  29,
	Vehicle: Vehicle{
		Make:  "Toyta",
		Model: "Sienna",
		Color: "Grey",
		Year:  "2011",
		LicensePlate: LicensePlate{
			Number:   "1AAA000",
			ExpMonth: "March",
		},
	},
}

func main() {
	r := NewDiffReporter(NewRedactedFields(LicensePlate{}, "Number"), NewRedactedFields(Driver{}, "Name"))

	cmp.Equal(driver1, driver2, cmp.Reporter(&r))
	fmt.Print(r.String())
}
