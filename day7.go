package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Hand struct {
  cards string
  bid int
  strength int // 0 for high card - 6 for five of a kind
}

func cardValue(c rune) int {
  if unicode.IsDigit(c) {
    return int(c - '0')
  } else if c == 'T' {
    return 10
  } else if c == 'J' {
    return 1
  // } else if c == 'J' {
  //   return 11
  } else if c == 'Q' {
    return 12
  } else if c == 'K' {
    return 13
  } else if c == 'A' {
    return 14
  }
  panic("Card value didn't resolve")

}

func secondOrderLess(a Hand, b Hand) bool {
  for i := 0; i < 5; i++ {
    ax := cardValue(rune(a.cards[i]))
    bx := cardValue(rune(b.cards[i]))
    if ax < bx {
      return true
    } else if ax > bx {
      return false
    }
  }
  panic("Second order comparison didn't resolve")
}

func less(a Hand, b Hand) bool {
  if (a.strength < b.strength) {
    return true
  } else if (a.strength > b.strength) {
    return false
  } else {
    return secondOrderLess(a, b)
  }
}

type HandHeap []Hand

func (h HandHeap) Len() int { return len(h) }
func (h HandHeap) Less(i, j int) bool { return less(h[i], h[j]) }
func (h HandHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *HandHeap) Push(x interface{}) {
  *h = append(*h , x.(Hand))
}
func (h *HandHeap) Pop() interface{} {
  old := *h
  x := old[0]
  *h = old[1:]
  return x
}

// Part 1
// func getHandStrength(cards string) int {
//   cardMap := map[rune]int{}
//   max := 0
//   for _, c := range cards {
//     if _, exists := cardMap[c]; exists {
//       cardMap[c] += 1
//     } else {
//       cardMap[c] = 1
//     }
//     if cardMap[c] > max {
//       max = cardMap[c]
//     }
//   }
//   
//   switch max {
//   case 1:
//     return 0
//   case 2:
//     if len(cardMap) == 4 { // One pair
//       return 1
//     } else {               // Two pair
//       return 2
//     }
//   case 3:
//     if len(cardMap) == 3 {
//       return 3
//     } else {
//       return 4
//     }
//   case 4:
//     return 5
//   case 5:
//     return 6
//   default:
//     panic("How'd you do that??")
//   }
// }

func getHandStrength(cards string) int {
  cardMap := map[rune]int{}
  max, maxCard := 0, '0'
  wildCards := 0
  for _, c := range cards {
    if c == 'J' {
      wildCards++
      continue
    }

    if _, exists := cardMap[c]; exists {
      cardMap[c] += 1
    } else {
      cardMap[c] = 1
    }

    if cardMap[c] > max {
      max = cardMap[c]
      maxCard = c
    }
  }

  cardMap[maxCard] += wildCards
  
  switch cardMap[maxCard] {
  case 1:
    return 0
  case 2:
    if len(cardMap) == 4 { // One pair
      return 1
    } else {               // Two pair
      return 2
    }
  case 3:
    if len(cardMap) == 3 { // Three of a king
      return 3
    } else {               // Full house
      return 4
    }
  case 4:
    return 5
  case 5:
    return 6
  default:
    panic("How'd you do that??")
  }
}

func parseHandSlice(scanner *bufio.Scanner) []Hand {
  handList := []Hand{}
  for scanner.Scan() {
    line := scanner.Text()
    fields := strings.Fields(line)
    cards := fields[0]
    bid, _ := strconv.Atoi(fields[1])
    handList = append(handList, Hand{cards: cards, bid: bid, strength: getHandStrength(cards)})
  }
  return handList
}

func main() {
  file, err := os.Open("inputs/day7.txt")
  if err != nil {
    fmt.Println(err)
    return
  }

  scanner := bufio.NewScanner(file)

  handSlice := parseHandSlice(scanner)

  handHeap := HandHeap(handSlice)
  heap.Init(&handHeap)

  total := 0
  handHeapLen := handHeap.Len()
  for i := 1; i <= handHeapLen; i++ {
    bid := handHeap.Pop().(Hand).bid
    total += i * bid
    heap.Init(&handHeap)
  }

  fmt.Printf("Total: %d\n", total)
}
