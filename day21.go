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

var (
  Up = Direction {X: 0, Y: -1}
  Down = Direction {X: 0, Y: 1}
  Left = Direction {X: -1, Y: 0}
  Right = Direction {X: 1, Y: 0}
)

var m, n int

func (p Point) isValid() bool {
  return p.X >= 0 && p.X < n && p.Y >= 0 && p.Y < m
} 

var directions = []Direction{Up, Down, Left, Right}

type Point struct {
  X int
  Y int
}

type Mem struct {
  Point
  remStep int
}

func copyEndPos(endPos map[Point]struct{}) map[Point]struct{} {
  copied := make(map[Point]struct{})
  for key, val := range endPos {
    copied[key] = val
  }
  return copied
}

func union(sets []map[Point]struct{}) map[Point]struct{} {
  unionSet := make(map[Point]struct{})
  for _, set := range sets {
    for key := range set {
      unionSet[key] = struct{}{}
    }
  }
  return unionSet
}

func walk(garden [][]rune, pos Point, remStep int, endPos map[Point]struct{}, memMap map[Mem]map[Point]struct{}) map[Point]struct{} {
  // Check all four directions
  if remStep == 0 {
    endPos[pos] = struct{}{}
    return endPos
  } else if memEndPos, exists := memMap[Mem{pos, remStep}]; exists {
    return memEndPos
  }

  allEndPos := []map[Point]struct{}{}
  for _, dir := range directions {
    nextPoint := Point{pos.X + dir.X, pos.Y + dir.Y}
    if nextPoint.isValid() && garden[nextPoint.Y][nextPoint.X] != '#' {
      // supply a new endPos set, bootstrapped from current
      allEndPos = append(allEndPos, walk(garden, nextPoint, remStep-1, copyEndPos(endPos), memMap))
    }
  }

  // sum all end pos
  endPos = union(allEndPos)
  memMap[Mem{pos, remStep}] = endPos
  // add to mem map
  return endPos
}

func parseGarden() ([][]rune, Point) {
  file, _ := os.Open("./inputs/day21.txt")
  scanner := bufio.NewScanner(file)
  garden := [][]rune{}
  start := Point{}
  i := 0
  for scanner.Scan() {
    row := []rune{}
    for j, c := range scanner.Text() {
      if c == 'S' {
        c = '.'
        start.X = j
        start.Y = i
      }
      row = append(row, c)
    }
    garden = append(garden, row)
    i++
  }
  m = len(garden)
  n = len(garden[0])
  return garden, start
}

func main() {
  garden, start := parseGarden()
  memMap := map[Mem]map[Point]struct{}{}
  endPos := walk(garden, start, 64, map[Point]struct{}{}, memMap)

  fmt.Printf("Pos: %d\n", len(endPos))
}
