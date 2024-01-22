package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FloatPoint struct {
  X float64
  Y float64
}

type FloatVector struct {
  X float64
  Y float64
}

type Hail struct {
  FloatPoint
  FloatVector
  m float64
  b float64
}

var min, max = 200000000000000, 400000000000000
// var min, max = 7, 27
func (p FloatPoint) inBounds() bool {
  return p.X >= float64(min) && p.X <= float64(max) && p.Y >= float64(min) && p.Y <= float64(max)
}

func (p FloatPoint) inFuture(h Hail) bool {
  return ((p.X - h.FloatPoint.X) * h.FloatVector.X) >= 0 && ((p.Y - h.FloatPoint.Y) * h.FloatVector.Y) >= 0
}

func parseHail() []Hail{
  file, _ := os.Open("./inputs/day24.txt")
  scanner := bufio.NewScanner(file)
  hail := []Hail{}
  for scanner.Scan() {
    line := scanner.Text()
    sides := strings.Split(line, "@")

    pointStrings := strings.Split(sides[0], ",")
    x, _ := strconv.ParseFloat(strings.TrimSpace(pointStrings[0]), 64)
    y, _ := strconv.ParseFloat(strings.TrimSpace(pointStrings[1]), 64)
    point := FloatPoint{X: x, Y: y}

    vectorStrings := strings.Split(sides[1], ",")
    dx, _ := strconv.ParseFloat(strings.TrimSpace(vectorStrings[0]), 64)
    dy, _ := strconv.ParseFloat(strings.TrimSpace(vectorStrings[1]), 64)
    vector := FloatVector{X: dx, Y: dy}
    m := dy / dx
    b := y - (m * x)

    hail = append(hail, Hail{point, vector, m, b})
  }
  return hail
}

func getInterctionCount(hail []Hail) int {
  count := 0
  for i := 0; i < len(hail)-1; i++ {
    for j := i+1; j < len(hail); j++ {
      x := (hail[j].b - hail[i].b) / (hail[i].m - hail[j].m)
      y := (hail[i].m * x) + hail[i].b
      intersection := FloatPoint{x, y}
      if intersection.inBounds() && intersection.inFuture(hail[i]) && intersection.inFuture(hail[j]){
        count++
      }
    }
  }
  return count
}

func main() {
  hail := parseHail()

  fmt.Printf("Total: %d\n", getInterctionCount(hail))

}
