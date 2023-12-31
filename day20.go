package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PulseType int

type Module interface {
  Pulse(PulseType, Module) []ModulePulse
  setDestination(Module)
  getDestinations() []Module
  getId() string
}

const (
  low PulseType = 0
  high PulseType = 1
)

type BaseModule struct {
  id string
  Destinations []Module
}

type FlipFlopModule struct {
  BaseModule
  On bool
}

type ConjunctionModule struct {
  BaseModule
  SourcePulses map[Module]PulseType
}

type BroadcastModule struct {
  BaseModule
}

type UntypedModule struct {
  BaseModule
}

type ModulePulse struct {
  Module
  PulseType
  Source Module
}

func (m *FlipFlopModule) Pulse(pt PulseType, source Module) []ModulePulse { 
  if pt == high {
    return nil
  }

  if m.On {
    m.On = false
    pt = low 
  } else {
    m.On = true
    pt = high
  }

  destPulses := []ModulePulse{}
  for _, dest := range m.Destinations {
    destPulses = append(destPulses, ModulePulse{dest, pt, m})
  }
  return destPulses
}

func (m *FlipFlopModule) setDestination(source Module) {
  m.BaseModule.Destinations = append(m.BaseModule.Destinations, source)
}

func (m *FlipFlopModule) getDestinations() []Module {
  return m.BaseModule.Destinations
}

func (m *FlipFlopModule) getId() string {
  return m.BaseModule.id
}

func (m *ConjunctionModule) Pulse(pt PulseType, source Module) []ModulePulse {
  m.SourcePulses[source] = pt

  pt = low
  for _, sourcePt := range m.SourcePulses {
    if sourcePt == low {
      pt = high
      break
    }
  }

  destPulses := []ModulePulse{}
  for _, dest := range m.Destinations {
    destPulses = append(destPulses, ModulePulse{dest, pt, m})
  }
  return destPulses
}

func (m *ConjunctionModule) setDestination(source Module) {
  m.BaseModule.Destinations = append(m.BaseModule.Destinations, source)

  // if _, exists := source.SourcePulses[m]; !exists {
  //   source.SourcePulses[m] = low
  // }
}

func (m *ConjunctionModule) getDestinations() []Module {
  return m.BaseModule.Destinations
}

func (m *ConjunctionModule) getId() string {
  return m.BaseModule.id
}

func (m *UntypedModule) Pulse(pt PulseType, source Module) []ModulePulse { 
  return nil
}

func (m *UntypedModule) setDestination(source Module) {
  return
}

func (m *UntypedModule) getDestinations() []Module {
  return nil
}

func (m *UntypedModule) getId() string {
  return m.BaseModule.id
}

var pulseCount = map[PulseType]int{
  low: 0,
  high: 0,
}

func (m *BroadcastModule) Pulse(pt PulseType, source Module) []ModulePulse {
  // pulseQueue := []ModulePulse{ModulePulse{m, pt}}

  // queue up broadcaster pulse
  pulseCount[pt] += 1
  pulseQueue := []ModulePulse{}
  for _, dest := range m.getDestinations() {
    pulseCount[pt] += 1
    pulseQueue = append(pulseQueue, dest.Pulse(pt, m)...)
  }

  for len(pulseQueue) > 0 {
    pulse := pulseQueue[0]
    pulseCount[pulse.PulseType] += 1
    pulseQueue = pulseQueue[1:]
    pulseQueue = append(pulseQueue, pulse.Module.Pulse(pulse.PulseType, pulse.Source)...)
  }
  return nil
}

func (m *BroadcastModule) setDestination(source Module) {
  m.BaseModule.Destinations = append(m.BaseModule.Destinations, source)
}

func (m *BroadcastModule) getDestinations() []Module {
  return m.BaseModule.Destinations
}

func (m *BroadcastModule) getId() string {
  return m.BaseModule.id
}

var moduleMap = map[string]Module{}
var broadcaster *BroadcastModule
var rx *Module
var simpleDestinationMap = map[Module][]string{}

func parseModule(moduleLine string) {
  var module Module
  var moduleName string

  fields := strings.Fields(moduleLine)
  moduleString := fields[0]
  
  if moduleString == "broadcaster" {
    moduleName = moduleString
    module = &BroadcastModule{BaseModule{id: moduleName}}
    broadcaster = module.(*BroadcastModule)
  } else {
    moduleType := moduleString[0]
    moduleName = moduleString[1:]
    switch moduleType {
    case '&':
      module = &ConjunctionModule{BaseModule{id: moduleName}, map[Module]PulseType{}}
    case '%':
      module = &FlipFlopModule{BaseModule{id: moduleName}, false}
    default:
      panic("Invalid module type")
    }
  }
  moduleMap[moduleName] = module
  for _, dest := range fields[2:] {
    simpleDestinationMap[module] = append(simpleDestinationMap[module], strings.TrimRight(dest, ","))
  }
}

func connectModules() {
  for _, module := range moduleMap {
    for _, dest := range simpleDestinationMap[module] {
      if _, exists := moduleMap[dest]; !exists {
        moduleMap[dest] = &UntypedModule{BaseModule{id: dest}}
      }
      // if type of moduleMap[dest] == ConjunctionModule, need to add module as a source
      if c, ok := moduleMap[dest].(*ConjunctionModule); ok {
        c.SourcePulses[module] = low
      }
      module.setDestination(moduleMap[dest])
    }
  }
}

func parseCommsSystem() {
  file, _ := os.Open("./inputs/day20.txt")
  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    parseModule(scanner.Text())
  }

  connectModules()
}

func main() {
  parseCommsSystem()

  for i := 0; i < 1000; i++ {
  // for i := 0; i < 4; i++ {
    broadcaster.Pulse(low, broadcaster)
    // fmt.Println()
  }

  fmt.Printf("Result: %d\n", pulseCount[low]*pulseCount[high])
}
