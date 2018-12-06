package renren

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli"
)

var urlChannel = make(chan string, 200)

// GetUSMovie 抓取指定美剧资源
func GetUSMovie(c *cli.Context) {
	go spider("http://www.msj1.com/archives/6121.html")

	file, err := os.Create("./link.text")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

	for url := range urlChannel {
		fmt.Println("下载链接 => ", url)
		_, err := file.WriteString(url + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}

var userAgent = []string{
	"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
	"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
	"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
	"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
	"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
	"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
	"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandomUserAgent() string {
	return userAgent[r.Intn(len(userAgent))]
}

//以Must前缀的方法或函数都是必须保证一定能执行成功的,否则将引发一次panic
var aTagRegExp = regexp.MustCompile(`<a[^>]+[(href)|(HREF)]\s*\t*\n*=\s*\t*\n*[(".+")|('.+')][^>]*>[^<]*</a>`)

func getHref(aTag string) (href, content string) {
	inputReader := strings.NewReader(aTag)
	decoder := xml.NewDecoder(inputReader)

	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		// 处理元素标签开始部分
		case xml.StartElement:
			for _, attr := range token.Attr {
				attrName := attr.Name.Local
				attrValue := attr.Value
				if strings.EqualFold(attrName, "href") {
					href = attrValue
				}
			}
		// 处理元素标签结束部分
		case xml.EndElement:
		// 处理元素字符数据部分
		case xml.CharData:
			content = string([]byte(token))
		default:
			href = ""
			content = ""
		}
	}
	return href, content
}

func spider(url string) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", getRandomUserAgent())
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(res.StatusCode)
	if res.StatusCode == 200 {

		body := res.Body
		defer body.Close()

		bodyByte, _ := ioutil.ReadAll(body)
		resStr := string(bodyByte)
		aTag := aTagRegExp.FindAllString(resStr, -1)

		for _, a := range aTag {
			href, _ := getHref(a)
			if strings.Contains(href, "ed2k://") {
				urlChannel <- href
			}
		}
	}
}
