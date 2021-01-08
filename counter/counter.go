package counter

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

type lastCounterStat struct {
	Name       string
	Count      int
	SinceLast  time.Duration
	SinceFirst time.Duration
}

// statSorter joins a SortBy function and a slice of LastCounterStat to be sorted.
type statSorter struct {
	stats []lastCounterStat
	by    func(p1, p2 *lastCounterStat) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *statSorter) Len() int {
	return len(s.stats)
}

// Swap is part of sort.Interface.
func (s *statSorter) Swap(i, j int) {
	s.stats[i], s.stats[j] = s.stats[j], s.stats[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *statSorter) Less(i, j int) bool {
	return s.by(&s.stats[i], &s.stats[j])
}

// LastCounterStats holds information regarding last occurrence of specified event
type LastCounterStats map[string]*lastCounterStat

// LastCounter keeps track of occurrences and last time something happened
type lastCounter interface {
	Add(name string, timeSince *time.Duration)
	PrintStats()
}

// SortBy implements selecting a correct field for sorting.
type SortBy func(d1, d2 *lastCounterStat) bool

func formatDurationToDays(duration time.Duration) string {
	return fmt.Sprintf("%.0f", duration.Hours()/24.0)
}

// Sort is a method on the function type, SortBy, that sorts the argument slice according to the function.
func (sortBy SortBy) Sort(stats []lastCounterStat) {
	ps := &statSorter{
		stats: stats,
		by:    sortBy,
	}
	sort.Sort(ps)
}

// Add adds a new instance to the counter.
func (p LastCounterStats) Add(name string, timeSince *time.Duration) {
	_, ok := p[name]
	if !ok {
		p[name] = &lastCounterStat{name, 0, *timeSince, *timeSince}
	}
	if *timeSince < p[name].SinceLast {
		p[name].SinceLast = *timeSince
	}
	if *timeSince > p[name].SinceFirst {
		p[name].SinceFirst = *timeSince
	}
	p[name].Count++

}

// PrintStats prints tabulated statistics to stdout
func (p LastCounterStats) PrintStats(sortBy string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Nimi", "Kertoja", "Edellinen päivää sitten", "Ensimmäinen päivää sitten"})
	t.AppendSeparator()
	sl := make([]lastCounterStat, len(p))
	i := 0
	for _, stat := range p {
		sl[i] = *stat
		i++
	}
	nameSort := func(s1, s2 *lastCounterStat) bool {
		return s1.Name < s2.Name
	}
	countSort := func(s1, s2 *lastCounterStat) bool {
		return s1.Count < s2.Count
	}
	sinceFirstSort := func(s1, s2 *lastCounterStat) bool {
		return s1.SinceFirst < s2.SinceFirst
	}
	sinceLastSort := func(s1, s2 *lastCounterStat) bool {
		return s1.SinceLast < s2.SinceLast
	}
	switch sortBy {
	case "name":
		SortBy(nameSort).Sort(sl)
	case "count":
		SortBy(countSort).Sort(sl)
	case "sinceFirst":
		SortBy(sinceFirstSort).Sort(sl)
	case "sinceLast":
		SortBy(sinceLastSort).Sort(sl)
	default:
		fmt.Println("Invalid sort flag", sortBy, ". Showing entries in random order.")
	}
	for i, stat := range sl {
		t.AppendRow([]interface{}{i + 1, stat.Name, stat.Count, formatDurationToDays(stat.SinceLast), formatDurationToDays(stat.SinceFirst)})
	}
	t.Render()
	fmt.Println("Yhteensä", len(p))
}
