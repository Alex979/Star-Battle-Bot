package main

import (
	"fmt"
	"time"
)

const GridSize = 8
const NumStars = 1

func printGrid(grid [][]uint, rowStars []uint, columnStars []uint, shapeStars []uint) {
	for _, row := range grid {
		for _, val := range row {
			switch val {
			case 0:
				fmt.Print("  ")
			case 1:
				fmt.Print("★ ")
			case 2:
				fmt.Print("∙ ")
			}
		}
		fmt.Println()
	}
}

func copyGrid(grid [][]uint, rowStars []uint, columnStars []uint, shapeStars []uint) ([][]uint, []uint, []uint, []uint) {
	// Copy the grid values
	newGrid := make([][]uint, len(grid))
	for i := range newGrid {
		newGrid[i] = make([]uint, GridSize)
		copy(newGrid[i], grid[i])
	}
	// Copy the remaining slices
	newRowStars := make([]uint, len(rowStars))
	copy(newRowStars, rowStars)
	newColumnStars := make([]uint, len(columnStars))
	copy(newColumnStars, columnStars)
	newShapeStars := make([]uint, len(shapeStars))
	copy(newShapeStars, shapeStars)

	return newGrid, newRowStars, newColumnStars, newShapeStars
}

func placeStar(grid [][]uint, shapes [][]uint, shapeHash map[uint][]int, rowStars []uint, columnStars []uint, shapeStars []uint, i int, j int) {
	// Set grid to star
	grid[i][j] = 1

	// Update row and column totals
	rowStars[i] += 1
	columnStars[j] += 1

	// Update shape total
	shapeIndex := shapes[i][j]
	shapeStars[shapeIndex] += 1

	// Remove spots around the star
	for y := i - 1; y <= i+1; y++ {
		for x := j - 1; x <= j+1; x++ {
			if y >= 0 && y < len(grid) && x >= 0 && x < len(grid[y]) && grid[y][x] == 0 {
				grid[y][x] = 2
			}
		}
	}

	// Check if row or column is full, and block them out if so
	if rowStars[i] == NumStars {
		for x := range grid {
			if grid[i][x] == 0 {
				grid[i][x] = 2
			}
		}
	}
	if columnStars[j] == NumStars {
		for y := range grid {
			if grid[y][j] == 0 {
				grid[y][j] = 2
			}
		}
	}

	// Check if shape is full, and if so block out all it's squares
	if shapeStars[shapeIndex] == NumStars {
		coords := shapeHash[shapeIndex]
		for _, flattenedCoord := range coords {
			y := flattenedCoord / GridSize
			x := flattenedCoord % GridSize
			if grid[y][x] == 0 {
				grid[y][x] = 2
			}
		}
	}
}

func isPossible(grid [][]uint, shapes [][]uint, shapeHash map[uint][]int) bool {
	// Check if any rows cannot contain NumStars stars
	for i := range grid {
		numPossibleStars := 0
		for j := range grid[i] {
			if grid[i][j] != 2 {
				numPossibleStars++
			}
		}
		if numPossibleStars < NumStars {
			return false
		}
	}

	// Check if any columns cannot contain NumStars stars
	for j := range grid[0] {
		numPossibleStars := 0
		for i := range grid {
			if grid[i][j] != 2 {
				numPossibleStars++
			}
		}
		if numPossibleStars < NumStars {
			return false
		}
	}

	// Check if any shapes cannot contain NumStars stars
	// Different logic for 1 vs 2 stars
	switch NumStars {
	case 2:
		for _, coords := range shapeHash {
			maxDist := 0
			for _, coord1 := range coords {
				for _, coord2 := range coords {
					i1 := coord1 / GridSize
					j1 := coord1 % GridSize
					i2 := coord2 / GridSize
					j2 := coord2 % GridSize

					if grid[i1][j1] != 2 && grid[i2][j2] != 2 {
						xDist := i1 - i2
						if xDist < 0 {
							xDist = -xDist
						}
						yDist := j1 - j2
						if yDist < 0 {
							yDist = -yDist
						}

						if xDist > maxDist {
							maxDist = xDist
						}
						if yDist > maxDist {
							maxDist = yDist
						}
					}
				}
			}

			if maxDist < 2 {
				return false
			}
		}
	case 1:
		for _, coords := range shapeHash {
			numPossibleStars := 0
			for _, coord := range coords {
				i := coord / GridSize
				j := coord % GridSize
				if grid[i][j] != 2 {
					numPossibleStars++
				}
			}
		}
	}

	return true
}

func generateGrid(grid [][]uint, shapes [][]uint, shapeHash map[uint][]int, rowStars []uint, columnStars []uint, shapeStars []uint, totalStars uint) bool {
	if totalStars == GridSize*NumStars {
		printGrid(grid, rowStars, columnStars, shapeStars)
		return true
	}

	// Iterate through possible star locations
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 0 {
				newGrid, newRowStars, newColumnStars, newShapeStars := copyGrid(grid, rowStars, columnStars, shapeStars)
				placeStar(newGrid, shapes, shapeHash, newRowStars, newColumnStars, newShapeStars, i, j)

				if isPossible(newGrid, shapes, shapeHash) {
					success := generateGrid(newGrid, shapes, shapeHash, newRowStars, newColumnStars, newShapeStars, totalStars+1)
					if success {
						return true
					}
				}
			}
		}
	}

	return false
}

func main() {
	grid := make([][]uint, GridSize)
	for i := range grid {
		grid[i] = make([]uint, GridSize)
	}
	rowStars := make([]uint, GridSize)
	columnStars := make([]uint, GridSize)
	shapeStars := make([]uint, GridSize)
	shapes := [][]uint{
		{0, 0, 1, 2, 3, 4, 4, 4},
		{0, 0, 1, 2, 3, 3, 3, 3},
		{0, 0, 2, 2, 3, 3, 3, 3},
		{0, 0, 0, 2, 2, 3, 3, 3},
		{5, 5, 5, 5, 6, 6, 3, 3},
		{5, 5, 5, 5, 6, 6, 3, 3},
		{7, 7, 5, 5, 7, 7, 7, 7},
		{7, 7, 7, 7, 7, 7, 7, 7},
	}
	// shapes := [][]uint{
	// 	{0, 0, 0, 1, 1, 1, 2, 2, 2, 2},
	// 	{0, 0, 0, 1, 1, 1, 1, 1, 1, 2},
	// 	{3, 3, 0, 0, 2, 2, 2, 2, 2, 2},
	// 	{3, 3, 3, 4, 4, 4, 2, 5, 5, 5},
	// 	{3, 3, 3, 4, 6, 4, 2, 2, 5, 5},
	// 	{3, 3, 4, 4, 6, 6, 6, 6, 6, 6},
	// 	{3, 3, 4, 4, 4, 7, 7, 7, 8, 8},
	// 	{3, 3, 4, 8, 8, 8, 8, 8, 8, 8},
	// 	{8, 8, 8, 8, 8, 8, 9, 9, 8, 8},
	// 	{8, 8, 8, 8, 8, 8, 9, 9, 9, 8},
	// }

	shapeHash := make(map[uint][]int)
	for i := range shapes {
		for j, shape := range shapes[i] {
			flattendCoord := i*GridSize + j
			if _, ok := shapeHash[shape]; ok {
				shapeHash[shape] = append(shapeHash[shape], flattendCoord)
			} else {
				shapeHash[shape] = []int{flattendCoord}
			}
		}
	}

	start := time.Now()
	generateGrid(grid, shapes, shapeHash, rowStars, columnStars, shapeStars, 0)
	end := time.Now()

	elapsed := end.Sub(start)
	fmt.Printf("Runtime: %v\n", elapsed)
}
