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
	Program   string    `xml:"program,attr"`
	Version   string    `xml:"version,attr"`
	Settings  Settings  `xml:"settings"`
	Divesites Divesites `xml:"divesites"`
	Dives     Dives     `xml:"dives"`
}

// Settings has general per-divelog settings, such as dive computer info.
type Settings struct {
	XMLName        xml.Name         `xml:"settings"`
	DiveComputerID []DiveComputerID `xml:"divecomputerid"`
}

// DiveComputerID is per-log information about a specific dive computer
type DiveComputerID struct {
	XMLName  xml.Name `xml:"divecomputerid"`
	Model    string   `xml:"model,attr"`
	DeviceID string   `xml:"deviceid,attr"`
	Serial   string   `xml:"serial,attr"`
	Firmware string   `xml:"firmware,attr"`
}

// Divesites holds generic information about each divesite
type Divesites struct {
	XMLName xml.Name   `xml:"divesites"`
	Site    []Divesite `xml:"site"`
}

// Divesite describes a single dive site with no information about related dives.
type Divesite struct {
	XMLName     xml.Name      `xml:"site"`
	UUID        string        `xml:"uuid,attr"`
	Name        string        `xml:"name,attr"`
	GPS         string        `xml:"gps,attr"`
	Description string        `xml:"description,attr"`
	Notes       string        `xml:"notes"`
	Geo         []DivesiteGEO `xml:"geo"`
}

// DivesiteGEO holds category information for dive sites.
type DivesiteGEO struct {
	XMLName xml.Name `xml:"geo"`
	Cat     string   `xml:"cat,attr"`
	Origin  string   `xml:"origin,attr"`
	Value   string   `xml:"value,attr"`
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

// SubsurfaceTime holds parsed time information
type SubsurfaceTime struct {
	Value time.Time
}

// UnmarshalXMLAttr Parses XML attribute to time
func (t *SubsurfaceTime) UnmarshalXMLAttr(attr xml.Attr) error {
	const timeFormat = "15:04:05"
	parsedValue, err := time.Parse(timeFormat, attr.Value)
	if err != nil {
		return err
	}
	*t = SubsurfaceTime{parsedValue}
	return nil
}

// MarshalXMLAttr outputs parsed time object to a string
func (t *SubsurfaceTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: t.Value.Format("15:04:05")}, nil
}

// SubsurfaceDate holds parsed date object
type SubsurfaceDate struct {
	Value time.Time
}

// UnmarshalXMLAttr Parses XML attribute to date
func (t *SubsurfaceDate) UnmarshalXMLAttr(attr xml.Attr) error {
	const dateFormat = "2006-01-02"
	parsedValue, err := time.Parse(dateFormat, attr.Value)
	if err != nil {
		return err
	}
	*t = SubsurfaceDate{parsedValue}
	return nil
}

// MarshalXMLAttr formats parsed date back to string
func (t *SubsurfaceDate) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: t.Value.Format("2006-01-02")}, nil
}

func (t SubsurfaceTime) Duration() time.Duration {
	return time.Duration(time.Duration(t.Value.Hour())*time.Hour + time.Duration(t.Value.Minute())*time.Minute + time.Duration(t.Value.Second())*time.Second)
}

// Trip is a collection of dives.
type Trip struct {
	Date     string `xml:"date,attr"`
	Time     string `xml:"time,attr"`
	Location string `xml:"location,attr"`
	Dives    []Dive `xml:"dive"`
	Notes    string `xml:"notes"`
}

// Dive has information about a single dive.
type Dive struct {
	XMLName         xml.Name              `xml:"dive"`
	TripFlag        string                `xml:"tripflag,attr,omitempty"`
	Divemaster      string                `xml:"divemaster"`
	Number          string                `xml:"number,attr"`
	Tags            Tags                  `xml:"tags,attr,omitempty"`
	DiveSiteID      string                `xml:"divesiteid,attr,omitempty"`
	Date            SubsurfaceDate        `xml:"date,attr,omitempty"`
	Time            SubsurfaceTime        `xml:"time,attr,omitempty"`
	RawDuration     string                `xml:"duration,attr,omitempty"`
	Buddy           string                `xml:"buddy"`
	Cylinders       []Cylinder            `xml:"cylinder"`
	Invalid         string                `xml:"invalid,attr,omitempty"`
	DiveTemperature ManualDiveTemperature `xml:"divetemperature"`
	DiveComputer    DiveComputer          `xml:"divecomputer"`
	Rating          string                `xml:"rating,attr,omitempty"`
	CNS             string                `xml:"cns,attr,omitempty"`
	SAC             string                `xml:"sac,attr,omitempty"`
	Notes           string                `xml:"notes"`
	OTU             string                `xml:"otu,attr,omitempty"`
	Visibility      string                `xml:"visibility,attr,omitempty"`
	Current         string                `xml:"current,attr,omitempty"`
	Suit            string                `xml:"suit"`
	WeightSystem    []WeightSystem        `xml:"weightsystem"`
}

// ManualDiveTemperature holds manually added dive temperature information
type ManualDiveTemperature struct {
	XMLName xml.Name `xml:"divetemperature"`
	Water   string   `xml:"water,attr,omitempty"`
	Air     string   `xml:"air,attr,omitempty"`
}

// WeightSystem has weight system information (weights, where those were deployed to)
type WeightSystem struct {
	XMLName     xml.Name `xml:"weightsystem"`
	Weight      string   `xml:"weight,attr,omitempty"`
	Description string   `xml:"description,attr,omitempty"`
}

// Tags is a list of tags entered by user
type Tags struct {
	Value []string
}

func (t *Tags) UnmarshalXMLAttr(attr xml.Attr) error {
	tags := strings.Split(attr.Value, ", ")
	*t = Tags{tags}
	return nil
}

func (t *Tags) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: strings.Join(t.Value, ", ")}, nil
}

func (d Dive) IsInvalid() bool {
	return d.Invalid == "1"
}

// DiveComputer holds information imported from a dive computer.
type DiveComputer struct {
	XMLName        xml.Name        `xml:"divecomputer"`
	Model          string          `xml:"model,attr,omitempty"`
	LastManualTime string          `xml:"last-manual-time,attr,omitempty"`
	ManualDate     string          `xml:"date,attr,omitempty"`
	ManualTime     string          `xml:"time,attr,omitempty"`
	Depth          DiveDepth       `xml:"depth"`
	Temperature    DiveTemperature `xml:"temperature"`
	DeviceID       string          `xml:"deviceid,attr,omitempty"`
	DiveID         string          `xml:"diveid,attr,omitempty"`
	Surface        Surface         `xml:"surface"`
	Events         []DiveEvent     `xml:"event"`
	Samples        []DiveSample    `xml:"sample"`
	ExtraData      []ExtraData     `xml:"extradata"`
	Water          WaterDetails    `xml:"water"`
}

// WaterDetails contains information about water
type WaterDetails struct {
	XMLName  xml.Name `xml:"water"`
	Salinity string   `xml:"salinity,attr,omitempty"`
}

// ExtraData describes any unstructured values provided by the dive computer.
type ExtraData struct {
	XMLName xml.Name `xml:"extradata"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

// DiveEvent is a specific event not describe by samples, such as gas changes.
type DiveEvent struct {
	XMLName  xml.Name `xml:"event"`
	Time     string   `xml:"time,attr,omitempty"`
	Type     string   `xml:"type,attr,omitempty"`
	Flags    string   `xml:"flags,attr,omitempty"`
	Name     string   `xml:"name,attr,omitempty"`
	Cylinder string   `xml:"cylinder,attr,omitempty"`
	Value    string   `xml:"value,attr,omitempty"`
}

// DiveSample is a sample provided by the dive computer. Only time is a mandatory field; everything else is optional
type DiveSample struct {
	XMLName     xml.Name `xml:"sample"`
	Time        string   `xml:"time,attr"`
	Depth       string   `xml:"depth,attr,omitempty"`
	Temperature string   `xml:"temp,attr,omitempty"`
	Pressure    string   `xml:"pressure,attr,omitempty"`
	RBT         string   `xml:"rbt,attr,omitempty"`
	NDL         string   `xml:"ndl,attr,omitempty"`
	CNS         string   `xml:"cns,attr,omitempty"`
	StopTime    string   `xml:"stoptime,attr,omitempty"`
	StopDepth   string   `xml:"stopdepth,attr,omitempty"`
	InDeco      string   `xml:"in_deco,attr,omitempty"`
}

// Surface contains the surface pressure.
type Surface struct {
	XMLName  xml.Name `xml:"surface"`
	Pressure string   `xml:"pressure,attr,omitempty"`
}

// DepthReading is a parsed depth reading
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

func (d *DepthReading) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{Name: name, Value: fmt.Sprintf("%f m", d.Value)}, nil
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
	Size         string   `xml:"size,attr,omitempty"`
	WorkPressure string   `xml:"workpressure,attr,omitempty"`
	Description  string   `xml:"description,attr,omitempty"`
	O2           string   `xml:"o2,attr,omitempty"`
	He           string   `xml:"he,attr,omitempty"`
	Start        string   `xml:"start,attr,omitempty"`
	End          string   `xml:"end,attr,omitempty"`
	Depth        string   `xml:"depth,attr,omitempty"`
}

// DiveTemperature has water and air temperature information.
type DiveTemperature struct {
	XMLName xml.Name    `xml:"temperature"`
	Water   Temperature `xml:"water,attr,omitempty"`
	Air     Temperature `xml:"air,attr,omitempty"`
}

// Temperature holds temperature information, including whether temperature was valid (in order to avoid outputting 0 C).
type Temperature struct {
	Value float64
	Valid bool
}

// UnmarshalXMLAttr parses temperature information. Only celsius is supported.
func (t *Temperature) UnmarshalXMLAttr(attr xml.Attr) error {
	if !strings.HasSuffix(attr.Value, " C") {
		fmt.Println("Invalid water temperature:", attr.Value)
		return nil
	}
	r := strings.Split(attr.Value, " ")
	convertedTemperature, _ := strconv.ParseFloat(r[0], 64)
	*t = Temperature{convertedTemperature, true}
	return nil
}

// MarshalXMLAttr outputs temperature information back to XML. Only celsius is supported.
func (t *Temperature) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if t.Valid {
		return xml.Attr{Name: name, Value: fmt.Sprintf("%f C", t.Value)}, nil
	}
	return xml.Attr{Name: name, Value: ""}, nil
}
