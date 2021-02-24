package region

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleGetTiles() {
	fromLvl, toLvl, accuracy := 10, 10, -1
	fld := float64(5)
	coords := &Coordinates{MaxLat: 60.120754, MaxLon: 30.699476, MinLat: 59.805323, MinLon: 30.024240}

	for lvl := fromLvl; lvl <= toLvl; lvl++ {
		tiles := GetTiles(lvl, accuracy, fld, coords)
		fmt.Printf("For level %d\n", tiles.Level)
		fmt.Printf("y range = %d - %d, x range = %d - %d\n", tiles.Range.MinY, tiles.Range.MaxY, tiles.Range.MinX, tiles.Range.MaxX)
		fmt.Printf("Min coordinates: %v, %v\n", tiles.Coordinates.MinLat, tiles.Coordinates.MinLon)
		fmt.Printf("Max coordinates: %v, %v\n", tiles.Coordinates.MaxLat, tiles.Coordinates.MaxLon)
	}

	// Output:
	// For level 10
	// y range = 4261 - 4270, x range = 5974 - 5993
	// Min coordinates: 59.80078125, 30.0234375
	// Max coordinates: 60.15234375, 30.7265625
}

func TestCornerCoordinates(t *testing.T) {
	tests := []struct {
		name       string
		curLvl     int
		accuracy   int
		fld        float64
		origCoords *Coordinates
		want       *Coordinates
	}{
		{
			name:       "correct curLvl=10 accuracy=-1 fld=5",
			curLvl:     10,
			accuracy:   -1,
			fld:        5,
			origCoords: &Coordinates{MaxLat: 60.120754, MaxLon: 30.699476, MinLat: 59.805323, MinLon: 30.024240},
			want:       &Coordinates{MaxLat: 60.15234375, MaxLon: 30.7265625, MinLat: 59.80078125, MinLon: 30.0234375},
		},
		{
			name:       "correct curLvl=10 accuracy=-1 fld=2",
			curLvl:     10,
			accuracy:   -1,
			fld:        2,
			origCoords: &Coordinates{MaxLat: 60.120754, MaxLon: 30.699476, MinLat: 59.805323, MinLon: 30.024240},
			want:       &Coordinates{MaxLat: 60.205078125, MaxLon: 30.76171875, MinLat: 59.765625, MinLon: 29.970703125},
		},
		{
			name:       "correct curLvl=10 accuracy=6 fld=5",
			curLvl:     10,
			accuracy:   6,
			fld:        5,
			origCoords: &Coordinates{MaxLat: 60.120754, MaxLon: 30.699476, MinLat: 59.805323, MinLon: 30.024240},
			want:       &Coordinates{MaxLat: 60.1875, MaxLon: 30.9375, MinLat: 59.625, MinLon: 29.8125},
		},
		{
			name:       "correct curLvl=10 accuracy=6 fld=2",
			curLvl:     10,
			accuracy:   6,
			fld:        2,
			origCoords: &Coordinates{MaxLat: 60.120754, MaxLon: 30.699476, MinLat: 59.805323, MinLon: 30.024240},
			want:       &Coordinates{MaxLat: 60.46875, MaxLon: 30.9375, MinLat: 59.0625, MinLon: 29.53125},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cornerCoords := CornerCoordinates(tc.curLvl, tc.accuracy, tc.fld, tc.origCoords)
			assert.Equal(t, tc.want, cornerCoords)
		})
	}

}
