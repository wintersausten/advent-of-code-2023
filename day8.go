package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type LoopInfo struct {
  length int
  ends []int
}

type Node struct {
  val string
  left *Node
  right *Node
}

type Directions struct {
  directions string
  i int
}

func NewDirections(d string) *Directions {
  return &Directions{directions: d, i: -1}
}

func (d *Directions) next() rune {
  d.i += 1
  if d.i % len(d.directions) == 0 {
    d.i = 0
  }
  return rune(d.directions[d.i])
}

func parseDirections(scanner *bufio.Scanner) *Directions {
  scanner.Scan()
  line := scanner.Text()
  scanner.Scan() // moves ver the empty line under the directions
  return NewDirections(line)
}

func parseNodes(scanner *bufio.Scanner) []*Node {
  var starts []*Node
  nodeMap := map[string]*Node{}
  for scanner.Scan() {
    line := scanner.Text()
    // fmt.Sscanf(line, "%s = (%s, %s)", &i, &l, &r)
    i := line[:strings.Index(line, " ")]
    l := line[strings.Index(line, "(")+1: strings.Index(line, ",")]
    r := line[strings.Index(line, ",")+2:strings.Index(line, ")")]

    var newNode *Node
    if node, exists := nodeMap[i]; exists {
      newNode = node
    } else {
      newNode = &Node{val: i}
      nodeMap[i] = newNode
    }

    if node, exists := nodeMap[l]; exists {
      newNode.left = node
    } else {
      newNode.left = &Node{val: l}
      nodeMap[l] = newNode.left
    }
    if node, exists := nodeMap[r]; exists {
      newNode.right = node
    } else {
      newNode.right = &Node{val: r}
      nodeMap[r] = newNode.right
    }

    if i[2] == 'A' {
      starts = append(starts, newNode)
    }
  }
  return starts
}

func day8() {
  file, err := os.Open("inputs/day8.txt")
  if err != nil {
    fmt.Println(err)
  }
  
  scanner := bufio.NewScanner(file)

  // parseDirections & return some sort of circular ll implementation
  directions := parseDirections(scanner)  

  f := *directions
  fastDirections := &f

  // Parse nodes & return all starts (__A)
  nodeItr := parseNodes(scanner)

  steps := 0
  
  fastNodeItr := make([]*Node, len(nodeItr))
  copy(fastNodeItr, nodeItr)

  // confirm loops
  loops := make([]bool, len(nodeItr))
  for completeLoops := 0; completeLoops != len(nodeItr); {
    // fmt.Printf("Node: %s\n", nodeItr[1].val)
    
    d := directions.next()
    for i, node := range nodeItr {
      if d == 'R' {
        nodeItr[i] = node.right
      } else {
        nodeItr[i] = node.left
      }

      // if end && starts[i].val[2] != 'Z' {
      //   end = false
      // }
    }
    // First fast move
    d2 := fastDirections.next()
    for i, node := range fastNodeItr {
      if d2 == 'R' {
        fastNodeItr[i] = node.right
      } else {
        fastNodeItr[i] = node.left
      }
    }
    // Second fast move
    d2 = fastDirections.next()
    for i, node := range fastNodeItr {
      if d2 == 'R' {
        fastNodeItr[i] = node.right
      } else {
        fastNodeItr[i] = node.left
      }
    }
    // check for loop
    for i, node := range nodeItr {
      if node == fastNodeItr[i] && fastDirections.i == directions.i && !loops[i]{
        loops[i] = true
        completeLoops += 1
      }
    }
  }

  // find loop lengths & end options 
  loops = make([]bool, len(nodeItr))
  starts := make([]*Node, len(nodeItr))
  copy(starts, nodeItr)
  loopInfo := make([]LoopInfo, len(nodeItr))
  for completeLoops, startStep, directionStart := 0, steps, directions.i; completeLoops != len(nodeItr); { 
    // fmt.Printf("Node: %s\n", nodeItr[1].val)
    steps ++
    d := directions.next()

    for i, node := range nodeItr {
        if d == 'R' {
        nodeItr[i] = node.right
      } else {
        nodeItr[i] = node.left
      }

      if !loops[i] {
        if nodeItr[i].val[2] == 'Z' { 
          loopInfo[i].ends = append(loopInfo[i].ends, steps)
        }
        if nodeItr[i] == starts[i] && directionStart == directions.i {
          loops[i] = true
          completeLoops += 1
          loopInfo[i].length = steps - startStep
        }
      }
    }
  }

  // bring next closest to 0 until all are 0 at same
  nextEndInLoop := make([]int, len(loopInfo))
  for i, loopInfo := range loopInfo {
    min := int(math.MaxInt64)
    for _, end := range loopInfo.ends {
      toNextEnd := loopInfo.length - ((steps - end) % loopInfo.length)
      nextEnd := toNextEnd + steps
      if nextEnd < min {
        min = nextEnd
      }
    }
    nextEndInLoop[i] = min
  }

  for {
    min, minIndex, max := nextEndInLoop[0], 0, nextEndInLoop[0]
    for i, nextEnd := range nextEndInLoop {
      if nextEnd > max {
        max = nextEnd
      }
      if nextEnd < min {
        min = nextEnd
        minIndex = i
      }
    }
    if min == max {
      fmt.Printf("Steps: %d\n", min)
      break
    }
    // just assuming 1 end in each loop because I've confirmed in debugger & I'm tired
    nextEndInLoop[minIndex] = min + loopInfo[minIndex].length
  }


  // return steps
}
