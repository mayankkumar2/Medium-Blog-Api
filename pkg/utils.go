package pkg
import (
	"net/http"
	"io/ioutil"
	"github.com/mmcdole/gofeed/rss"
	"strings"
	"regexp"
	"encoding/json"
)
type Article struct{
	Title string `json:"title"`
	Link string `json:"link"`
	ThumbnailRef string `json:"thumbnailref"`
	Creators string `json:"creators"`
	PubDate string `json:"PublishDateTime"`
}
type Articles struct{
	Status string `json:"status"`
	Arts []Article `json:"articles"`
}
func GetArticles(userId string) []byte {
	url := "https://medium.com/feed/" + userId
	data,_ := http.Get(url)
	defer data.Body.Close()
	if data.StatusCode == 404 {
		resp, _ := json.Marshal(Articles{
			Status: "fail",
		})
		return resp
	}
	body, _ := ioutil.ReadAll(data.Body)
	bodyAsString := string(body)
	fp := rss.Parser{}
	rssFeed, _ := fp.Parse(strings.NewReader(bodyAsString))
	articles := make([]Article,0,10000)
	for _, value := range (rssFeed.Items) {
		re := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
		imageRefs := re.FindAllStringSubmatch(value.Content, -1)
		creators := ""
		for x,val := range value.Extensions["dc"]["creator"] {
			creators += val.Value
			if x != len(value.Extensions["dc"]["creator"])-1{
				creators += ","
			}
		}
		articles = append(articles, Article{Title: value.Title,Creators: creators, Link: value.Link, ThumbnailRef: imageRefs[0][1],PubDate:value.PubDate})
	}
	arts := Articles{
		Status: "success",
		Arts: articles,
	}
	responseJson,_ := json.Marshal(arts)
	return responseJson
}