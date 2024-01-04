package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Orientation int

const (
  X Orientation = 0
  Y Orientation = 1
  Z Orientation = 2
)

type Brick struct {
  Orientation
  start Point
  end Point
  above []*Brick
  below []*Brick
}

func NewBrick(start Point, end Point, o Orientation) *Brick {
  return &Brick{
    Orientation: o,
    start: start,
    end: end,
    above: []*Brick{},
    below: []*Brick{},
  }
}

type Snapshot [][][]*Brick

func NewSnapshot(x int, y int, z int) Snapshot {
  snapshot := Snapshot{}
  for i := 0; i <= x; i++ {
    layer := [][]*Brick{}
    for j := 0; j <= y; j++ {
      row := make([]*Brick, z+1)
      layer = append(layer, row)
    }
    snapshot = append(snapshot, layer)
  }
  return snapshot
}

type Point struct {
  X int
  Y int
  Z int
}

func relate(above *Brick, below *Brick) {
  above.below = append(above.below, below)
  below.above = append(below.above, above)
}

func (s Snapshot) remove(b *Brick) {
  switch b.Orientation {
  case X:
    y, z := b.start.Y, b.start.Z
    for x := b.start.X; x <= b.end.X; x++ {
      s[x][y][z] = nil
    }
  case Y:
    x, z := b.start.X, b.start.Z
    for y := b.start.Y; y <= b.end.Y; y++ {
      s[x][y][z] = nil
    }
  default: 
    x, y := b.start.X, b.start.Y
    for z := b.start.Z; z <= b.end.Z; z++ {
      s[x][y][z] = nil
    }
  }
}

func (s Snapshot) place(b *Brick) {
  switch b.Orientation {
  case X:
    y, z := b.start.Y, b.start.Z
    for x := b.start.X; x <= b.end.X; x++ {
      s[x][y][z] = b
    }
  case Y:
    x, z := b.start.X, b.start.Z
    for y := b.start.Y; y <= b.end.Y; y++ {
      s[x][y][z] = b
    }
  default: 
    x, y := b.start.X, b.start.Y
    for z := b.start.Z; z <= b.end.Z; z++ {
      s[x][y][z] = b
    }
  }
}

// drop brick by checking one row down Z until there's at least one brick or it hits ground
func (s Snapshot) drop(b *Brick) {
  below := map[*Brick]struct{}{}
  switch b.Orientation {
  case X:
    y := b.start.Y
    z := b.start.Z - 1
    for z >= 1 && len(below) == 0 {
      for x := b.start.X; x <= b.end.X; x++ {
        if _, exists := below[s[x][y][z]]; !exists && s[x][y][z] != nil {
          below[s[x][y][z]] = struct{}{}
        }
      }
      if len(below) == 0 {
        z--
      }
    }
    z++ // move up - block must rest on whatever is blocking below
    if z != b.start.Z {
      s.remove(b)
      b.start.Z = z
      b.end.Z = z
      s.place(b)
    }
    for blockBelow, _ := range below {
      relate(b, blockBelow)
    }
  case Y:
    x := b.start.X
    z := b.start.Z - 1
    for z >= 1 && len(below) == 0 {
      for y := b.start.Y; y <= b.end.Y; y++ {
        if _, exists := below[s[x][y][z]]; !exists && s[x][y][z] != nil {
          below[s[x][y][z]] = struct{}{}
        }
      }
      if len(below) == 0 {
        z--
      }
    }
    z++
    if z != b.start.Z {
      s.remove(b)
      b.start.Z = z
      b.end.Z = z
      s.place(b)
    }
    for blockBelow, _ := range below {
      relate(b, blockBelow)
    }
  default:
    x := b.start.X
    y := b.start.Y
    z := b.start.Z - 1
    for z >= 1 && len(below) == 0 {
      if _, exists := below[s[x][y][z]]; !exists && s[x][y][z] != nil {
        below[s[x][y][z]] = struct{}{}
      }
      if len(below) == 0 {
        z--
      }
    }
    z++
    // move block
    if z != b.start.Z {
      s.remove(b)
      zDiff := b.start.Z - z
      b.start.Z = b.start.Z - zDiff
      b.end.Z = b.end.Z - zDiff
      s.place(b)
    }
    // relate to blocks below
    for blockBelow, _ := range below {
      relate(b, blockBelow)
    }
  }
}

// resolve snapshot
// resolve from the ground up
func (s Snapshot) resolve() {
  resolved := map[*Brick]struct{}{}
  m, n, o := len(s), len(s[0]), len(s[0][0])

  for z := 1; z < o; z++ {
    for y := 0; y < n; y++ {
      for x := 0; x < m; x++ {
        // if the brick is not in resolved & isn't nil
        if _, exists := resolved[s[x][y][z]]; !exists && s[x][y][z] != nil {
          s.drop(s[x][y][z])
          resolved[s[x][y][z]] = struct{}{}
        }
      }
    }
  }
}

// parse snapshot
func parsePoint(side string) Point {
  coords := strings.Split(side, ",")
  x, _ := strconv.Atoi(coords[0])
  y, _ := strconv.Atoi(coords[1])
  z, _ := strconv.Atoi(coords[2])
  return Point{x, y, z}
}

func parseBrick(line string) *Brick{
  sides := strings.Split(line, "~")
  a := parsePoint(sides[0])
  b := parsePoint(sides[1])

  dx, dy, dz := a.X - b.X, a.Y - b.Y, a.Z - b.Z
  if dx != 0 {
    if dx < 0 {
      return NewBrick(a, b, X)
    } else {
      return NewBrick(b, a, X)
    }
  } else if dy != 0 {
    if dy < 0 {
      return NewBrick(a, b, Y)
    } else {
      return NewBrick(b, a, Y)
    }
  } else {
    if dz < 0 {
      return NewBrick(a, b, Z)
    } else {
      return NewBrick(b, a, Z)
    }
  }
}

func parseSnapshot() (Snapshot, []*Brick) {
  file, _ := os.Open("./inputs/day22.txt")
  scanner := bufio.NewScanner(file)
  snapshot := NewSnapshot(9, 9, 303)
  bricks := []*Brick{}
  // mx, my, mz := 0, 0, 0
  for scanner.Scan() {
    line := scanner.Text()
    brick := parseBrick(line)
    // if brick.end.X > mx {
    //   mx = brick.end.X
    // }
    // if brick.end.Y > my {
    //   my = brick.end.Y
    // }
    // if brick.end.Z > mz {
    //   mz = brick.end.Z
    // }
    snapshot.place(brick)
    bricks = append(bricks, brick)
  }
  // fmt.Printf("%d, %d, %d", mx, my, mz)
  return snapshot, bricks
}

func countRemovableBricks(bricks []*Brick) int {
  count := 0
  BrickLoop:
  for _, b := range bricks {
    //if for evern brick above, current isn't the only brick below
    for _, above := range b.above {
      if len(above.below) == 1 && above.below[0] == b {
        continue BrickLoop
      }
    }
    count ++
  }
  return count
}

func main() {
  // parse 
  snapshot, bricks := parseSnapshot()

  // resolve snapshot
  snapshot.resolve() 
  
  fmt.Printf("Total: %d\n", countRemovableBricks(bricks))
}
