package region

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleGetLvls() {
	levels := "7 - 9"
	fromLvl, toLvl, err := GetLvls(levels)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("from level: %d\nto level: %d\n", fromLvl, toLvl)

	// Output:
	// from level: 7
	// to level: 9
}

func ExampleGetTiles() {
	fromLvl, toLvl, accuracy := 10, 10, -1
	fld := float64(5)
	coords := &Coordinates{"maxLat": 60.120754, "maxLon": 30.699476, "minLat": 59.805323, "minLon": 30.024240}

	for lvl := fromLvl; lvl <= toLvl; lvl++ {
		tiles := GetTiles(lvl, accuracy, fld, coords)
		fmt.Printf("For level %d\n", tiles.Level)
		fmt.Printf("y range = %d - %d, x range = %d - %d\n", tiles.Range.MinLat, tiles.Range.MaxLat, tiles.Range.MinLon, tiles.Range.MaxLon)
		fmt.Printf("Min coordinates: %v, %v\n", tiles.Coordinates["minLat"], tiles.Coordinates["minLon"])
		fmt.Printf("Max coordinates: %v, %v\n", tiles.Coordinates["maxLat"], tiles.Coordinates["maxLon"])
	}

	// Output:
	// For level 10
	// y range = 4261 - 4270, x range = 5974 - 5993
	// Min coordinates: 59.80078125, 30.0234375
	// Max coordinates: 60.15234375, 30.7265625
}

func TestAbsCoordinates(t *testing.T) {
	tests := []struct {
		name       string
		curLvl     int
		accuracy   int
		fld        float64
		origCoords Coordinates
		want       Coordinates
	}{
		{
			name:       "correct curLvl=10 accuracy=-1 fld=5",
			curLvl:     10,
			accuracy:   -1,
			fld:        5,
			origCoords: Coordinates{"maxLat": 60.120754, "maxLon": 30.699476, "minLat": 59.805323, "minLon": 30.024240},
			want:       Coordinates{"maxLat": 60.15234375, "maxLon": 30.7265625, "minLat": 59.80078125, "minLon": 30.0234375},
		},
		{
			name:       "correct curLvl=10 accuracy=-1 fld=2",
			curLvl:     10,
			accuracy:   -1,
			fld:        2,
			origCoords: Coordinates{"maxLat": 60.120754, "maxLon": 30.699476, "minLat": 59.805323, "minLon": 30.024240},
			want:       Coordinates{"maxLat": 60.205078125, "maxLon": 30.76171875, "minLat": 59.765625, "minLon": 29.970703125},
		},
		{
			name:       "correct curLvl=10 accuracy=6 fld=5",
			curLvl:     10,
			accuracy:   6,
			fld:        5,
			origCoords: Coordinates{"maxLat": 60.120754, "maxLon": 30.699476, "minLat": 59.805323, "minLon": 30.024240},
			want:       Coordinates{"maxLat": 60.1875, "maxLon": 30.9375, "minLat": 59.625, "minLon": 29.8125},
		},
		{
			name:       "correct curLvl=10 accuracy=6 fld=2",
			curLvl:     10,
			accuracy:   6,
			fld:        2,
			origCoords: Coordinates{"maxLat": 60.120754, "maxLon": 30.699476, "minLat": 59.805323, "minLon": 30.024240},
			want:       Coordinates{"maxLat": 60.46875, "maxLon": 30.9375, "minLat": 59.0625, "minLon": 29.53125},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			absCoords := AbsCoordinates(tc.curLvl, tc.accuracy, tc.fld, tc.origCoords)
			assert.Equal(t, tc.want, absCoords)
		})
	}

}

func TestGetLvls(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		fromLvl int
		toLvl   int
		err     error
	}{
		{name: "case 7-9", input: "7-9", fromLvl: 7, toLvl: 9, err: nil},
		{name: "case 0-11", input: "0-11", fromLvl: 0, toLvl: 11, err: nil},
		{name: "case 0 - 11", input: "0-11", fromLvl: 0, toLvl: 11, err: nil},
		{name: "err with levels parsing", input: "712", fromLvl: 0, toLvl: 0, err: ErrLvlParse},
		{name: "err with fromLvl", input: "Z-9", fromLvl: 0, toLvl: 0, err: strconv.ErrSyntax},
		{name: "err with toLvl", input: "7-Z", fromLvl: 0, toLvl: 0, err: strconv.ErrSyntax},
		{name: "err with range", input: "7-20", fromLvl: 0, toLvl: 0, err: ErrLvlRange},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fromLvl, toLvl, err := GetLvls(tc.input)
			if fromLvl != tc.fromLvl || toLvl != tc.toLvl || !errors.Is(err, tc.err) {
				t.Errorf("%s: expected: %d %d, %v; got: %v %v, %v",
					tc.name, tc.fromLvl, tc.toLvl, tc.err, fromLvl, toLvl, err)
			}
		})
	}
}

func TestGetCoordinates(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		coords *Coordinates
		err    error
	}{
		{
			name:   "correct",
			input:  []string{"59.805323, ", "30.024240", "60.120754, ", "30.699476"},
			coords: &Coordinates{"minLat": 59.805323, "minLon": 30.024240, "maxLat": 60.120754, "maxLon": 30.699476},
			err:    nil,
		},
		{
			name:   "error",
			input:  []string{"59.805323, ", "30.0Z4240", "60.120754, ", "30.699476"},
			coords: nil,
			err:    strconv.ErrSyntax,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			coords, err := GetCoordinates(tc.input)
			if errors.Is(err, tc.err) {
				assert.Equal(t, tc.coords, coords, "coordinates should be equal")
			} else {
				t.Errorf("%s: expected: %#v, %v; got: %v, %v", tc.name, tc.coords, tc.err, coords, err)
			}
		})
	}
}

func TestRemoveCharacters(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		remove string
		want   string
	}{
		{name: "letter", input: "test string", remove: "t", want: "es sring"},
		{name: "letters", input: "test string", remove: "ts", want: "e ring"},
		{name: "whitespace", input: "test string", remove: " ", want: "teststring"},
		{name: "hyphen", input: "test-string", remove: "-", want: "teststring"},
		{name: "punctuation mark", input: "test, string", remove: ",", want: "test string"},
		{name: "punctuation marks", input: "test, string.", remove: ",.", want: "test string"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := RemoveCharacters(tc.input, tc.remove)
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("%s: expected: %s; got: %v", tc.name, tc.want, got)
			}
		})
	}
}
