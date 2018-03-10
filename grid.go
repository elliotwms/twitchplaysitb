package main

const tileWidth = 112
const tileHeight = 82
const boardWidth = tileWidth * 8
const boardHeight = tileHeight * 8

// GetCoordinates takes a set of A-H 1-8 grid coordinates and returns pixel coordinates
func GetCoordinates(a, n string) (x, y int) {
	x = gameWidth / 2
	y = gameHeight - 12 - (tileHeight / 2) // the bottom tip of the bottom tile is 12 pixels up from the bottom

	as := []string{"A", "B", "C", "D", "E", "F", "G", "H"} // Alpha grid coordinates
	ns := []string{"1", "2", "3", "4", "5", "6", "7", "8"} // Numeric grid coordinates

	for k, v := range as {
		if v == a {
			x += (tileWidth / 2) * k  // Move right half a tile for every row
			y -= (tileHeight / 2) * k // Move up half a tile for every row
		}
	}

	for k, v := range ns {
		if v == n {
			x -= (tileWidth / 2) * k  // Move left half a tile for every row
			y -= (tileHeight / 2) * k // Move up half a tile for every row
		}
	}

	return x, y
}
