package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Lens struct {
  label string
  focalLength int
}

type Box struct {
  lenses []Lens
}

func (b Box) findLensIndex (lens Lens) int{
  for i, l := range b.lenses {
    if l.label == lens.label {
      return i
    }
  }
  return -1
}

func (b *Box) removeLens (lens Lens) {
  i := b.findLensIndex(lens)
  if i != -1 {
    b.lenses = append(b.lenses[:i], b.lenses[i+1:]...)
  }
}

func (b *Box) addLens (lens Lens) {
  i := b.findLensIndex(lens)
  if i != -1 {
    b.lenses[i] = lens
  } else {
    b.lenses = append(b.lenses, lens)
  }
}

func hash(s string) int {
  val := 0
  for _, c := range s {
    code := int(c)
    val += code
    val *= 17
    val %= 256
  }
  return val
}

func calculateFocusingPower(boxes []Box) int {
  power := 0
  for i, box := range boxes {
    for j, lens := range box.lenses {
      power += (i+1) * (j+1) * lens.focalLength
    }
  }
  return power
}

func main() {
  file, _ := os.Open("./inputs/day15.txt")
  seq, _ := io.ReadAll(file)
  // read into 

  steps := strings.Split(string(seq[:len(seq)-1]), ",")
  
  // initialize boxes
  boxes := make([]Box, 256)

  for _, step := range steps {
    if i:=strings.Index(step, "-"); i!=-1 {
      lens := Lens{label: step[:i]}
      boxes[hash(lens.label)].removeLens(lens)
    } else if i:=strings.Index(step, "="); i!=-1 {
      lens := Lens{label: step[:i], focalLength: int(step[i+1] - '0')}
      boxes[hash(lens.label)].addLens(lens)
    } else {
      panic("No instruction!")
    }
  }

  // calculate result
  power := calculateFocusingPower(boxes)

  fmt.Printf("Power: %d\n", power)
}
