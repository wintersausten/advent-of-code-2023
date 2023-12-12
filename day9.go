package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

// could pontentially cash levels if calucating many repeats
func predictNextValue(set []int) int {
  level := make([]int, len(set))
  copy(level, set)
  lastValues := []int{level[len(level)-1]}

  // refactor, could end on a 0 last value without being done
  static := false 
  for !static {
    newLevel := make([]int, len(level)-1)

    static = true
    prev := level[0]
    for i := 1; i < len(level); i++ {
      newVal := level[i] - prev
      if newVal != 0 {
        static = false
      }
      newLevel[i-1] = newVal
      prev = level[i]

      if i == len(level) - 1 {
        lastValues = append(lastValues, newVal)
      }
    }
    level = newLevel
  } 

  next := 0
  for i := len(lastValues) - 2; i >= 0; i-- {
    next = lastValues[i] + next
  }
  
 return next
}

func main() {
  file, _ := os.Open("./inputs/day9.txt")

  content, _ := io.ReadAll(file)


  lines := strings.Split(string(content), "\n")
  lines = lines[:len(lines)-1]
  sets := make([][]int, len(lines))
  for i, line := range lines {
    set := []int{}
    fields := strings.Fields(line)
    for _, field := range fields {
      val, _ := strconv.Atoi(field)
      set = append(set, val)
    }
    sets[i] = set
  }

    
  result := 0
  for _, set := range sets {
    slices.Reverse(set)
    result += predictNextValue(set)
  }


  fmt.Printf("Total: %d\n", result)
}
