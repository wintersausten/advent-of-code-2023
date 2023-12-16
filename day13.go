package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pattern [][]rune

func compareSlices(slice1, slice2 []rune) bool {
  if len(slice1) != len(slice2) {
    return false
  }
  for i, r1 := range slice1 {
    if r1 != slice2[i] {
      return false
    }
  }
  return true
}

func numDifferences(slice1, slice2 []rune) int {
  n := 0
  for i, r1 := range slice1 {
    if r1 != slice2[i] {
      n += 1
    }
  }
  return n
}


// func (p Pattern) horizontalPattern() int {
//   val := 0
//   for i := 1; i < len(p[0]); i++ {
//     for j := 0; (i - j - 1 >= 0) && (i + j < len(p[0])); j++ {
//       if !compareSlices(p[i+j], p[i-j-1]) {
//         break
//       }
//     }
//     // didn't break when comparing reflection, reflection found
//     val = i
//     break
//   }
//   return val
// }

func (p Pattern) rotateLeft() [][]rune {
  rows, cols := len(p), len(p[0])
  rotated := make([][]rune, cols)
  for i := 0; i < cols; i++ {
    rotated[i] = make([]rune, rows)
    for j := 0; j < rows; j++ {
      rotated[i][rows-1-j] = p[j][i]
    }
  }
  return rotated
}

// TODO: ignore exactly 1 flaw in pattern -- there MUST be 1 flaw
func (p Pattern) findPattern() int {
  val := 0
  RowLoop:
  for i := 1; i < len(p); i++ {
    smudges := 0
    for j := 0; (i - j - 1 >= 0) && (i + j < len(p)); j++ {
      smudges += numDifferences(p[i+j], p[i-j-1])
      if smudges > 1{
        continue RowLoop
      }
    }
    // didn't break when comparing reflection, reflection found
    if smudges == 1 {
      val = i
      break
    }
  }

  return val
}

func (p Pattern) analyze() int{
  if val := p.findPattern(); val != 0 {
    return val * 100
  } 

  p = p.rotateLeft()

  if val := p.findPattern(); val != 0 {
    return val
  } 

  panic("No pattern found")
}

func parsePatterns() []Pattern{
  file, _ := os.Open("./inputs/day13.txt")
  scanner := bufio.NewScanner(file)

  patterns := []Pattern{}
  i := 0
  for scanner.Scan() {
    patterns = append(patterns, Pattern{}) 
    line := scanner.Text()
    for line != "" {
      patterns[i] = append(patterns[i], []rune(line))

      scanner.Scan()
      line = scanner.Text()
    }
    i++

  }
  return patterns
}

func main() {
  patterns := parsePatterns()

  total := 0
  for _, pattern := range patterns {
    total += pattern.analyze()
  }

  fmt.Printf("Total: %d\n", total)
}
