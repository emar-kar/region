// Package region implements functions to calculate tiles range according to the given coordinates.
package region

import (
	"math"
)

// Coordinates represents map coordinates.
type Coordinates struct {
	MaxLat float64
	MaxLon float64
	MinLat float64
	MinLon float64
}

// Tiles represents data about current tiles map.
type Tiles struct {
	Level       int
	Range       *TilesRange
	Coordinates *Coordinates
}

// TilesRange represents range of the tiles.
type TilesRange struct {
	MaxY int
	MaxX int
	MinY int
	MinX int
}

// CornerCoordinates returns coordinates of the nearest tile corner.
func CornerCoordinates(curLvl, accuracy int, fld float64, originalCoordinates *Coordinates) *Coordinates {
	absCoords := Coordinates{}

	if accuracy != -1 {
		curLvl = accuracy
	}

	rangeX := 360 / ((fld * (math.Pow(2, float64(curLvl)))) * 2)
	rangeY := 180 / (fld * (math.Pow(2, float64(curLvl))))

	absCoords.MaxLat = math.Ceil((originalCoordinates.MaxLat+90)/rangeY)*rangeY - 90
	absCoords.MaxLon = math.Ceil((originalCoordinates.MaxLon+180)/rangeX)*rangeX - 180
	absCoords.MinLat = math.Floor((originalCoordinates.MinLat+90)/rangeY)*rangeY - 90
	absCoords.MinLon = math.Floor((originalCoordinates.MinLon+180)/rangeX)*rangeX - 180

	return &absCoords
}

// GetTiles returns tiles range for given parameters.
func GetTiles(lvl, accuracy int, fld float64, coords *Coordinates) *Tiles {
	cornerCoords := CornerCoordinates(lvl, accuracy, fld, coords)

	rangeX := 360 / ((fld * (math.Pow(2, float64(lvl)))) * 2)
	rangeY := 180 / (fld * (math.Pow(2, float64(lvl))))

	tilesRange := TilesRange{}

	tilesRange.MaxY = int((cornerCoords.MaxLat+90)/rangeY) - 1
	tilesRange.MaxX = int((cornerCoords.MaxLon+180)/rangeX) - 1
	tilesRange.MinY = int((cornerCoords.MinLat + 90) / rangeY)
	tilesRange.MinX = int((cornerCoords.MinLon + 180) / rangeX)

	return &Tiles{Level: lvl, Range: &tilesRange, Coordinates: cornerCoords}
}
