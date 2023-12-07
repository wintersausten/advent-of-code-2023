package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type StartMarker struct {
  startRange int
  offset int
}

func parseSeed(scanner *bufio.Scanner) []int {
  scanner.Scan()
  line := scanner.Text()
  colonIndex := strings.Index(line, ":")
  seedStringSlice := strings.Fields(line[colonIndex+1:])
  seedSlice := make([]int, 0, len(seedStringSlice))
  for _, seedString := range seedStringSlice {
    seed, _ := strconv.Atoi(seedString)
    seedSlice = append(seedSlice, seed)
  }

  return seedSlice
}

func parseParameterMaps(scanner *bufio.Scanner) [][]StartMarker {
  allParameterMaps := [][]StartMarker{}
  for {
    startMarkerMap := map[int]StartMarker{0: {startRange: 0, offset: 0}}
    stillScanning := false
    for scanner.Scan() && strings.Index(scanner.Text(), ":") == -1 {}
    for scanner.Scan() && scanner.Text() != "" {
      stillScanning = true
      fields := strings.Fields(scanner.Text())
      destStart, _ := strconv.Atoi(fields[0])
      sourceStart, _ := strconv.Atoi(fields[1])
      mapLength, _ := strconv.Atoi(fields[2])
      
      startMarker := StartMarker{startRange: sourceStart, offset: destStart - sourceStart}
      startMarkerMap[sourceStart] = startMarker

      nextMapStart := sourceStart + mapLength
      // if the next mapping isn't already set, then it should be a direct 0 offset mapping
      if _, exists := startMarkerMap[nextMapStart]; !exists {
        startMarkerMap[nextMapStart] = StartMarker{startRange: nextMapStart, offset: 0}
      }
    }
    if !stillScanning {
      break
    }
    singleParameterMap := make([]StartMarker, 0, len(startMarkerMap))
    for _, value := range startMarkerMap {
      singleParameterMap = append(singleParameterMap, value)
    }
    sort.Slice(singleParameterMap, func(i, j int) bool {
      return singleParameterMap[i].startRange < singleParameterMap[j].startRange
    })
    allParameterMaps = append(allParameterMaps, singleParameterMap)
  }
  return allParameterMaps
}

func getLocation(seed int, parameterMaps [][]StartMarker) int {
  source := seed
  for _, parameterMap := range parameterMaps {
    // find corresponding StartMarker
    index := sort.Search(len(parameterMap), func(i int) bool {
      return parameterMap[i].startRange >= source
    })

    var offset int
    if index < len(parameterMap) {
        if parameterMap[index].startRange == source {
            offset = parameterMap[index].offset
        } else if index > 0 {
            offset = parameterMap[index-1].offset
        } else {
            offset = parameterMap[0].offset
        }
    } else {
        offset = parameterMap[len(parameterMap)-1].offset
    }
    source = source + offset
  }
  return source
}

func day5() {
  file, err := os.Open("inputs/day5.txt")
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println("Start")

  scanner := bufio.NewScanner(file)

  seedSlice := parseSeed(scanner)
  parameterMaps := parseParameterMaps(scanner)

  min := int(math.MaxInt64)
  // part 2: took 3 min
  start := time.Now()
  count := 0
  for i := 0; i < len(seedSlice); i += 2 {
    rangeStart, rangeLen := seedSlice[i], seedSlice[i+1]
    for seed := rangeStart; seed < (rangeStart + rangeLen); seed++ {
      if (seed - rangeStart) % 1000000 == 0 {
        count += 1000000
        fmt.Printf("Done with %d\n in %v\n", count, time.Since(start))
      }
      location := getLocation(seed, parameterMaps) 
      if location < min {
        min = location
      }
    }
  }
  // part1
  // for _, seed := range seedSlice {
  //   location := getLocation(seed, parameterMaps) 
  //   if location < min {
  //     min = location
  //   }
  // }

  fmt.Printf("Lowest Location: %d\n", min)
}
