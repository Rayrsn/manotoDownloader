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
		fmt.Println("Failed to get search results")
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
		fmt.Println("Failed to get search results list")
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
		fmt.Println("Failed to get season results")
		fmt.Println(err)
		os.Exit(1)
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
		fmt.Println("Failed to get episode results list")
		fmt.Println(err)
		os.Exit(1)
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
		fmt.Println("Failed to download episode")
		fmt.Println(err)
		os.Exit(1)
	}
	var episode map[string]interface{}
	json.NewDecoder(response.Body).Decode(&episode)
	episodeDetails := episode["details"].(map[string]interface{})
	episodeTitle := episodeDetails["analyticsEpisodeTitle"].(string)
	m3u8VodUrl := episodeDetails["videoM3u8Url"].(string)
	m3u8VodBaseUrl := strings.Replace(m3u8VodUrl, ".m3u8", "", 1)

	if quality == 1080 {
		downloadWithFFmpeg(m3u8VodBaseUrl, "_2500.m3u8", episodeTitle+"_1080.mp4", "1080")

	} else if quality == 720 {
		downloadWithFFmpeg(m3u8VodBaseUrl, "_750.m3u8", episodeTitle+"_720.mp4", "720")
	}

}

func downloadWithFFmpeg(Url string, Extension string, Filename string, Quality string) {
	if Quality == "1080" {
		switch Extension {
		case "_2500.m3u8":
			ffmpegHandler := ffmpeg.Input(Url + "_2500.m3u8").
				Output(Filename).OverWriteOutput().Run()
			if ffmpegHandler != nil {
				fmt.Printf("Couldn't download with %s quality: %v\n", Extension, ffmpegHandler)
				fmt.Printf("\nTrying with %s quality\n", "1500")
				downloadWithFFmpeg(Url, "_1500.m3u8", Filename, "1080")
			} else {
				fmt.Println("Download Successful!")
			}
		case "_1500.m3u8":
			ffmpegHandler := ffmpeg.Input(Url + "_1500.m3u8").
				Output(Filename).OverWriteOutput().Run()
			if ffmpegHandler != nil {
				fmt.Printf("Couldn't download with %s quality: %v\n", Extension, ffmpegHandler)
				fmt.Printf("\nTrying with %s quality\n", "1000")
				downloadWithFFmpeg(Url, "_1000.m3u8", Filename, "1080")
			} else {
				fmt.Println("Download Successful!")
			}
		case "_1000.m3u8":
			ffmpegHandler := ffmpeg.Input(Url + "_1000.m3u8").
				Output(Filename).OverWriteOutput().Run()
			if ffmpegHandler != nil {
				fmt.Printf("Couldn't download with %s quality: %v\n", Extension, ffmpegHandler)
				fmt.Printf("\nTrying with %s quality\n", "750")
				downloadWithFFmpeg(Url, "_750.m3u8", Filename, "1080")
			} else {
				fmt.Println("Download Successful!")
			}
		case "_750.m3u8":
			ffmpegHandler := ffmpeg.Input(Url + "_750.m3u8").
				Output(Filename).OverWriteOutput().Run()
			if ffmpegHandler != nil {
				fmt.Printf("Couldn't download with %s quality: %v\n", Extension, ffmpegHandler)
				fmt.Printf("\nTrying with %s quality\n", "500")
				downloadWithFFmpeg(Url, "_500.m3u8", Filename, "1080")
			} else {
				fmt.Println("Download Successful!")
			}
		case "_500.m3u8":
			ffmpegHandler := ffmpeg.Input(Url + "_500.m3u8").
				Output(Filename).OverWriteOutput().Run()
			if ffmpegHandler != nil {
				fmt.Printf("Couldn't download with %s quality: %v\n", Extension, ffmpegHandler)
			} else {
				fmt.Println("Download Successful!")
			}
		}
	} else if Quality == "720" {
		switch Extension {
		case "_750.m3u8":
			ffmpegHandler := ffmpeg.Input(Url + "_750.m3u8").
				Output(Filename).OverWriteOutput().Run()
			if ffmpegHandler != nil {
				fmt.Printf("Couldn't download with %s quality: %v\n", Extension, ffmpegHandler)
				fmt.Printf("\nTrying with %s quality\n", "1000")
				downloadWithFFmpeg(Url, "_1000.m3u8", Filename, "720")
			} else {
				fmt.Println("Download Successful!")
			}
		case "_1000.m3u8":
			ffmpegHandler := ffmpeg.Input(Url + "_1000.m3u8").
				Output(Filename).OverWriteOutput().Run()
			if ffmpegHandler != nil {
				fmt.Printf("Couldn't download with %s quality: %v\n", Extension, ffmpegHandler)
				fmt.Printf("\nTrying with %s quality\n", "1500")
				downloadWithFFmpeg(Url, "_1500.m3u8", Filename, "720")
			} else {
				fmt.Println("Download Successful!")
			}
		case "_1500.m3u8":
			ffmpegHandler := ffmpeg.Input(Url + "_1500.m3u8").
				Output(Filename).OverWriteOutput().Run()
			if ffmpegHandler != nil {
				fmt.Printf("Couldn't download with %s quality: %v\n", Extension, ffmpegHandler)
				fmt.Printf("\nTrying with %s quality\n", "2500")
				downloadWithFFmpeg(Url, "_2500.m3u8", Filename, "720")
			} else {
				fmt.Println("Download Successful!")
			}
		case "_2500.m3u8":
			ffmpegHandler := ffmpeg.Input(Url + "_2500.m3u8").
				Output(Filename).OverWriteOutput().Run()
			if ffmpegHandler != nil {
				fmt.Printf("Couldn't download with %s quality: %v\n", Extension, ffmpegHandler)
			} else {
				fmt.Println("Download Successful!")
			}
		}
	}
}
