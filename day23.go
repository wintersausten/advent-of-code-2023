package main

import (
	"bufio"
	"fmt"
	"os"
)

type Direction struct {
  X int
  Y int
}

type Point struct {
  X int
  Y int
}

var (
  Up = Direction {X: 0, Y: -1}
  Down = Direction {X: 0, Y: 1}
  Left = Direction {X: -1, Y: 0}
  Right = Direction {X: 1, Y: 0}
)

var directions = []Direction{Up, Down, Left, Right}

var m, n int
func (p Point) isValid() bool {
  return p.X >= 0 && p.X < n && p.Y >= 0 && p.Y < m
} 

var exit Point
func parseForest() [][]rune{
  file, _ := os.Open("./inputs/day23.txt")
  scanner := bufio.NewScanner(file)
  forest := [][]rune{}
  for scanner.Scan() {
    forest = append(forest, []rune(scanner.Text()))
  }
  m = len(forest)
  n = len(forest[0])
  exit = Point{X: n-2, Y: m-1}
  return forest
}

func findLongestPath(forest [][]rune) int {
  return travel(forest, Point{X: 1, Y: 0}, 0, map[Point]struct{}{}, -1)
}

func getAllowedDirections(r rune) []Direction{
  switch r {
  case '>':
    return []Direction{Right}
  case '<':
    return []Direction{Left}
  case '^':
    return []Direction{Up}
  case 'v':
    return []Direction{Down}
  default:
    return directions
  }
}

func travel(forest [][]rune, itr Point, distance int, seen map[Point]struct{}, max int) int {
  if itr == exit {
    return distance
  }

  seen[itr] = struct{}{}

  // allowedDirections := getAllowedDirections(forest[itr.Y][itr.X])
  allowedDirections := directions

  for _, dir := range allowedDirections {
    next := Point{X: itr.X + dir.X, Y: itr.Y + dir.Y}
    if _, exists := seen[next]; (!exists && next.isValid() && forest[next.Y][next.X] != '#') {
      subMax := travel(forest, next, distance + 1, seen, max)
      if subMax > max {
        max = subMax
      }
    }
  }

  delete(seen, itr)
  return max
}

func main() {
  forest := parseForest()

  fmt.Printf("Longest: %d\n", findLongestPath(forest))
}
