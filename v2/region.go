// Package utils implements functions to calculate tiles range according to the given coordinates.
package region

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// ErrLvlRange raises in case if one of the levels is out of range 0-16.
var ErrLvlRange = errors.New("levels range is broken")

// ErrLvlParse raises in case if the level string is incorrect.
var ErrLvlParse = errors.New("cannot parse levels")

// Coordinates represents map with coordinates.
type Coordinates map[string]float64

type Tiles struct {
	Level       int
	Range       *TilesRange
	Coordinates Coordinates
}

type TilesRange struct {
	MaxLat int
	MaxLon int
	MinLat int
	MinLon int
}

// AbsCoordinates returns coordinates of the nearest tile corner.
func AbsCoordinates(curLvl, accuracy int, fld float64, originalCoordinates Coordinates) Coordinates {
	absCoords := Coordinates{}

	if accuracy != -1 {
		curLvl = accuracy
	}

	rangeX := 360 / ((fld * (math.Pow(2, float64(curLvl)))) * 2)
	rangeY := 180 / (fld * (math.Pow(2, float64(curLvl))))

	absCoords["minLat"] = math.Floor((originalCoordinates["minLat"]+90)/rangeY)*rangeY - 90
	absCoords["minLon"] = math.Floor((originalCoordinates["minLon"]+180)/rangeX)*rangeX - 180
	absCoords["maxLat"] = math.Ceil((originalCoordinates["maxLat"]+90)/rangeY)*rangeY - 90
	absCoords["maxLon"] = math.Ceil((originalCoordinates["maxLon"]+180)/rangeX)*rangeX - 180

	return absCoords
}

// GetCoordinates parses coordinates from the given parameters.
func GetCoordinates(args []string) (*Coordinates, error) {
	var coordinatesVars = [4]string{"minLat", "minLon", "maxLat", "maxLon"}
	coords := Coordinates{}

	for ind, el := range coordinatesVars {
		current, err := strconv.ParseFloat(RemoveCharacters(args[ind], ", "), 64)
		if err != nil {
			return nil, err
		}

		coords[el] = current
	}
	return &coords, nil
}

// GetLvls parses levels from parameters.
func GetLvls(s string) (int, int, error) {
	lvlsString := strings.Split(RemoveCharacters(s, " "), "-")
	if len(lvlsString) != 2 {
		return 0, 0, ErrLvlParse
	}

	fromLvl, err := strconv.Atoi(lvlsString[0])
	if err != nil {
		return 0, 0, err
	}

	toLvl, err := strconv.Atoi(lvlsString[1])
	if err != nil {
		return 0, 0, err
	}

	if toLvl > 16 || fromLvl < 0 || toLvl < fromLvl {
		return 0, 0, ErrLvlRange
	}

	return fromLvl, toLvl, nil
}

// GetTiles returns tiles range for given parameters.
func GetTiles(lvl, accuracy int, fld float64, coords *Coordinates) *Tiles {
	absCoords := AbsCoordinates(lvl, accuracy, fld, *coords)

	rangeX := 360 / ((fld * (math.Pow(2, float64(lvl)))) * 2)
	rangeY := 180 / (fld * (math.Pow(2, float64(lvl))))

	tilesRange := TilesRange{}

	tilesRange.MinLat = int((absCoords["minLat"] + 90) / rangeY)
	tilesRange.MinLon = int((absCoords["minLon"] + 180) / rangeX)
	tilesRange.MaxLat = int((absCoords["maxLat"]+90)/rangeY) - 1
	tilesRange.MaxLon = int((absCoords["maxLon"]+180)/rangeX) - 1

	return &Tiles{Level: lvl, Range: &tilesRange, Coordinates: absCoords}

	// fmt.Printf("For level %d\n", lvl)
	// fmt.Printf("y range = %d - %d, x range = %d - %d\n", minLat, maxLat-1, minLon, maxLon-1)
	// fmt.Printf("Min coordinates: %v, %v\n", absCoords["minLat"], absCoords["minLon"])
	// fmt.Printf("Max coordinates: %v, %v\n", absCoords["maxLat"], absCoords["maxLon"])
}

// RemoveCharacters removes determined characters from the given string.
func RemoveCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if !strings.ContainsRune(characters, r) {
			return r
		}
		return -1
	}
	return strings.Map(filter, input)
}
