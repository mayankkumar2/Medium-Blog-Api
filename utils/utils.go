package utils
import (
	"net/http"
	"io/ioutil"
	"github.com/mmcdole/gofeed/rss"
	"sort"
	"strings"
	"regexp"
	"encoding/json"
	"github.com/gocolly/colly"
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
func GetArticles(userId string) (resp []byte,statusCode int) {
	url := "https://medium.com/feed/" + userId
	data,_ := http.Get(url)
	defer data.Body.Close()
	if data.StatusCode != 200 {
		resp, _ := json.Marshal(Articles{
			Status: "fail",
		})
		return resp,data.StatusCode
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
	return responseJson,200
}
type User struct {
	Name string `json:"name"`
	UserName string `json:"username"`
	Avatar string `json:"avatar"`
}
type Post struct {
	Creator User `json:"creators"`
	Link string	`json:"link"`
	Title string `json:"title"`
	Subtitle string `json:"subtitle"`
	Timestamp int64 `json:"PublishDateTime"`
	Thumbnail string `json:"thumbnailref"`
}
func GetArticlesforOrganizationV2(orgId string) (respose interface{}, statusCode int) {
	col := colly.NewCollector()
	var pos *[]Post
	statusCode = 200
	col.OnError(func(e *colly.Response, err error) {
		statusCode = e.StatusCode
	})
	col.OnHTML("script", func (e *colly.HTMLElement) {
		v := len(e.Text)
		if v > 32 {
			if ("// <![CDATA[\nwindow[\"obvInit\"]({") ==(e.Text[:32]) {
				var l map[string] interface{}
				json.Unmarshal([]byte(e.Text[31:v-8]),&l)

				rootMap := (l["references"].(map[string] interface{})) //Root interface
				_users := rootMap["User"].(map[string] interface{}) //JSON interface
				users := make(map[string] User)
				for k := range _users {
					_user := _users[k].(map[string] interface{})
					user := User{
						Name:_user["name"].(string),
						UserName: _user["username"].(string),
						Avatar: "https://cdn-images-1.medium.com/" +  _user["imageId"].(string),
					}
					users[k] = user
				}
				_posts := rootMap["Post"].(map[string] interface{})
				posts := make([]Post,0,len(_posts))
				for k := range _posts {
					//	//post := Posts{}
					_post := _posts[k].(map[string] interface{})
					usr := users[_post["creatorId"].(string)]
					title := _post["title"].(string)
					link := "https://medium.com/"+orgId+"/"+_post["uniqueSlug"].(string)
					//subtitle := _post["subtitle"].(string)
					virtuals := _post["virtuals"].(map[string] interface{})
					prevImage := virtuals["previewImage"].(map[string] interface{})
					thumbnail := "https://cdn-images-1.medium.com/" + prevImage["imageId"].(string)
					subtitle := virtuals["subtitle"].(string)
					timestamp := int64(_post["createdAt"].(float64))
					post := Post{
						Creator: usr,
						Title: title,
						Link: link,
						Thumbnail: thumbnail,
						Subtitle: subtitle,
						Timestamp: timestamp,
					}
					posts = append(posts, post)
				}
				pos = &posts;
			}
		}

	})
	col.Visit("https://medium.com/"+orgId)
	if statusCode != 200 {
		return map[string] string {
			"status": "fail",
			"error": "An error occured.",
		}, statusCode
	} else {
		sort.Slice(*pos, func (i,j int) bool {
			return (*pos)[i].Timestamp > (*pos)[j].Timestamp
		})
		return (map[string] interface{} {
			"status" : "success",
			"articles" : *pos,
		}), statusCode
	}
}