package main

import (
	"log"
	"net/http"
	"os"
	"math/rand"
	"github.com/gin-gonic/gin"
	"time"
	"strings"
	"io/ioutil"
)

var count = 0

func download(){
	UpdateCheck()
	DownloadPlugin()

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan bool)
	for {
		select {
		case  <- ticker.C:
			UpdateCheck()
			DownloadPlugin()
			count++
		case <- quit:
			ticker.Stop()
			return
		}
	}
}

type Plugin struct {
	PluginSlug  string `json:"Plugin Slug"`
	Name        string `json:"Name"`
	PluginURI   string `json:"PluginURI"`
	Version     string `json:"Version"`
	Description string `json:"Description"`
}

type PluginRequestField struct {
	Plugins string `json:"plugins"`
	Active  []string          `json:"active"`
}

const PLUGIN_UPDATE_CHECK_ENDPOINT string = "http://api.wordpress.org/plugins/update-check/1.1/"
const PLUGIN_DOWNLOAD_ENDPOINT string = "https://downloads.wordpress.org/plugin/wp-vn-oembed.zip"

func GenerateRandomDomain(strlen int) (string) {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return "http://" + string(result) + ".vn"
}

func DownloadPlugin() {
	log.Println("Start Download Plugin.")
	req, _ := http.NewRequest("GET", PLUGIN_DOWNLOAD_ENDPOINT, nil)
	req.Header.Add("Range", "bytes=0-1023")
	var client http.Client
	client.Do(req)
	log.Println("Download finish.")
}

func UpdateCheck() {
	payload := strings.NewReader(`plugins={"plugins":{"wp-vn-oembed\/plugin.php":{"Name":"Wordpress VN oEmbed","PluginURI":"http:\/\/laptrinh.senviet.org","Version":"1.1.0","Description":"T\u1ef1 \u0111\u1ed9ng nh\u00fang player cho c\u00e1c trang nh\u1ea1c \u1edf Vi\u1ec7t Nam. B\u1ea1n c\u00f3 th\u1ec3 xem h\u01b0\u1edbng d\u1eabn t\u1ea1i \u0111\u00e2y <a href=\"http:\/\/laptrinh.senviet.org\/wordpress-plugin\/wordpress-oembed-ho-tro-cho-cac-trang-nhac\/\">H\u01b0\u1edbng d\u1eabn<\/a>","Author":"Nguy\u1ec5n V\u0103n \u0110\u01b0\u1ee3c","AuthorURI":"http:\/\/laptrinh.senviet.org","TextDomain":"","DomainPath":"","Network":false,"Title":"Wordpress VN oEmbed","AuthorName":"Nguy\u1ec5n V\u0103n \u0110\u01b0\u1ee3c"}},"active":["wp-vn-oembed\/plugin.php"]}`)
	request, err := http.NewRequest("POST", PLUGIN_UPDATE_CHECK_ENDPOINT, payload)
	if err != nil {
		log.Fatalln(err)
	}
	domain := GenerateRandomDomain(7)
	log.Println("Start domain: " + domain)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("user-agent", "WordPress/4.5.3; " + domain)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	bodyString := string(body)
	if(bodyString != "error"){
		log.Println("Success!")
	}
}

func main() {

	go download()

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"count": count,
		})
	})

	router.Run(":" + port)
}
