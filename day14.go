package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func rotateRight(a [][]rune) [][]rune {
  rows, cols := len(a), len(a[0])
  rotated := make([][]rune, cols)
  for i := 0; i < cols; i++ {
    rotated[i] = make([]rune, rows)
    for j := 0; j < rows; j++ {
      rotated[i][j] = a[j][cols-1-i]
    }
  }
  return rotated
}

func rotateLeft(a [][]rune) [][]rune {
  rows, cols := len(a), len(a[0])
  rotated := make([][]rune, cols)
  for i := 0; i < cols; i++ {
    rotated[i] = make([]rune, rows)
    for j := 0; j < rows; j++ {
      rotated[i][rows-1-j] = a[j][i]
    }
  }
  return rotated
}

func parsePlatform() [][]rune {
  file, _ := os.Open("./inputs/day14.txt")
  scanner := bufio.NewScanner(file)

  platform := [][]rune{}
  for scanner.Scan() {
    line := scanner.Text()
    platform = append(platform, []rune(line))

  }

  return platform
}

func tiltWest(platform [][]rune) {
  for _, row := range platform {
    open := -1
    for i, val := range row {
      switch val {
      case 'O':
        if open != -1 {
          row[open] = 'O'
          row[i] = '.'
          open++
        }
      case '#':
        open = -1
      default:
        if open == -1 {
          open = i
        }
      } 
    }
  }
}

func tiltEast(platform [][]rune) {
  for _, row := range platform {
    open := -1
    // for i, val := range row {
    for i := len(row)-1; i >= 0; i-- {
      switch row[i] {
      case 'O':
        if open != -1 {
          row[open] = 'O'
          row[i] = '.'
          open--
        }
      case '#':
        open = -1
      default:
        if open == -1 {
          open = i
        }
      } 
    }
  }
}

func tiltNorth(platform [][]rune) {
  for i:=0; i<len(platform[0]); i++ {
    open := -1
    for j:=0; j<len(platform); j++ {
      switch platform[j][i] {
      case 'O':
        if open != -1 {
          platform[open][i] = 'O'
          platform[j][i] = '.'
          open++
        }
      case '#':
        open = -1
      default:
        if open == -1 {
          open = j
        }
      } 
    }
  }
}

func tiltSouth(platform [][]rune) {
  for i:=0; i<len(platform[0]); i++ {
    open := -1
    for j:=len(platform)-1; j>=0; j-- {
      switch platform[j][i] {
      case 'O':
        if open != -1 {
          platform[open][i] = 'O'
          platform[j][i] = '.'
          open--
        }
      case '#':
        open = -1
      default:
        if open == -1 {
          open = j
        }
      } 
    }
  }
}

func calculateNorthLoad(platform [][]rune) int {
  load := 0
  for i:=0; i<len(platform[0]); i++{
    for j:=len(platform)-1; j>=0; j-- {
      if platform[j][i] == 'O' {
        load += len(platform)-j
      }
    }
  }
  return load
}

func platformString(platform[][]rune) string {
  s := strings.Builder{}
  for _, row := range platform {
    s.WriteString(string(row))
    s.WriteString("\n")
  }
  return s.String()
}

// Part 2 Thoughts
// can't use rotate hack anything if repeating so many  times  
// generalize tilt function to go around each side
// can either change calculate load to work for north or still preform a single rotate to start

func main() {
  platform := parsePlatform()
  cycleTarget := 1000000000

  memo := map[string]int{}
  OuterLoop:
  for i:=0; i<cycleTarget; i++ {
    s := platformString(platform)
    if val, exists := memo[s]; exists {
      loopLen := i - val
      offset := (cycleTarget - val) % loopLen
      for j:=0; j<offset; j++ {
        tiltNorth(platform)
        tiltWest(platform)
        tiltSouth(platform)
        tiltEast(platform)
      }
      break OuterLoop
    } 
    memo[s] = i
    if i%1000000==0 {
      fmt.Printf("i: %d\n", i)
    }
    tiltNorth(platform)
    tiltWest(platform)
    tiltSouth(platform)
    tiltEast(platform)

  }
  load := calculateNorthLoad(platform)

  fmt.Printf("Load: %d\n", load)
}
