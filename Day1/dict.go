package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}
type HuoShanRequest struct {
	Text     string `json:"text"`
	Language string `json:"language"`
}
type HuoShanResponse struct {
	Words []struct {
		Source  int    `json:"source"`
		Text    string `json:"text"`
		PosList []struct {
			Type      int `json:"type"`
			Phonetics []struct {
				Type int    `json:"type"`
				Text string `json:"text"`
			} `json:"phonetics"`
			Explanations []struct {
				Text     string `json:"text"`
				Examples []struct {
					Type      int `json:"type"`
					Sentences []struct {
						Text      string `json:"text"`
						TransText string `json:"trans_text"`
					} `json:"sentences"`
				} `json:"examples"`
				Synonyms []interface{} `json:"synonyms"`
			} `json:"explanations"`
			Relevancys []interface{} `json:"relevancys"`
		} `json:"pos_list"`
	} `json:"words"`
	Phrases  []interface{} `json:"phrases"`
	BaseResp struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"base_resp"`
}
type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

func queryFromCaiyun(word string) {
	client := &http.Client{}
	//序列化请求
	request := DictRequest{TransType: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	//创建请求
	//创建一个data流
	var data = bytes.NewBuffer(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	//设置请求头
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("os-version", "")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
	req.Header.Set("app-name", "xy")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("device-id", "")
	req.Header.Set("os-type", "web")
	req.Header.Set("X-Authorization", "token:qgemv4jr1y38jyq6vhvi")
	req.Header.Set("Origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "_ym_uid=16456948721020430059; _ym_d=1645694872")
	//发起请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//读响应 延迟关闭 Body的本质其实就是一个IO流
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//如果响应结果不是200
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode", resp.StatusCode, "body", string(bodyText))
	}
	//把json反序列化到结构提中
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("彩云翻译结果")
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
	wg.Done()
}

func queryFromHuoShan(word string) {
	client := &http.Client{}
	request := HuoShanRequest{Text: word, Language: "en"}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://translate.volcengine.com/web/dict/match/v1/?msToken=&X-Bogus=DFSzKwVLQDcCwf5fSvhUPKXAIQ5b&_signature=_02B4Z6wo00001gCAgoAAAIDDY4pC6Rk10gYAkIYAAOKTy9LdtcqqV9KwF5W1KFWZnsrYjcksZs2wQF2FnEACWZ6gkkYL0O8BgN2COC6pyM63mhyJP.U7-DEMVCN2WQPOUcIjZR8EkQUNUmiv35", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "translate.volcengine.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("cookie", "x-jupiter-uuid=16550192397566480; s_v_web_id=verify_l4azq7le_2HLzMBe3_TaN5_4463_Athh_4x0guqJJo12f; _tea_utm_cache_2018=undefined; ttcid=29073d5e671d463190394cb24212de9c14; i18next=translate; tt_scid=xHY8rC2iqTrgXS2iJ6wZMuZAFP-TlOLfzlgU7gK9W8yA6lpbpxBeUHCjwWtUjTUmc3e4")
	req.Header.Set("origin", "https://translate.volcengine.com")
	req.Header.Set("referer", "https://translate.volcengine.com/translate?category=&home_language=zh&source_language=detect&target_language=zh&text=hello")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="102", "Google Chrome";v="102"`)
	req.Header.Set("sec-ch-ua-mobile", "?1")
	req.Header.Set("sec-ch-ua-platform", `"Android"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Mobile Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Bad StatusCode:", resp.StatusCode, "body", string(bodyText))
	}
	var response HuoShanResponse
	err = json.Unmarshal(bodyText, &response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("火山翻译结果")
	fmt.Println(word, response.Words[0].PosList[0].Phonetics[0].Text)
	for _, item := range response.Words[0].PosList[0].Explanations {
		fmt.Println(item.Text)
		for _, str := range item.Examples {
			fmt.Println(str.Sentences)
		}
	}
	wg.Done()
}
func main() {
	//第一个参数是go文件名 第二个参数是单词
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
	example: simpleDict hello
			`)
		os.Exit(1)
	}
	//	//拿到第二个参数 进行query函数的调用
	word := os.Args[1]
	wg.Add(2)
	go queryFromCaiyun(word)
	go queryFromHuoShan(word)
	wg.Wait()
}
