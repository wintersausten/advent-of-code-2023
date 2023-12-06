package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CubeColor string

const (
  Red CubeColor = "red"
  Green CubeColor = "green"
  Blue CubeColor = "blue"
)

type CubeSet map[CubeColor]int

func NewCubeSet(r int, g int, b int) CubeSet {
  return map[CubeColor]int {
    Red: r,
    Green: g,
    Blue: b,
  }
}

var limit = NewCubeSet(12, 13, 14)
func (cs CubeSet) isValid() bool {
  for color, count := range cs {
    if count > limit[color] {
      return false
    }
  }
  return true
}

func parseSet(setString string) CubeSet {
  set := CubeSet{}

  for _, cubeCountSegment := range strings.Split(setString, ",") {
    cubeCountSegment = cubeCountSegment[1:] // trim leading space
    spaceIndex := strings.Index(cubeCountSegment, " ")
    count, _ := strconv.Atoi(cubeCountSegment[:spaceIndex])
    color := CubeColor(cubeCountSegment[spaceIndex + 1:])
    set[color] = count 
  }

  return set
}

func parseGame(line string) (int, []CubeSet) {
  spaceIndex := strings.Index(line, " ")  
  colonIndex := strings.Index(line, ":")  

  gameNum, _ := strconv.Atoi(line[spaceIndex + 1:colonIndex])

  setList := []CubeSet{}
  for _, setString := range strings.Split(line[colonIndex+1:], ";") {
    set := parseSet(setString) 
    setList = append(setList, set)
  }
  
  return gameNum, setList
}

func isValidSetList(setList []CubeSet) bool {
  for _, set := range setList {
    if !set.isValid() {
      return false
    }
  } 
  return true
}

func getMinimumSet(setList []CubeSet) CubeSet {
  minSet := CubeSet{}

  for _, set := range setList {
    for color, count := range set {
      if minSet[color] < count {
        minSet[color] = count
      }
    }
  }

  return minSet
}

// func main() {
func day2() {
  file, err := os.Open("inputs/day2.txt")
  if err != nil {
    fmt.Println(err)
    return
  }

  part1Total := 0
  part2Total := 0

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()

    // parse (should return game # and a list of sets)
    gameNum, setList := parseGame(line)
    // part 1: test list of sets against limit
    if isValidSetList(setList) {
      part1Total += gameNum
    }

    // get minimum of each cube that would support the games sets
    minSet := getMinimumSet(setList)
    product := 1
    for _, count := range minSet {
      product *= count
    }
    part2Total += product
  }

  fmt.Printf("Part 1: %d\n", part1Total)
  fmt.Printf("Part 2: %d\n", part2Total)
  return
}
