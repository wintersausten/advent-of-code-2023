package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Record struct {
  springs []rune 
  groups []int
}

func countDamagedGroups(runes []rune) []int {
    var groupSizes []int
    count := 0

    for _, r := range runes {
        if r == '#' {
            count++
        } else if count > 0 {
            groupSizes = append(groupSizes, count)
            count = 0
        }
    }

    if count > 0 {
        groupSizes = append(groupSizes, count)
    }

    return groupSizes
}

func parseSpringRecords() []Record {
  file, _ := os.Open("./inputs/day12.txt")
  scanner := bufio.NewScanner(file)

  records := []Record{}
  for scanner.Scan() {
    line := scanner.Text()
    spaceIndex := strings.Index(line, " ")

    springs := []rune(line[:spaceIndex])

    groupsStrings := strings.Split(line[spaceIndex+1:], ",")
    groups := []int{}
    for _, s := range groupsStrings {
      i, _ := strconv.Atoi(s)
      groups = append(groups, i)
    }

    originalGroups := make([]int, len(groups))
    copy(originalGroups, groups)
    for i := 0; i < 4; i++ {
        groups = append(groups, originalGroups...)
    }

    originalSprings := make([]rune, len(springs))
    copy(originalSprings, springs)
    for i:= 0; i < 4; i++ {
      springs = append(springs, '?') 
      springs = append(springs, originalSprings...)
    }

    records = append(records, Record{springs, groups})
  }
  
  return records
}

// TODO: doesn't need to include the passed parts of record
func generateMemoKey (springsRem []rune, groups []int, currentDamage int) string {
 var sb strings.Builder
  for _, r := range springsRem {
    sb.WriteRune(r)
  }
  sb.WriteString("-")
  for _, i := range groups {
    sb.WriteString(fmt.Sprintf("%d-", i))
  }
  sb.WriteString("|")
  sb.WriteString(fmt.Sprintf("%d-", currentDamage))
  return sb.String()
}

func generateArrangements(record Record, i int, currentDamaged int, targetDamaged int, memo *map[string]int) int {
  memoKey := generateMemoKey(record.springs[i:], record.groups, currentDamaged)
  if val, exists := (*memo)[memoKey]; exists {
    return val
  }

  // move i to next ? or end
  for i < len(record.springs) && record.springs[i] != '?'  {
    if record.springs[i] == '#' {
      // start/continue damage group
      currentDamaged++
      // if exceed target group size, stop 
      if currentDamaged > targetDamaged {
        (*memo)[memoKey] = 0
        return 0
      }
    } else {
      if currentDamaged != 0 {
        // damage group ended, make sure size is right (if not, break)
        if currentDamaged != targetDamaged {
          (*memo)[memoKey] = 0
          return 0
        }

        record.groups = record.groups[1:]
        currentDamaged = 0
        if len(record.groups) > 0 {
          targetDamaged = record.groups[0]
        } else {
          targetDamaged = 0
        }
      }
    }
    i++ 
  }

  // completed record
  if i == len(record.springs) {
    // finished last group or no more & not in a damaged group
    if (currentDamaged == targetDamaged) && len(record.groups) <= 1 {
      (*memo)[memoKey] = 1
      return 1
    }
    (*memo)[memoKey] = 0
    return 0
  }

  // generate arrangements with ? changed to # and . 
  record.springs[i] = '#'
  possibleArrangements := 0
  possibleArrangements += generateArrangements(record, i, currentDamaged, targetDamaged, memo)
  record.springs[i] = '.'
  possibleArrangements += generateArrangements(record, i, currentDamaged, targetDamaged, memo)

  record.springs[i] = '?'
  (*memo)[memoKey] = possibleArrangements
  return possibleArrangements
}

func getTotalDamaged(groups []int) int {
  total := 0
  for _, size := range groups {
    total += size
  }
  return total
}

func getArrangementCount(record Record, memo *map[string]int) int {
  return generateArrangements(record, 0, 0, record.groups[0], memo)
}

func main() {
  records := parseSpringRecords()

  total := 0
  memo := make(map[string]int)
  for i, record := range records {
    total += getArrangementCount(record, &memo)
    if i % 100 == 0 {
      fmt.Printf("i: %d  ----   total: %d\n", i, total)
    }
  }

  fmt.Printf("Total: %d\n", total)
}
