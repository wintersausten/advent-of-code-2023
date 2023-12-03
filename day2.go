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

func main() {
  file, err := os.Open("inputs/day2.txt")
  if err != nil {
    fmt.Println(err)
    return
  }

  total := 0

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()

    // parse (should return game # and a list of sets)
    gameNum, setList := parseGame(line)
    // test list of sets against limit
    if isValidSetList(setList) {
      total += gameNum
    }
  }

  fmt.Printf("Total: %d\n", total)
  return
}
