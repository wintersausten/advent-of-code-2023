package main

import (
	"bufio"
	"fmt"
	"os"
)
var file, _ = os.Open("./inputs/day10.txt")

type Coord struct {
  x int
  y int
}

func parsePipeLayout(file *os.File) (Coord, [][]rune) {
  var start Coord
  pipeLayout := [][]rune{}
  scanner := bufio.NewScanner(file)

  for y := 0; scanner.Scan(); y++ {
    line := []rune(scanner.Text())
    for x, pipe := range line {
      if pipe == 'S' {
        start = Coord{x, y}
      }
    }
    pipeLayout = append(pipeLayout, line)
  }

  return start, pipeLayout
}

var start, pipeLayout = parsePipeLayout(file)

/*##################################*/

type Offset Coord

var (
  NorthOffset = Offset{0, -1}
  SouthOffset = Offset{0, 1}
  WestOffset = Offset{-1, 0}
  EastOffset = Offset{1, 0}
)

var offsets = []Offset {NorthOffset, SouthOffset, WestOffset, EastOffset}



// is b east relative to a
func (a Coord) isEast(b *Coord) bool {
  return (a.x == b.x+1) && (a.y == b.y)
}

// is b west relative to a
func (a Coord) isWest(b *Coord) bool {
  return (a.x == b.x-1) && (a.y == b.y)
}

// is b North relative to a
func (a Coord) isNorth(b *Coord) bool {
  return (a.x == b.x) && (a.y == b.y-1)
}

// is b south relative to a
func (a Coord) isSouth(b *Coord) bool {
  return (a.x == b.x) && (a.y == b.y+1)
}

func (c *Coord) offset(offset Offset) {
  c.x = c.x + offset.x
  c.y = c.y + offset.y
}

func (a *Coord) toNext(prev *Coord) {
  var offset Offset
  switch pipeLayout[a.y][a.x] {
  case '|':
    if a.isNorth(prev) {
      offset = NorthOffset
    } else {
      offset = SouthOffset
    }
  case '-':
    if a.isWest(prev) {
      offset = WestOffset
    } else {
      offset = EastOffset
    }
  case 'L':
    if a.isWest(prev) {
      offset = NorthOffset
    } else {
      offset = EastOffset
    }
  case 'J':
    if a.isEast(prev) {
      offset = NorthOffset
    } else {
      offset = WestOffset
    }
  case '7':
    if a.isEast(prev) {
      offset = SouthOffset
    } else {
      offset = WestOffset
    }
  case 'F':
    if a.isWest(prev) {
      offset = SouthOffset
    } else {
      offset = EastOffset
    }
  default: // S
    panic("Unexpected path")
  }

  *prev = *a 
  // part 2 
  pipeLayout[a.y][a.x] = ' '
  a.offset(offset)
}
//
//
//
// // is the pipe at a connected to the pipe at b
// // not: does not mean pipe at b is connect to pipe at a
func (a Coord) isConnected(b *Coord) bool {
  switch pipeLayout[a.y][a.x] {
  case '|':
    return a.isNorth(b) || a.isSouth(b)
  case '-':
    return a.isWest(b) || a.isEast(b)
  case 'L':
    return a.isWest(b) || a.isSouth(b)
  case 'J':
    return a.isEast(b) || a.isSouth(b)
  case '7':
    return a.isEast(b) || a.isNorth(b)
  case 'F':
    return a.isWest(b) || a.isNorth(b)
  case '.':
    return false
  default: // S
    return true
  }
}


func main() {
  // determine find possible start
  offsets := []Coord{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}
  itrStartCopy := start
  itr := &itrStartCopy

  for _, offset := range offsets {
    if (itr.x + offset.x) < 0 || (itr.y + offset.y) < 0 {
      continue
    }
    check := Coord{itr.x + offset.x, itr.y + offset.y}
    if check.isConnected(itr) {
      *itr = check
      break
    }
  }


  // assuming that if you come in one side you must go out the other for the main pipe, we know that 
  // the other side must be connected
  startCopy := start
  prev := &startCopy
  // d := 1
  for *itr != start {
    itr.toNext(prev)
    // d++
  } 

  for _, col := 
  
  // result := int(d/2)
  // fmt.Println(result)
}
