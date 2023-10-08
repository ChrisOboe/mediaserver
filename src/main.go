package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ChrisOboe/tvhcc/tvhapi"
)

type ImageUrls struct {
	Primary   []string `json:"primary"`
	Clearart  []string `json:"clearart"`
	Banner    []string `json:"banner"`
	Disk      []string `json:"disk"`
	Logo      []string `json:"logo"`
	Miniature []string `json:"miniature"`
	Fanart    []string `json:"fanart"`
}

type Entry struct {
	Provider    string    `json:"provider"`
	Id          string    `json:"id"`
	Channel     string    `json:"channel"`
	Number      int       `json:"number"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	ImageUrls   ImageUrls `json:"imageUrls"`
	MediaUrl    string    `json:"playUrl"`
}

const tvhserver = "http://tv.chump.gob.zone"

func getEntries() ([]Entry, error) {
	tvhapi := tvhapi.Init(tvhserver)
	epg, err := tvhapi.GetEpg()
	if err != nil {
		fmt.Println(err.Error())
		return []Entry{}, err
	}
	entries := make([]Entry, 0)
	for _, epg := range epg.Entries {
		number, _ := strconv.Atoi(epg.ChannelNumber)

		entries = append(entries, Entry{
			Provider:    "tvheadend",
			Id:          string(epg.ChannelUuid),
			Channel:     epg.ChannelName,
			Number:      number,
			Name:        epg.Title,
			Tags:        []string{"live"},
			Description: epg.Description,
			MediaUrl:    tvhapi.GetStream(epg.ChannelUuid),
			ImageUrls: ImageUrls{
				Logo: []string{epg.ChannelIcon},
			},
		})
	}

	return entries, nil
}

func main() {
	http.HandleFunc("/entries", jsonHandler)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	entries, err := getEntries()
	if err != nil {
		fmt.Println(err.Error())
	}
	json.NewEncoder(w).Encode(entries)
}
