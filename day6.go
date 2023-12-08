package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseTimeLimits(scanner *bufio.Scanner) []int {
  scanner.Scan()
  line := scanner.Text()
  colonIndex := strings.Index(line, ":")
  timeStrings := strings.Fields(line[colonIndex+1:])
  times := []int{}
  for _, t := range timeStrings {
    time, _ := strconv.Atoi(t)
    times = append(times, time)
  }
  return times
}

func parseToBeats(scanner *bufio.Scanner, times []int) map[int]int {
  scanner.Scan()
  line := scanner.Text()
  colonIndex := strings.Index(line, ":")
  distanceStrings := strings.Fields(line[colonIndex+1:])
  toBeat := map[int]int{}
  for i, d := range distanceStrings {
    distance, _ := strconv.Atoi(d)
    toBeat[times[i]] = distance
  }
  return toBeat
}

func parseTime(scanner *bufio.Scanner) int {
  scanner.Scan()
  line := scanner.Text()
  colonIndex := strings.Index(line, ":")
  timeStrings := strings.Fields(line[colonIndex+1:])
  timeString := strings.Join(timeStrings, "")
  time, _ := strconv.Atoi(timeString)
  return time
}

func parseDistance(scanner *bufio.Scanner) int {
  scanner.Scan()
  line := scanner.Text()
  colonIndex := strings.Index(line, ":")
  distanceStrings := strings.Fields(line[colonIndex+1:])
  distanceString := strings.Join(distanceStrings, "")
  distance, _ := strconv.Atoi(distanceString)
  return distance
}

func day6() {
  file, err := os.Open("inputs/day6.txt")
  if err != nil {
    fmt.Println(err)
    return
  }
  scanner := bufio.NewScanner(file)

  timeLimit := parseTime(scanner)
  distanceMin := parseDistance(scanner)
  possible := 0
  for timeHeld := 0; timeHeld <= timeLimit; timeHeld++ {
    distance := timeHeld * (timeLimit - timeHeld)
    if distance > distanceMin {
      possible += 1
    }
  }
  fmt.Printf("Possible: %d\n", possible)

  // timeLimits := parseTimeLimits(scanner)
  // toBeats := parseToBeats(scanner, timeLimits)
  // total := -1
  // for timeLimit, distanceMin := range toBeat {
  //   possible := 0
  //   for timeHeld := 0; timeHeld <= timeLimit; timeHeld++ {
  //     distance := timeHeld * (timeLimit - timeHeld)
  //     if distance > distanceMin {
  //       possible += 1
  //     }
  //   }
  //   if total == -1 {
  //     total = possible
  //   } else {
  //     total *= possible
  //   }
  // }
  //
  // fmt.Printf("Total: %d\n", total)
}
