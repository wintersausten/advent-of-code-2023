package main

import (
	"advent-of-code-2023/utils"
	"bufio"
	"fmt"
	"os"
	"unicode"
)

var numMap map[string]int = map[string]int{
  "one": 1,
  "two": 2,
  "three": 3,
  "four": 4,
  "five": 5,
  "six": 6,
  "seven": 7,
  "eight": 8,
  "nine": 9,
}

var numSlice = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
var reverseNumSlice = []string{"eno", "owt", "eerht", "ruof", "evif", "xis", "neves", "thgie", "enin"}

// func main() {
func day1() {
  file, err := os.Open("inputs/day1.txt")
  if err != nil {
    fmt.Println(err)
    return
  }

  forwardTrie := utils.BuildTrie(numSlice)
  backwardTrie := utils.BuildTrie(reverseNumSlice)
  scanner := bufio.NewScanner(file)
  total := 0
  for scanner.Scan() {
    line := scanner.Text()
    runes := []rune(line)
    forwardTrie.Itr = forwardTrie.Root
    backwardTrie.Itr = backwardTrie.Root
    value := 0
    
    // find first number and add to value
    forwardTrie := utils.BuildTrie(numSlice)
    FowardLoop:
    for i, char := range runes {
      if unicode.IsDigit(char) {
        value += 10 * int(char - '0')
        break
      } else if forwardTrie.ItrStep(char) {
        for j := i + 1; j < len(runes) && forwardTrie.ItrStep(runes[j]); j++ {
          if forwardTrie.Itr.End {
            value += 10 * numMap[forwardTrie.GetPath()]
            break FowardLoop
          } 
        }
      }
    }

    // find last number and add to value
    backwardTrie = utils.BuildTrie(reverseNumSlice)
    BackwardLoop:
    for i := len(runes) - 1; i >= 0; i-- {
      char := runes[i]
      if unicode.IsDigit(char) {
        value += int(char - '0')
        break
      } else if backwardTrie.ItrStep(char) {
        for j := i - 1; j >= 0 && backwardTrie.ItrStep(runes[j]); j-- {
          if backwardTrie.Itr.End {
            value += numMap[utils.ReverseString(backwardTrie.GetPath())]
            break BackwardLoop 
          } 
        }
      }
    }

    total += value
  }

  fmt.Printf("Total: %d\n", total)
  return
}

