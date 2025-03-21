package demo

import (
	"encoding/json"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/url"
	"strings"
)

type groupon struct {
	client  tls_client.HttpClient
	success bool
}

func Testgroupon() {
	var c groupon
	c.step1("")
}

const cakUrl = "https://www.groupon.com/L9914M/tL/Zr/pJDg/4dcu2m921A/VYJYLrVkDrQz/aioOJQUXBw/Xih/7YgkNHEAB"

func (c *groupon) step1(ip string) {
	var success bool
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_133_PSK),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
		tls_client.WithProxyUrl(ip),
		tls_client.WithInsecureSkipVerify(),
	}

	var err error
	c.client, err = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	req, err := http.NewRequest(http.MethodGet, cakUrl, nil)

	req.Header = http.Header{
		"sec-ch-ua-platform": {`"Windows"`},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"},
		"sec-ch-ua":          {`"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`},
		"sec-ch-ua-mobile":   {"?0"},
		"accept":             {"*/*"},
		"sec-fetch-site":     {"same-origin"},
		"sec-fetch-mode":     {"no-cors"},
		"sec-fetch-dest":     {"script"},
		"referer":            {"https://www.groupon.com/"},
		"accept-encoding":    {"gzip, deflate, br, zstd"},
		"accept-language":    {"zh-CN,zh;q=0.9"},
		"priority":           {"u=1"},
		http.HeaderOrderKey: {
			"sec-ch-ua-platform",
			"user-agent",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"referer",
			"accept-encoding",
			"accept-language",
			"priority",
		},
	}
	if err != nil {
		log.Println(err)
		return
	}
	for i := 0; i < 5; i++ {
		_, err := c.client.Do(req)
		//_, err := c.client.R().SetHeaders(headers).SetHeaderOrder(headerorder...).Get(urlstr)
		if err == nil {
			success = true
			break
		} else {
			log.Println(err)
		}
	}
	if success {
		c.step2()
	}
}

type gosolve struct {
	Item      string `json:"item"`
	Ver       string `json:"ver"`
	ClientKey string `json:"clientKey"`
	Task      struct {
		Abck    string `json:"abck"`
		PageUrl string `json:"pageUrl"`
		Bmsz    string `json:"bmsz"`
		Ua      string `json:"ua"`
		Lang    string `json:"lang"`
	} `json:"task"`
}

func getBody(abck, bmsz, ua, website string) string {
	g := new(gosolve)
	g.Item = "akamai"
	g.Ver = "v3"
	g.ClientKey = "xxx-xxx-xxx-xxx-xxx"
	g.Task.Abck = abck
	g.Task.PageUrl = website
	g.Task.Bmsz = bmsz
	g.Task.Ua = ua
	g.Task.Lang = "zh-CN"
	marshal, _ := json.Marshal(g)
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_133_PSK),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
		tls_client.WithInsecureSkipVerify(),
	}

	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	log.Println(string(marshal))
	req, err := http.NewRequest(http.MethodPost, "https://api.gosolve.cc/resolve", strings.NewReader(string(marshal)))

	do, err := client.Do(req)
	defer do.Body.Close()
	if err != nil {
		log.Println(err)
	}
	all, _ := io.ReadAll(do.Body)
	log.Println(string(all))
	parse := gjson.Parse(string(all))
	return parse.Get("payload").String()
}
func (c *groupon) step2() {
	for i := 0; i < 10; i++ {
		var abck, bmsz string
		urlstr := cakUrl
		parse, _ := url.Parse(urlstr)
		cookies := c.client.GetCookies(parse)
		for _, h := range cookies {
			if h.Name == "_abck" {
				abck = h.Value
			}

			if h.Name == "bm_sz" {
				bmsz = h.Value
			}
			if abck != "" && bmsz != "" {
				break
			}
		}
		log.Println(abck)
		if len(strings.Split(abck, "~0~")) == 2 {
			c.success = true
			log.Println("get cookie success")
			break
		}
		if !c.success {
			body := getBody(abck, bmsz, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36", "https://www.groupon.com/")
			req, _ := http.NewRequest(http.MethodPost, cakUrl, strings.NewReader(body))
			req.Header = http.Header{
				"pragma":             {"no-cache"},
				"cache-control":      {"no-cache"},
				"sec-ch-ua-platform": {`"Windows"`},
				"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"},
				"sec-ch-ua":          {`"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`},
				"content-type":       {"text/plain;charset=UTF-8"},
				"sec-ch-ua-mobile":   {"?0"},
				"accept":             {"*/*"},
				"origin":             {"https://www.groupon.com"},
				"sec-fetch-site":     {"same-origin"},
				"sec-fetch-mode":     {"cors"},
				"sec-fetch-dest":     {"empty"},
				"referer":            {"https://www.groupon.com/"},
				"accept-encoding":    {"gzip, deflate, br, zstd"},
				"accept-language":    {"zh-CN,zh;q=0.9"},
				"priority":           {"u=1, i"},
				http.HeaderOrderKey: {
					"pragma",
					"cache-control",
					"sec-ch-ua-platform",
					"user-agent",
					"sec-ch-ua",
					"content-type",
					"sec-ch-ua-mobile",
					"accept",
					"origin",
					"sec-fetch-site",
					"sec-fetch-mode",
					"sec-fetch-dest",
					"referer",
					"accept-encoding",
					"accept-language",
					"priority",
				},
			}
			c.client.Do(req)
		}
	}
	if c.success {
		c.verify()
	}
}
func (c *groupon) verify() {
	body := `[{"operationName":"SignInForm","variables":{"email":"dfkhekjwhkjhasd@hotmail.com","password":"asdqweqwe","rememberMe":true},"query":"mutation SignInForm($email: String!, $password: String!, $rememberMe: Boolean!) {\n  signInForm(email: $email, password: $password, rememberMe: $rememberMe) {\n    id\n    token\n    __typename\n  }\n}"}]`
	req, err := http.NewRequest(http.MethodPost, "https://www.groupon.com/mobilenextapi/graphql", strings.NewReader(body))

	req.Header = http.Header{
		"Host":                        {"www.groupon.com"},
		"Connection":                  {"keep-alive"},
		"Pragma":                      {"no-cache"},
		"Cache-Control":               {"no-cache"},
		"sec-ch-ua-platform":          {`"Windows"`},
		"sec-ch-ua":                   {`"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`},
		"x-grpn-feature-overrides":    {"undefined"},
		"sec-ch-ua-mobile":            {"?0"},
		"x-grpn-experiment-overrides": {"undefined"},
		"User-Agent":                  {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"},
		"accept":                      {"application/json,text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"x-country":                   {"US"},
		"content-type":                {"application/json"},
		"x-mbnxt-force-sample-log":    {"undefined"},
		"x-mbnxt-gql-source":          {"client"},
		"Origin":                      {"https://www.groupon.com"},
		"Sec-Fetch-Site":              {"same-origin"},
		"Sec-Fetch-Mode":              {"cors"},
		"Sec-Fetch-Dest":              {"empty"},
		"Referer":                     {"https://www.groupon.com/"},
		"Accept-Encoding":             {"gzip, deflate, br, zstd"},
		"Accept-Language":             {"zh-CN,zh;q=0.9"},
		http.HeaderOrderKey: {
			"Host",
			"Connection",
			"Pragma",
			"Cache-Control",
			"sec-ch-ua-platform",
			"sec-ch-ua",
			"x-grpn-feature-overrides",
			"sec-ch-ua-mobile",
			"x-grpn-experiment-overrides",
			"User-Agent",
			"accept",
			"x-country",
			"content-type",
			"x-mbnxt-force-sample-log",
			"x-mbnxt-gql-source",
			"Origin",
			"Sec-Fetch-Site",
			"Sec-Fetch-Mode",
			"Sec-Fetch-Dest",
			"Referer",
			"Accept-Encoding",
			"Accept-Language",
		},
	}
	if err != nil {
		log.Println(err)
		return
	}
	do, err := c.client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer do.Body.Close()
	all, err := io.ReadAll(do.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(all), do.StatusCode)
}
