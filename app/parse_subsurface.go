package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ojarva/subsurface-statistics/counter"

	"github.com/ojarva/subsurface-statistics/subsurfacetypes"
)

const unknownDiveSite string = "unknown"

var filenameFlag = flag.String("filename", "filename.ssrf", "Filename to be parsed")
var sortByFlag = flag.String("sort", "count", "Field used for sorting")

type statsContainerMap map[statType]counter.LastCounterStats

func (scm *statsContainerMap) Add(statType statType, name string, timeSince *time.Duration) {
	_, exists := (*scm)[statType]
	if !exists {
		(*scm)[statType] = make(counter.LastCounterStats)
	}
	(*scm)[statType].Add(name, timeSince)
}

type statType int

//go:generate stringer -type=statType
const (
	DiveLength statType = iota
	Buddies
	Cylinders
	MeanDepth
	MaxDepth
	Temperature
	DiveSite
	TagStat
)

type diveSiteMap map[string]string

func (dsm diveSiteMap) FetchByID(id string) string {
	diveSiteName, found := dsm[id]
	if found {
		return diveSiteName
	}
	return unknownDiveSite
}

func diveReceiver(c chan subsurfacetypes.Dive, wg *sync.WaitGroup, diveSites *diveSiteMap) {
	defer wg.Done()
	statsContainer := make(statsContainerMap)
	for dive := range c {
		processDive(&dive, &statsContainer, diveSites)
	}
	for _, stats := range statsContainer {
		stats.PrintStats(*sortByFlag)
	}
}

func processDive(dive *subsurfacetypes.Dive, statsContainer *statsContainerMap, diveSites *diveSiteMap) {
	if dive.IsInvalid() {
		return
	}
	timeSinceDive := dive.TimeSince()
	buddies := dive.BuddyList()
	for _, buddy := range buddies {
		(*statsContainer).Add(Buddies, buddy, &timeSinceDive)
	}
	usedCylinders := map[string]bool{}
	for _, cylinder := range dive.Cylinders {
		// Deduplicate cylinders used in a single dive; subsurface occasionally creates duplicate cylinders.
		// This won't work well for multiple stages with the same size but it's good enough for most cases.
		_, ok := usedCylinders[cylinder.Size]
		if ok {
			continue
		}
		usedCylinders[cylinder.Size] = true
		(*statsContainer).Add(Cylinders, cylinder.Size, &timeSinceDive)
	}
	(*statsContainer).Add(DiveLength, subsurfacetypes.DurationToSlot(dive.Duration()), &timeSinceDive)
	(*statsContainer).Add(MeanDepth, subsurfacetypes.MeanDepthToSlot(dive.DiveComputer.Depth.Mean.Value), &timeSinceDive)
	(*statsContainer).Add(MaxDepth, subsurfacetypes.MaxDepthToSlot(dive.DiveComputer.Depth.Max.Value), &timeSinceDive)
	(*statsContainer).Add(Temperature, subsurfacetypes.TemperatureToSlot(dive.DiveComputer.Temperature.Water.Value), &timeSinceDive)
	diveSiteID := strings.TrimSpace(dive.DiveSiteID)
	(*statsContainer).Add(DiveSite, diveSites.FetchByID(diveSiteID), &timeSinceDive)
	for _, tag := range dive.Tags.Value {
		(*statsContainer).Add(TagStat, tag, &timeSinceDive)
	}
}

func diveSiteReceiver(c chan subsurfacetypes.Divesite, wg *sync.WaitGroup, diveSites *diveSiteMap) {
	for diveSite := range c {
		u := strings.TrimSpace(diveSite.UUID)
		(*diveSites)[u] = diveSite.Name
	}
	wg.Done()
}

func readAndUnmarshal(filename string) subsurfacetypes.Divelog {
	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer xmlFile.Close()
	rawXMLValue, _ := ioutil.ReadAll(xmlFile)
	var divelog subsurfacetypes.Divelog
	err = xml.Unmarshal(rawXMLValue, &divelog)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	return divelog
}

func processDiveSites(divelog *subsurfacetypes.Divelog) diveSiteMap {
	var wg sync.WaitGroup
	diveSites := make(diveSiteMap)
	wg.Add(1)
	c := make(chan subsurfacetypes.Divesite)
	go diveSiteReceiver(c, &wg, &diveSites)
	for _, diveSite := range divelog.Divesites.Site {
		c <- diveSite
	}
	close(c)
	wg.Wait()
	return diveSites
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	divelog := readAndUnmarshal(*filenameFlag)
	diveSites := processDiveSites(&divelog)
	c := make(chan subsurfacetypes.Dive, 100)

	wg.Add(1)
	go diveReceiver(c, &wg, &diveSites)

	for _, trip := range divelog.Dives.Trips {
		for _, dive := range trip.Dives {
			c <- dive
		}
	}
	for _, dive := range divelog.Dives.Dives {
		c <- dive
	}
	close(c)
	wg.Wait()
}
