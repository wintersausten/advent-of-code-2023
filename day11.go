package main

import (
	"bufio"
	"fmt"
	"os"
)


type Point struct {
  X int 
  Y int
}

var (
  up = Point{-1, 0}
  down = Point{1, 0}
  right = Point{0, 1}
  left = Point{0, -1}
)


// insertCopies inserts copies of a value into a slice at specified indices.
func insertCopies[T any](slice []T, value T, indices []int) []T {
    result := make([]T, len(slice) + len(indices))
    resultIdx, sliceIdx := 0, 0

    for _, index := range indices {
        // Copy elements up to the current index
        for ; sliceIdx < index; sliceIdx++ {
            result[resultIdx] = slice[sliceIdx]
            resultIdx++
        }
        // Insert the value
        result[resultIdx] = value
        resultIdx++
    }
    // Copy any remaining elements
    copy(result[resultIdx:], slice[sliceIdx:])
    return result
}

func expandSpace(space *[][]rune) (map[int]struct{}, map[int]struct{}){
  // expand rows
  expandRows := map[int]struct{}{}
  for i := 0; i < len(*space); i++ {
    expand := true
    for _, val := range (*space)[i] {
      if val != '.' {
        expand = false
        break
      }
    }
    if expand {
      expandRows[i] = struct{}{}
    }
  }

  // expand cols
  expandCols := map[int]struct{}{}
  for i := 0; i < len((*space)[0]); i++ {
    expand := true
    for _, row := range *space {
      val := row[i]
      if val != '.' {
        expand = false
        break
      }
    }
    if expand {
      expandCols[i] = struct{}{}
    }
  }

  return expandCols, expandRows
}

// Check if the point is within the space and not an obstacle
func isValid(space [][]rune, visited map[Point]Point, point Point) bool {
    rows, cols := len(space), len(space[0])
    _, beenVisited := visited[point]
    return point.X >= 0 && point.X < cols && point.Y >= 0 && point.Y < rows && !beenVisited
}

func oppositeDirection (dir Point) Point {
  switch dir {
  case up:
    return down
  case down:
    return up
  case right:
    return left
  case left:
    return right
  default:
    panic("Tried to get opposite of not-direction")
  }
}

func distanceBack (space [][]rune, visited map[Point]Point, end Point, start Point, expandCols map[int]struct{}, expandRows map[int]struct{}) int {
  d := 0
  itr := end
  crossedCols := map[int]struct{}{}
  crossedRows := map[int]struct{}{}
  for itr != start {
    if _, expandedCol := expandCols[itr.X]; expandedCol {
      if _, crossed := crossedCols[itr.X]; !crossed {
        d += 999999
        crossedCols[itr.X] = struct{}{}
      }
    }
    if _, expandedRow := expandRows[itr.Y]; expandedRow {
       if _, crossed := crossedRows[itr.Y]; !crossed {
         d += 999999
         crossedRows[itr.Y] = struct{}{}
       }
    }
    d++

    dir, _ := visited[itr]
    itr.X += dir.X
    itr.Y += dir.Y
  }
  return d
}

func bfs(space [][]rune, start Point, galaxiesProcessed map[rune]struct{}, total *int, expandCols map[int]struct{}, expandRows map[int]struct{}) {
  // Directions: up, down, left, right
  directions := []Point{up, down, left, right}

  // map of point in space to direction 
  visited := map[Point]Point{}

  queue := []Point{start}
  visited[start] = Point{0, 0}

  distance := -1
  // can speed up by stopping when all necessary galaxies are calculated
  for len(queue) > 0 {
    distance++
    layerSize := len(queue)
    for i := 0; i < layerSize; i++ {
      point := queue[0]
      queue = queue[1:]

      // If point is unprocessed galaxy
      if _, exists := galaxiesProcessed[space[point.Y][point.X]]; distance != 0 && !exists && space[point.Y][point.X] != '.' {
        *total += distanceBack(space, visited, point, start, expandCols, expandRows)
      }

      // Check all four directions
      for _, dir := range directions {
          nextPoint := Point{point.X + dir.X, point.Y + dir.Y}
          if isValid(space, visited, nextPoint) {
              queue = append(queue, nextPoint)
              visited[nextPoint] = oppositeDirection(dir)
          }
      }
    }
  }
}

func totalDistance(space [][]rune, galaxies []Point, expandCols map[int]struct{}, expandRows map[int]struct{}) int {
  var total int
  
  // go through every galaxy
  galaxiesProcessed := map[rune]struct{}{}
  for _, galaxyPos := range galaxies {
    bfs(space, galaxyPos, galaxiesProcessed, &total, expandCols, expandRows)
    galaxiesProcessed[space[galaxyPos.Y][galaxyPos.X]] = struct{}{}
  }

  return total
}

func findGalaxies(space [][]rune) []Point {
  galaxies := []Point{}
  for Y, row := range space {
    for X, val := range row {
      if val == '#' {
        space[Y][X] = rune(len(galaxies) + 1) + '0'
        galaxies = append(galaxies, Point{X, Y})
      }
    }
  }
  return galaxies
}

func parseSpace() ([][]rune) {
  var file, _ = os.Open("./inputs/day11.txt")
  var scanner = bufio.NewScanner(file)

  space := [][]rune{}
  for row := 0; scanner.Scan(); row++{
    line := scanner.Text()
    space = append(space, []rune(line)) 
  }
  return space 
}

func main() {
  space := parseSpace()
  expandCols, expandRows := expandSpace(&space)
  galaxies := findGalaxies(space)

  total := totalDistance(space, galaxies, expandCols, expandRows)
 
  fmt.Printf("Total: %d\n", total)
}
