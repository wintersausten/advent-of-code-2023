package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Coord struct {
  x int
  y int
}

func getMatchingCoords(file *os.File, matcher func(rune) bool) []Coord {
  scanner := bufio.NewScanner(file)

  mathingCoords := []Coord{}

  for y := 0; scanner.Scan(); y++ {
    for x, char := range scanner.Text() {
      if matcher(char) {
        mathingCoords = append(mathingCoords, Coord{x: x, y: y})
      }
    }
  }

  return mathingCoords
}

func symbolMatcher(char rune) bool {
  if !unicode.IsDigit(char) && char != '.' {
    return true
  }
  return false
}

func gearMatcher(char rune) bool {
  if char == '*' {
    return true
  }
  return false
}

// part 1 implementation of getAdjacentCoords
// func getAdjacentCoords(coordList []Coord) map[Coord]struct{} {
//   adjacentCoords := map[Coord]struct{}{}
//   for _, source := range coordList {
//     for dx := -1; dx <= 1; dx++ {
//       for dy := -1; dy <= 1; dy++ {
//         if dx == 0 && dy == 0 {
//           continue
//         }
//
//         // don't need to worry about out of bounds, as long as this set is a superset of the 'true' adjacent coords
//         adjacentCoords[Coord{x: source.x + dx, y: source.y + dy}] = struct{}{}
//       }
//     }
//   } 
//   return adjacentCoords
// }

// gets all adjacent coordinates relative to a source coord list, returns in the form of a map of Coord: list of symbol ids that the coord is adjacent to 
func getAdjacentCoords(coordList []Coord) map[Coord][]int {
  adjacentCoords := map[Coord][]int{}
  for sourceId, source := range coordList {
    for dx := -1; dx <= 1; dx++ {
      for dy := -1; dy <= 1; dy++ {
        if dx == 0 && dy == 0 {
          continue
        }

        // don't need to worry about out of bounds, as long as this set is a superset of the 'true' adjacent coords
        adjacent := Coord{x: source.x + dx, y: source.y + dy}
        if sourceIdList, exists := adjacentCoords[adjacent]; exists {
          adjacentCoords[adjacent] = append(sourceIdList, sourceId)
        } else {
          adjacentCoords[Coord{x: source.x + dx, y: source.y + dy}] = []int{sourceId}
        }
      }
    }
  } 
  return adjacentCoords
}

func getPartNumber(source Coord, line []rune) (int, []Coord) {
  coords := []Coord{}
  partNumberString := strings.Builder{}
  i := source.x
  for ; i > 0 && unicode.IsDigit(line[i-1]); i-- {}
  for ; i < len(line) && unicode.IsDigit(line[i]); i++ {
    coords = append(coords, Coord{x: i, y: source.y}) 
    partNumberString.WriteRune(line[i])
  }
 
  partNumber, _ := strconv.Atoi(partNumberString.String())
  return partNumber, coords
}

func totalPartNumbers(file *os.File) int {
  symbolCoords := getMatchingCoords(file, symbolMatcher)
  symbolAdjacentCoords := getAdjacentCoords(symbolCoords)

  file.Seek(0, io.SeekStart)
  scanner := bufio.NewScanner(file)
  invalid := map[Coord]struct{}{}
  total := 0
  for y := 0; scanner.Scan(); y++ {
    for x, char := range scanner.Text() {
      // if the char is not a digit, we can ignore it
      if !unicode.IsDigit(char) {
        continue
      }

      checkCoord := Coord{x: x, y: y}

      // ignore char is the coordinate has been invalidated (number used already)
      if _, isInvalid := invalid[checkCoord]; isInvalid {
        continue
      }

      // if the character is a digit, is still valid, AND is symbol adjacent
      if _, isSymbolAdjacent := symbolAdjacentCoords[checkCoord]; isSymbolAdjacent {
        // parse out the part number
        partNumber, partNumberCoords := getPartNumber(checkCoord, []rune(scanner.Text()))

        total += partNumber
        // invalidate used part number coords
        for _, coord := range partNumberCoords {
          invalid[coord] = struct{}{}
        }
      }
    } 
  }
  return total
}

func totalGearRatio(file *os.File) int {
  gearCoords := getMatchingCoords(file, gearMatcher)
  gearAdjacentCoords := getAdjacentCoords(gearCoords)

  file.Seek(0, io.SeekStart)
  scanner := bufio.NewScanner(file)
  gearPartMap := map[int][]int{}
  total := 0
  for y := 0; scanner.Scan(); y++ {
    for x, char := range scanner.Text() {
      // if the char is not a digit, we can ignore it
      if !unicode.IsDigit(char) {
        continue
      }

      checkCoord := Coord{x: x, y: y}

      // if the character is a digit, is still valid, AND is gear adjacent
      if gearIdList, isGearAdjacent := gearAdjacentCoords[checkCoord]; isGearAdjacent {
        // parse out the part number
        partNumber, partNumberCoords := getPartNumber(checkCoord, []rune(scanner.Text()))

        // add the part number to every gear it is adjacent to
        for _, gearId := range gearIdList {
          if partList, exists := gearPartMap[gearId]; exists {
            gearPartMap[gearId] = append(partList, partNumber)
          } else {
            gearPartMap[gearId] = []int{partNumber}
          }
          
          // remove gear from adjacent list for other cords that are part of the part number
          for _, partCoord := range partNumberCoords {
            // this pretty messy but I'm sick of this puzzle so...
            gearList := gearAdjacentCoords[partCoord]
            for i, searchGearId := range gearList {
              if searchGearId == gearId {
                gearAdjacentCoords[partCoord] = append(gearList[:i], gearList[i+1:]...)
              } 
            }
          }
        }
      }
    } 
  }

  // add up all the gear ratios where a gear must have 2 parts adjacent
  for _, partList := range gearPartMap {
    if len(partList) == 2 {
      total += (partList[0] * partList[1])
    }
  }

  return total
}

func main () {
  file, err := os.Open("inputs/day3.txt")
  if err != nil {
    fmt.Println(err)
    return
  }

  // total := totalPartNumbers(file)
  total := totalGearRatio(file)

  fmt.Printf("Total: %d\n", total)
}
