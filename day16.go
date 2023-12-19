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

type DirectionTracker map[Direction]struct{}

type Point struct {
  X int
  Y int
}

func (p Point) isValid() bool {
  return p.X >= 0 && p.X < n && p.Y >= 0 && p.Y < m
} 

func (p Point) getPoint(d Direction) Point {
  return Point{p.X + d.X, p.Y + d.Y}
}

var m, n int 

var (
  Up = Direction {X: 0, Y: -1}
  Down = Direction {X: 0, Y: 1}
  Left = Direction {X: -1, Y: 0}
  Right = Direction {X: 1, Y: 0}
)

func parseGrid() [][]rune {
  file, _ := os.Open("./inputs/day16.txt")
  scanner := bufio.NewScanner(file)
  grid := [][]rune{}
  for scanner.Scan() {
    line := scanner.Text()
    grid = append(grid, []rune(line))
  }
  m = len(grid)
  n = len(grid[0])
  return grid
}

func beam(grid [][]rune, energized map[Point]DirectionTracker, point Point, dir Direction) {
  // stop if point is invalid or it's already been traversed in this direction
  // TODO energized shoudl be a map of directions as val
  if !point.isValid() {
    return
  }
  if prevDir, exists := energized[point]; exists {
    if _, exists := prevDir[dir]; exists {
      return
    }
  } else {
    energized[point] = DirectionTracker{}
  }

  energized[point][dir] = struct{}{}

  switch grid[point.Y][point.X] {
  case '|':
    if (dir == Left || dir == Right) {
      beam(grid, energized, point.getPoint(Up), Up)
      beam(grid, energized, point.getPoint(Down), Down)
    } else {
      beam(grid, energized, point.getPoint(dir), dir)
    }
  case '-':
    if (dir == Up || dir == Down) {
      beam(grid, energized, point.getPoint(Left), Left)
      beam(grid, energized, point.getPoint(Right), Right)
    } else {
      beam(grid, energized, point.getPoint(dir), dir)
    }
  case '/':
    switch dir {
    case Up:
      beam(grid, energized, point.getPoint(Right), Right)
    case Down:
      beam(grid, energized, point.getPoint(Left), Left)
    case Left:
      beam(grid, energized, point.getPoint(Down), Down)
    case Right:
      beam(grid, energized, point.getPoint(Up), Up)
    }
  case '\\':
    switch dir {
    case Up:
      beam(grid, energized, point.getPoint(Left), Left)
    case Down:
      beam(grid, energized, point.getPoint(Right), Right)
    case Left:
      beam(grid, energized, point.getPoint(Up), Up)
    case Right:
      beam(grid, energized, point.getPoint(Down), Down)
    }
  default:
    beam(grid, energized, point.getPoint(dir), dir)
  }
}

func visualizeEnergy(energized map[Point]DirectionTracker) {
  for i := 0; i < m; i++ {
    line := []rune{}
    for j := 0; j < n; j++ {
      if _, exists := energized[Point{j, i}]; exists{
        line = append(line, '#')
      } else {
        line = append(line, '.')
      }
    }
    fmt.Println(string(line))
  }
}

func main() {
  grid := parseGrid() 

  max := 0
  for y := 0; y < m; y++ {
    energized := map[Point]DirectionTracker{}
    beam(grid, energized, Point{0, y}, Right)
    if len(energized) > max {
      max = len(energized)
    }

    energized = map[Point]DirectionTracker{}
    beam(grid, energized, Point{n-1, y}, Left)
    if len(energized) > max {
      max = len(energized)
    }
  }

  for x := 0; x < m; x++ {
    energized := map[Point]DirectionTracker{}
    beam(grid, energized, Point{x, 0}, Down)
    if len(energized) > max {
      max = len(energized)
    }

    energized = map[Point]DirectionTracker{}
    beam(grid, energized, Point{x, m-1}, Up)
    if len(energized) > max {
      max = len(energized)
    }
  }

  fmt.Printf("Total: %d\n", max)
  // visualizeEnergy(energized)
}
