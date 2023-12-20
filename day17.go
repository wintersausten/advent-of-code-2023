package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

type Direction struct {
  X int
  Y int
}

var (
  Up = Direction {X: 0, Y: -1}
  Down = Direction {X: 0, Y: 1}
  Left = Direction {X: -1, Y: 0}
  Right = Direction {X: 1, Y: 0}
)

var Directions = []Direction{Up, Down, Left, Right}

func (a Direction) opposite(b Direction) bool {
  switch a {
  case Up:
    return b == Down
  case Down:
    return b == Up
  case Left:
    return b == Right
  case Right:
    return b == Left
  default:
    return false // only for start
  }
}

type Point struct {
  X int
  Y int
}

type PointDiscovery struct {
  Point
  Direction
  mag int // magnitude of direction discovered in
}

// Node represents a point in the A* search.
type Node struct {
  PointDiscovery
  fScore int
}

// PriorityQueue implements a priority queue for Nodes based on their fScore.
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { return pq[i].fScore < pq[j].fScore }

func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
  n := x.(*Node)
  *pq = append(*pq, n)
  heap.Fix(pq, len(*pq)-1)
}

func (pq *PriorityQueue) Pop() interface{} {
  old := *pq
  n := old[0] // Get the element with the highest priority
  last := old[len(old)-1]
  old[0] = last
  *pq = old[:len(old)-1]
  heap.Fix(pq, 0) // Fix the heap starting from the root
  return n
}

func getMinHeatLoss(grid [][]int, start PointDiscovery, goal Point) int {
  // Initialize open and closed sets
  openSet := PriorityQueue{}
  openSet.Push(&Node{start, 0})

  // Cost from start along best known path.
  gScore := map[PointDiscovery]int{}
  gScore[start] = 0

  // Estimated total cost from start to goal through y.
  fScore := map[PointDiscovery]int{}
  fScore[start] = heuristic(start.Point, goal)

  for openSet.Len() > 0 {
    // Find the node with the lowest fScore
    current := openSet.Pop().(*Node)

    // part1: Return the g score if goal reached
    // if current.PointDiscovery.Point == goal {
    //   return gScore[current.PointDiscovery]
    // }
    // part2: Return the g score if goal reached && have gone at least 4 so you can sotp
    if current.PointDiscovery.Point == goal && current.PointDiscovery.mag >= 4 {
      return gScore[current.PointDiscovery]
    }

    // Check every valid neighbor
    for _, neighbor := range getValidNeighbors(current.PointDiscovery, len(grid[0]), len(grid)) {
      // Calculate g score
      tentativeGScore := gScore[current.PointDiscovery] + int(grid[neighbor.Point.Y][neighbor.Point.X])

      // If g score is better
      if existingGScore, exists := gScore[neighbor]; !exists || tentativeGScore < existingGScore {
        // Set g & f score
        gScore[neighbor] = tentativeGScore
        fScore[neighbor] = gScore[neighbor] + heuristic(neighbor.Point, goal)

        // Add to open set if not there already, if it's already there then fix open set order
        if !contains(openSet, neighbor) {
          openSet.Push(&Node{neighbor, fScore[neighbor]})
        } else {
          heap.Fix(&openSet, findIndex(openSet, neighbor))
        }
      }
    }
  }

  // Open set is empty but goal was never reached
  fmt.Println("No path found!")
  return -1
}

func heuristic(n, goal Point) int {
  // Manhattan distance
  return int(math.Abs(float64(n.X-goal.X)) + math.Abs(float64(n.Y-goal.Y)))
}

func getValidNeighbors(p PointDiscovery, width, height int) []PointDiscovery {
  neighbors := []PointDiscovery{}
  for _, dir := range Directions {
    // continue if going in the direction came from
    if p.Direction.opposite(dir){
      continue
    }
    point := Point{p.Point.X + dir.X, p.Point.Y + dir.Y}
    if point.X < 0 || point.X >= width || point.Y < 0 || point.Y >= height {
      continue
    }
    if dir == p.Direction {
      // part1: if would go more than 3 in one direction, can't go this direction
      // if p.mag == 3 {
      //   continue
      // }
      // part2: if would go more than 10 in one direction, can't go this direction
      if p.mag == 10 {
        continue
      }
      neighbors = append(neighbors, PointDiscovery{Point: point, Direction: dir, mag: p.mag+1})
    } else {
      // part2: if haven't gone straight at least 4 times, can't go this direction
      if p.mag < 4 {
        continue
      }
      neighbors = append(neighbors, PointDiscovery{Point: point, Direction: dir, mag: 1})
    }
  }
  return neighbors
}

func contains(pq PriorityQueue, p PointDiscovery) bool {
  for _, node := range pq {
    if node.PointDiscovery == p {
      return true
    }
  }
  return false
}

func findIndex(pq PriorityQueue, p PointDiscovery) int {
  for i, node := range pq {
    if node.PointDiscovery == p {
      return i
    }
  }
  fmt.Println("Point not in pq, you're using this wrong!")
  return -1
}

func parseCityMap() [][]int{
  file, _ := os.Open("./inputs/day17.txt")
  scanner := bufio.NewScanner(file)
  cityMap := [][]int{}
  for scanner.Scan() {
    line := scanner.Text()
    mapRow := []int{}
    for _, c := range line {
      mapRow = append(mapRow, int(c - '0'))
    }
    cityMap = append(cityMap, mapRow)
  }
  return cityMap
}

func main () {
  cityMap := parseCityMap()
  // the direction its fed matters lol
  minHeatLoss := getMinHeatLoss(cityMap, PointDiscovery{Point: Point{0, 0}, Direction: Down}, Point{X: len(cityMap[0])-1, Y: len(cityMap)-1})

  fmt.Printf("Min Heat Loss: %d\n", minHeatLoss)
}
