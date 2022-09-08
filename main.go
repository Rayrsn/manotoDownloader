package main

import (
	"fmt"
	"os"
	"strconv"

	shows "github.com/Rayrsn/manotoDownloader/cmd"

	tm "github.com/buger/goterm"
)

const divider = "•"

func main() {
	clearTerm()
	var searchStringInput string
	fmt.Print("Enter show name: ")
	fmt.Scanln(&searchStringInput)
	searchresult := shows.SearchResultList(searchStringInput)
	fmt.Printf("Found %d shows.\n\n", len(searchresult))
	for i := range searchresult {
		fmt.Printf("%s %s %s\n", shows.GetSlideTitle(searchresult, i), divider, shows.GetSlideId(searchresult, i))
	}

	var showIdInput string
	fmt.Print("\nEnter show ID: ")
	fmt.Scanln(&showIdInput)

	clearTerm()

	fmt.Println("Available seasons:")
	seriesList := shows.GetSeriesList(showIdInput)
	for i := range seriesList {
		seasonIdString := strconv.FormatFloat(seriesList[i].(map[string]interface{})["id"].(float64), 'f', -1, 64)
		seasonNumber := seriesList[i].(map[string]interface{})["analyticsSeasonNumber"]
		fmt.Printf("Season %s, ID: %s.\n", seasonNumber, seasonIdString)
	}

	fmt.Println()

	var seasonIdInput string
	fmt.Print("Enter season ID: ")
	fmt.Scanln(&seasonIdInput)

	seasonId, err := strconv.Atoi(seasonIdInput)
	if err != nil {
		fmt.Println("Invalid input")
		os.Exit(1)
	}

	if seasonIdInput == "" || seasonId == 0 {
		fmt.Println("Invalid input")
		os.Exit(1)
	}

	if len(seasonIdInput) == 1 {
		seasonIdInput = "0" + seasonIdInput
	}

	clearTerm()

	fmt.Printf("Episodes for season %s: \n", seasonIdInput)
	episodeList := shows.GetEpisodeList(seasonIdInput)
	for i := range episodeList.([]interface{}) {
		episodeTitleEng := episodeList.([]interface{})[i].(map[string]interface{})["analyticsEpisodeTitle"]
		episodeTitleFa := episodeList.([]interface{})[i].(map[string]interface{})["formattedEpisodeTitle"]
		episodeTitleID := strconv.FormatFloat(episodeList.([]interface{})[i].(map[string]interface{})["id"].(float64), 'f', -1, 64)

		fmt.Printf("%s %s %s, ID: %s\n", episodeTitleEng, divider, episodeTitleFa, episodeTitleID)
	}

	var episodeIdInput string
	fmt.Print("\nEnter episode ID: ")
	fmt.Scanln(&episodeIdInput)

	var qualityInput string
	fmt.Println("1. 1080p")
	fmt.Println("2. 720p")
	fmt.Print("\nSelect Quality: ")
	fmt.Scanln(&qualityInput)

	clearTerm()
	if qualityInput == "1" || qualityInput == "1080" {
		quality := 1080
		fmt.Printf("Downloading episode %s in %vp quality\n", episodeIdInput, quality)
		shows.DownloadEpisode(episodeIdInput, quality)

	} else if qualityInput == "2" || qualityInput == "720" {
		quality := 720
		fmt.Printf("Downloading episode %s in %vp quality\n", episodeIdInput, quality)
		shows.DownloadEpisode(episodeIdInput, quality)

	} else {
		fmt.Println("Invalid input")
		os.Exit(1)
	}
}

func clearTerm() {
	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Flush()
}
