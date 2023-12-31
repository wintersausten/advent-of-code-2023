package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Workflow struct {
  tests [] func(p Part) string
  catch string
}

var workflows = map[string]Workflow{}

type Part struct {
  x int
  m int
  a int
  s int
}

func (p1 *Part) add(p2 Part) {
  p1.x += p2.x
  p1.m += p2.m
  p1.a += p2.a
  p1.s += p2.s
}

func (p Part) getProp(prop byte) int {
  switch prop {
  case byte('x'):
    return p.x
  case byte('m'):
    return p.m
  case byte('a'):
    return p.a
  case byte('s'):
    return p.s
  default:
    panic("Invalid prop")
  }
}

func (p *Part) setProp(prop byte, val int) {
  switch prop {
  case byte('x'):
    p.x = val
  case byte('m'):
    p.m = val
  case byte('a'):
    p.a = val
  case byte('s'):
    p.s = val
  default:
    panic("Invalid prop")
  }
}

func (p Part) accept() bool {
  wID := "in"
  for wID != "A" && wID != "R" {
    workflow := workflows[wID]
    wID = workflow.catch
    for _, test := range workflow.tests {
      if result := test(p); result != "" {
        wID = result
        break
      }
    }
  }
  return wID == "A"
}

func parseTest(testString string) func(Part) string {
  partProp := testString[0]
  operator := testString[1]

  colonIndex := strings.Index(testString, ":")
  resultId := testString[colonIndex+1:]

  val, _ := strconv.Atoi(testString[2:colonIndex])
  if operator == '>' {
    return func(p Part) string {
      if p.getProp(partProp) > val {
        return resultId
      }
      return ""
    }
  } else {
    return func(p Part) string {
      if p.getProp(partProp) < val {
        return resultId
      }
      return ""
    }
  }
}

func parseFile() []Part {
  file, _ := os.Open("./inputs/day19.txt")
  scanner := bufio.NewScanner(file)

  // parts workflows
  for scanner.Scan() && scanner.Text() != "" {
    line := scanner.Text()

    workflow := Workflow{}
    endIdIndex := strings.Index(line, "{")
    endTestsIndex := strings.Index(line, "}")

    wID := line[:endIdIndex]
    testStringSlice := strings.Split(line[endIdIndex+1:endTestsIndex], ",")

    workflow.tests = []func(p Part) string{}
    for i, testString := range testStringSlice {
      if i == len(testStringSlice) - 1 {
        workflow.catch = testString
      } else {
        workflow.tests = append(workflow.tests, parseTest(testString))
      }
    }

    workflows[wID] = workflow
  }

  // parse parts
  parts := []Part{}
  for scanner.Scan() && scanner.Text() != "" {
    line := scanner.Text()
    propVals := strings.Split(line[1:len(line)-1], ",")
    part := Part{}
    for _, propVal := range propVals {
      val, _ := strconv.Atoi(propVal[2:])
      part.setProp(propVal[0], val)
    }
    parts = append(parts, part)
  }
  return parts
}

func main() {
  parts := parseFile() 

  sumPart := Part{}
  for _, p := range parts {
    if p.accept() {
      sumPart.add(p)
    }
  }

  total := sumPart.x + sumPart.m + sumPart.a + sumPart.s
  fmt.Printf("Total: %d\n", total)

}
