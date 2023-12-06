package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func part1(file *os.File) int {
  scanner := bufio.NewScanner(file)
  total := 0
  for scanner.Scan() {
    line := scanner.Text()
    colonIndex := strings.Index(line, ":")
    barIndex := strings.Index(line, "|")

    cardSegment := strings.TrimSpace(line[colonIndex+1:barIndex])
    winningSegment := strings.TrimSpace(line[barIndex+1:])

    card := map[string]struct{}{}
    for _, num := range strings.Fields(cardSegment) {
      card[num] = struct{}{}
    }

    cardValue := 0
    for _, winningNum := range strings.Fields(winningSegment) {
      if _, exists := card[winningNum]; exists {
        if cardValue == 0 {
          cardValue = 1
        } else {
          cardValue *= 2
        }
      }
    }
    total += cardValue
  }
  return total
}

func part2(file *os.File) int {
  scanner := bufio.NewScanner(file)
  cardCopyMap := map[int]int{}
  cardNum := 1
  for ; scanner.Scan(); cardNum++{
    // initialize card copy if none is present
    if _, exists := cardCopyMap[cardNum]; !exists {
      cardCopyMap[cardNum] = 1
    }

    line := scanner.Text()

    colonIndex := strings.Index(line, ":")
    barIndex := strings.Index(line, "|")
    cardSegment := strings.TrimSpace(line[colonIndex+1:barIndex])
    winningSegment := strings.TrimSpace(line[barIndex+1:])

    // parse card numbers into map
    card := map[string]struct{}{}
    for _, num := range strings.Fields(cardSegment) {
      card[num] = struct{}{}
    }

    // identify matches from winning numbers
    matches := 0
    for _, winningNum := range strings.Fields(winningSegment) {
      if _, exists := card[winningNum]; exists {
        matches += 1
      }
    }

    // get a list of card numbers to copy, multiply by number of copies for current card
    for i := 1; i <= matches; i++ {
      offsetCardNum := cardNum + i
      if count, exists := cardCopyMap[offsetCardNum]; exists {
        cardCopyMap[offsetCardNum] = count + cardCopyMap[cardNum]
      } else {
        cardCopyMap[offsetCardNum] = 1 + cardCopyMap[cardNum]
      }
    }
  }
  // todo: prevent card copies past table end

  total := 0
  for cardNumKey, count := range cardCopyMap {
    // only count cards within range
    if cardNumKey < cardNum {
      total += count
    }
  }
  return total
}

func originalMain() {
  file, err := os.Open("inputs/day4.txt")
  if err != nil {
    fmt.Println(err)
    return
  }

  // total := part1(file)
  total := part2(file)

  fmt.Printf("Total: %d\n", total)
}
