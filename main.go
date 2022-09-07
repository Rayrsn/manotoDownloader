package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	tm "github.com/buger/goterm"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const searchApiUrl = "https://dt6ka97rrh6d5.cloudfront.net/api/v1/publicrole/searchmodule/searchitems"
const seriesListApiUrl = "https://dt6ka97rrh6d5.cloudfront.net/api/v1/publicrole/showmodule/serieslist?id="
const episodeListApiUrl = "https://dt6ka97rrh6d5.cloudfront.net/api/v1/publicrole/showmodule/episodelist?id="
const episodeDetailApiUrl = "https://dt6ka97rrh6d5.cloudfront.net/api/v1/publicrole/showmodule/episodedetails?id="

const divider = "•"

func main() {
	clearTerm()
	var searchStringInput string
	fmt.Print("Enter show name: ")
	fmt.Scanln(&searchStringInput)
	searchresult := SearchResultList(searchStringInput)
	fmt.Printf("Found %d shows.\n\n", len(searchresult))
	for i := range searchresult {
		fmt.Printf("%s %s %s\n", getSlideTitle(searchresult, i), divider, getSlideId(searchresult, i))
	}

	var showIdInput string
	fmt.Print("\nEnter show ID: ")
	fmt.Scanln(&showIdInput)

	clearTerm()

	fmt.Println("Available seasons:")
	seriesList := GetSeriesList(showIdInput)
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
	episodeList := getEpisodeList(seasonIdInput)
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
	fmt.Println("Select Quality:")
	fmt.Println("1. 1080p")
	fmt.Println("2. 720p")
	fmt.Scanln(&qualityInput)

	clearTerm()
	if qualityInput == "1" || qualityInput == "1080" {
		quality := 1080
		fmt.Printf("Downloading episode %s in %vp quality\n", episodeIdInput, quality)
		DownloadEpisode(episodeIdInput, quality)

	} else if qualityInput == "2" || qualityInput == "720" {
		quality := 720
		fmt.Printf("Downloading episode %s in %vp quality\n", episodeIdInput, quality)
		DownloadEpisode(episodeIdInput, quality)

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

func SearchForResults(searchString string) map[string]interface{} {
	payload := strings.NewReader(fmt.Sprintf(`{"SearchText":"%s","SlideType":"show","PageNumber":1,"PageSize":10,"SortBy":"az"}`, searchString))
	req, _ := http.NewRequest("POST", searchApiUrl, payload)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer res.Body.Close()
	var results map[string]interface{}
	json.NewDecoder(res.Body).Decode(&results)
	return results
}

func SearchResultList(searchString string) []interface{} {
	results := SearchForResults(searchString)
	var details = results["details"].(map[string]interface{})
	var list = details["list"].([]interface{})
	if len(list) == 0 {
		fmt.Println("No results found")
		os.Exit(1)
	}
	return list
}

func getSlideTitle(searchlist []interface{}, index int) string {
	return searchlist[index].(map[string]interface{})["SlideTitle"].(string)
}

func getSlideId(searchlist []interface{}, index int) string {
	return searchlist[index].(map[string]interface{})["SlideId"].(string)
}

func GetSeriesList(showId string) []interface{} {
	url := seriesListApiUrl + showId
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	var episodes map[string]interface{}
	json.NewDecoder(response.Body).Decode(&episodes)
	seriesList := episodes["details"].(map[string]interface{})["list"]
	return seriesList.([]interface{})
}

func GetSeriesId(seriesList interface{}, seasonNumber string) string {
	for i := range seriesList.([]interface{}) {
		if seriesList.([]interface{})[i].(map[string]interface{})["analyticsSeasonNumber"] == seasonNumber {
			return seriesList.([]interface{})[i].(map[string]interface{})["id"].(string)
		}
		fmt.Println(seriesList.([]interface{})[i].(map[string]interface{})["analyticsSeasonNumber"])
	}
	return "0"
}

func getEpisodeList(seriesId string) interface{} {
	url := episodeListApiUrl + seriesId
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	var episode map[string]interface{}
	json.NewDecoder(response.Body).Decode(&episode)
	episodeList := episode["details"].(map[string]interface{})["list"]
	return episodeList
}

func GetEpisodes(slideId string, index int) interface{} {
	seriesList := GetSeriesList(slideId)

	indexStr := strconv.Itoa(index)

	seriesId := GetSeriesId(seriesList, indexStr)
	episodeList := getEpisodeList(seriesId)
	return episodeList
}

func DownloadEpisode(episodeId string, quality int) {
	url := episodeDetailApiUrl + episodeId
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	var episode map[string]interface{}
	json.NewDecoder(response.Body).Decode(&episode)
	episodeDetails := episode["details"].(map[string]interface{})
	episodeTitle := episodeDetails["analyticsEpisodeTitle"].(string)
	m3u8VodUrl := episodeDetails["videoM3u8Url"].(string)
	m3u8VodBaseUrl := strings.Replace(m3u8VodUrl, ".m3u8", "", 1)
	m3u81080pUrl := m3u8VodBaseUrl + "_2500.m3u8"
	m3u8720pUrl := m3u8VodBaseUrl + "_750.m3u8"

	if quality == 1080 {
		ffmpegHandler := ffmpeg.Input(m3u81080pUrl).
			Output(episodeTitle + "_1080.mp4").OverWriteOutput().Run()
		if ffmpegHandler != nil {
			fmt.Println(ffmpegHandler)
		}
	} else if quality == 720 {
		ffmpegHandler := ffmpeg.Input(m3u8720pUrl).
			Output(episodeTitle + "_720.mp4").OverWriteOutput().Run()
		if ffmpegHandler != nil {
			fmt.Println(ffmpegHandler)
		}
	}

}
