package subsurfacetypes

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Divelog is a top level XML from subsurface.
type Divelog struct {
	XMLName   xml.Name  `xml:"divelog"`
	Divesites Divesites `xml:"divesites"`
	Dives     Dives     `xml:"dives"`
}

// Divesites holds generic information about each divesite
type Divesites struct {
	XMLName xml.Name   `xml:"divesites"`
	Site    []Divesite `xml:"site"`
}

// Divesite describes a single dive site with no information about related dives.
type Divesite struct {
	XMLName     xml.Name `xml:"site"`
	UUID        string   `xml:"uuid,attr"`
	Name        string   `xml:"name,attr"`
	GPS         string   `xml:"gps,attr"`
	Description string   `xml:"description,attr"`
	Notes       string   `xml:"notes"`
}

// Dives is a container for list of dives and trips.
type Dives struct {
	XMLName xml.Name `xml:"dives"`
	Dives   []Dive   `xml:"dive"`
	Trips   []Trip   `xml:"trip"`
}

func (d Dives) String() string {
	return fmt.Sprintf("Dives (%v, trips %v)", len(d.Dives), len(d.Trips))
}

type SubsurfaceTime struct {
	Value time.Time
}

func (t *SubsurfaceTime) UnmarshalXMLAttr(attr xml.Attr) error {
	const timeFormat = "15:04:05"
	parsedValue, err := time.Parse(timeFormat, attr.Value)
	if err != nil {
		return err
	}
	*t = SubsurfaceTime{parsedValue}
	return nil
}

type SubsurfaceDate struct {
	Value time.Time
}

func (t *SubsurfaceDate) UnmarshalXMLAttr(attr xml.Attr) error {
	const dateFormat = "2006-01-02"
	parsedValue, err := time.Parse(dateFormat, attr.Value)
	if err != nil {
		return err
	}
	*t = SubsurfaceDate{parsedValue}
	return nil
}

func (t SubsurfaceTime) Duration() time.Duration {
	return time.Duration(time.Duration(t.Value.Hour())*time.Hour + time.Duration(t.Value.Minute())*time.Minute + time.Duration(t.Value.Second())*time.Second)
}

// Trip is a collection of dives.
type Trip struct {
	//Date     SubsurfaceDate `xml:"date,attr"`
	Time     string `xml:"time,attr"`
	Location string `xml:"location,attr"`
	Dives    []Dive `xml:"dive"`
}

// Dive has information about a single dive.
type Dive struct {
	XMLName      xml.Name       `xml:"dive"`
	Number       string         `xml:"number,attr"`
	Tags         Tags           `xml:"tags,attr"`
	DiveSiteID   string         `xml:"divesiteid,attr"`
	Date         SubsurfaceDate `xml:"date,attr"`
	Time         SubsurfaceTime `xml:"time,attr"`
	RawDuration  string         `xml:"duration,attr"`
	Buddy        string         `xml:"buddy"`
	Cylinders    []Cylinder     `xml:"cylinder"`
	Invalid      string         `xml:"invalid,attr"`
	DiveComputer DiveComputer   `xml:"divecomputer"`
	//Notes           string          `xml:"notes"`
}

type Tags struct {
	Value []string
}

func (t *Tags) UnmarshalXMLAttr(attr xml.Attr) error {
	tags := strings.Split(attr.Value, ", ")
	*t = Tags{tags}
	return nil
}

func (d Dive) IsInvalid() bool {
	return d.Invalid == "1"
}

// DiveComputer holds information imported from a dive computer.
type DiveComputer struct {
	XMLName     xml.Name        `xml:"divecomputer"`
	Model       string          `xml:"model,attr"`
	Depth       DiveDepth       `xml:"depth"`
	Temperature DiveTemperature `xml:"temperature"`
}

type DepthReading struct {
	Value float64
}

func (d *DepthReading) UnmarshalXMLAttr(attr xml.Attr) error {
	if !strings.HasSuffix(attr.Value, " m") {
		fmt.Println("Invalid depth:", attr.Value)
		return nil
	}
	r := strings.Split(attr.Value, " ")
	val, _ := strconv.ParseFloat(r[0], 64)
	*d = DepthReading{val}
	return nil
}

// DiveDepth has information about max and mean depth for a single dive.
type DiveDepth struct {
	XMLName xml.Name     `xml:"depth"`
	Max     DepthReading `xml:"max,attr"`
	Mean    DepthReading `xml:"mean,attr"`
}

// TimeSince returns duration since dive was logged
func (d *Dive) TimeSince() time.Duration {
	diveDate := d.Date.Value.Add(d.Time.Duration())
	return time.Since(diveDate)
}

// BuddyList returns a list of buddies (or empty list)
func (d *Dive) BuddyList() []string {
	splitBuddies := strings.Split(d.Buddy, ",")
	for i := 0; i < len(splitBuddies); i++ {
		splitBuddies[i] = strings.Trim(splitBuddies[i], " ")
	}
	return splitBuddies
}

// Duration returns parsed dive duration
func (d *Dive) Duration() time.Duration {
	if strings.HasSuffix(d.RawDuration, " min") {
		a := strings.Split(d.RawDuration, " ")
		b := strings.Split(a[0], ":")
		secondsInt, err := strconv.Atoi(b[1])
		var secondsFraction float64
		if err == nil {
			secondsFraction = float64(secondsInt) / 60.0
		}
		minutesInt, _ := strconv.Atoi(b[0])
		durationFraction := float64(minutesInt) + secondsFraction
		duration, _ := time.ParseDuration(fmt.Sprintf("%.5f", durationFraction) + "m")
		return duration
	}
	zeroDuration, _ := time.ParseDuration("0s")
	return zeroDuration
}

// Cylinder has information about cylinders used on the dive.
type Cylinder struct {
	XMLName      xml.Name `xml:"cylinder"`
	Size         string   `xml:"size,attr"`
	WorkPressure string   `xml:"workpressure,attr"`
	Description  string   `xml:"description,attr"`
	O2           string   `xml:"o2,attr"`
	He           string   `xml:"he,attr"`
	Start        string   `xml:"start,attr"`
	End          string   `xml:"end,attr"`
	Depth        string   `xml:"depth,attr"`
}

// DiveTemperature has water and air temperature information.
type DiveTemperature struct {
	XMLName xml.Name    `xml:"temperature"`
	Water   Temperature `xml:"water,attr"`
}

type Temperature struct {
	Value float64
}

func (t *Temperature) UnmarshalXMLAttr(attr xml.Attr) error {
	if !strings.HasSuffix(attr.Value, " C") {
		fmt.Println("Invalid water temperature:", attr.Value)
		return nil
	}
	r := strings.Split(attr.Value, " ")
	convertedTemperature, _ := strconv.ParseFloat(r[0], 64)
	*t = Temperature{convertedTemperature}
	return nil
}
