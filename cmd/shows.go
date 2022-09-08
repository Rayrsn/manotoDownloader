package shows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const searchApiUrl = "https://dt6ka97rrh6d5.cloudfront.net/api/v1/publicrole/searchmodule/searchitems"
const seriesListApiUrl = "https://dt6ka97rrh6d5.cloudfront.net/api/v1/publicrole/showmodule/serieslist?id="
const episodeListApiUrl = "https://dt6ka97rrh6d5.cloudfront.net/api/v1/publicrole/showmodule/episodelist?id="
const episodeDetailApiUrl = "https://dt6ka97rrh6d5.cloudfront.net/api/v1/publicrole/showmodule/episodedetails?id="

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

func GetSlideTitle(searchlist []interface{}, index int) string {
	return searchlist[index].(map[string]interface{})["SlideTitle"].(string)
}

func GetSlideId(searchlist []interface{}, index int) string {
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

func GetEpisodeList(seriesId string) interface{} {
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
	episodeList := GetEpisodeList(seriesId)
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
