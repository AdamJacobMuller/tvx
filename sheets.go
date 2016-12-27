package main

import (
	"fmt"
	"strings"

	"github.com/garfunkel/go-tvdb"

	"gopkg.in/Iwark/spreadsheet.v2"
)

func main() {
	service, err := spreadsheet.NewService()
	if err != nil {
		panic(err)
	}

	spreadsheet, err := service.FetchSpreadsheet("1t9n09fuiQhjQOr2EbsHzVp-GKakLcfjr94jGniSh8Pw")
	if err != nil {
		panic(err)
	}

	show, err := spreadsheet.SheetByTitle("Shows")
	if err != nil {
		panic(err)
	}

	_ = fmt.Printf

	years := map[string]uint{}
	shows := map[uint]string{}
	imdbs := map[uint]string{}

	for _, r := range show.Rows[0][1:] {
		shows[r.Column] = r.Value
	}
	for _, r := range show.Rows[1][1:] {
		imdbs[r.Column] = r.Value
	}
	for _, r := range show.Columns[0][2:] {
		if r.Row == 0 {
			fmt.Printf("bad row %+v\n", r)
			continue
		}
		years[r.Value] = r.Row
	}
	fmt.Printf("shows: %+v\n", shows)
	fmt.Printf("imdbs: %+v\n", imdbs)
	fmt.Printf("years: %+v\n", years)
	for k := range shows {
		yearcount := map[string]int{}
		name := shows[k]
		imdb := imdbs[k]
		fmt.Printf("show %s (%s)\n", name, imdb)
		db, err := tvdb.GetSeriesByIMDBID(imdb)
		if err != nil {
			fmt.Printf("GetSeriesByIMDBID err = %s\n", err)
			continue
		}
		show.Update(0, int(k), db.SeriesName)
		err = db.GetDetail()
		if err != nil {
			fmt.Printf("GetDetail err = %s\n", err)
			continue
		}
		for k, season := range db.Seasons {
			for l, episode := range season {
				fmt.Printf("k:%+v l:%+v s:%+v\n", k, l, episode.EpisodeName)
				parts := strings.Split(episode.FirstAired, "-")
				if len(parts) == 0 {
					fmt.Printf("invalid date: %+v\n", episode.FirstAired)
					continue
				}
				_, ok := yearcount[parts[0]]
				if ok {
					yearcount[parts[0]] += 1
				} else {
					yearcount[parts[0]] = 1
				}
			}
		}
		for year, count := range yearcount {
			yearrow, ok := years[year]
			if !ok {
				fmt.Printf("missing year %s\n", year)
				continue
			}
			fmt.Printf("up (%d, %d) = %d\n", yearrow, k, count)
			show.Update(int(yearrow), int(k), fmt.Sprintf("%d", count))
		}

		err = show.Synchronize()
		if err != nil {
			panic(err)
		}
	}
}
