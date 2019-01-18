package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type config struct {
	username     string
	password     string
	semesterYear string
}

type course struct {
	courseType string
	courseName string
}

type reqPayload struct {
	pageNo   string
	pageSize string
	param    Param
}

type Param struct {
	collectionStatus     string
	hiddenConflictStatus string
	hiddenSelectedStatus string
	selectedCate         string
	selectedType         string
	semesterYear         string
}

func main() {

	cookie := getCaptchaPicAndCookie("https://cas.sysu.edu.cn/cas/captcha.jsp")
	// fmt.Println(cookie)
	token := getLoginHTML("https://cas.sysu.edu.cn/cas/login?service=https%3A%2F%2Fuems.sysu.edu.cn%2Fjwxt%2Fapi%2Fsso%2Fcas%2Flogin%3Fpattern%3Dstudent-login")
	login("https://cas.sysu.edu.cn/cas/login?service=https%3A%2F%2Fuems.sysu.edu.cn%2Fjwxt%2Fapi%2Fsso%2Fcas%2Flogin%3Fpattern%3Dstudent-login", token, cookie)
	// getCourseList("https://uems.sysu.edu.cn/jwxt/choose-course-front-server/classCourseInfo/selectCourseInfo?_t=1547712166")

}

func getCaptchaPicAndCookie(url string) string {
	captchaPic, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer captchaPic.Body.Close()

	//get cookie
	cookie := captchaPic.Header.Get("Set-Cookie")
	endIndex := strings.Index(cookie, ";")

	//save captcha picture
	pix, err := ioutil.ReadAll(captchaPic.Body)
	if err != nil {
		fmt.Println(err)
	}
	filePath, err := os.Create("./captcha.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer filePath.Close()
	_, err = io.Copy(filePath, bytes.NewReader(pix))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Captcha has been saved")
	}

	return cookie[:endIndex]
}

func login(loginURL string, token string, cookie string) {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter captcha: ")
	// captchaCode, _ := reader.ReadString('\n')
	captchaCode := "test"

	client := &http.Client{}

	// url.Values{"username": {"huangyt39"}, "password": {"chAnrd7ler"}, "captcha": {captchaCode}}
	req, err := http.NewRequest("POST", loginURL, nil)
	if err != nil {
		fmt.Println(err)
	}

	req.ParseForm()
	req.Form.Add("username", "huangyt39")
	req.Form.Add("password", "chAnrd7ler")
	req.Form.Add("captcha", captchaCode)
	req.Form.Add("execution", token)
	req.Form.Add("_eventId", "submit")

	req.Header.Set("Accept", "*/*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36(KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Set("origin", "https://cas.sysu.edu.cn")
	req.Header.Set("host", "cas.sysu.edu.cn")
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
	fmt.Println(resp.StatusCode)
}

func getCourseList(courseListURL string) {

	s := reqPayload{pageNo: "1", pageSize: "1000", param: Param{"0", "0", "0", "11", "1", "2018-2"}}
	b, _ := json.Marshal(s)
	fmt.Println("Here is the rePayload:")
	fmt.Println(b)

	resp, err := http.Post(courseListURL, "application/json;charset=UTF-8", strings.NewReader("heel="+string(b)))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func getLoginHTML(loginURL string) string {

	client := &http.Client{}

	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36(KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("host", "cas.sysu.edu.cn")
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	// req.Header.Set("Cookie", "safedog-flow-item=B66D062664F275F8EB5648510EB65F45")
	req.Header.Set("Referer", "https://uems.sysu.edu.cn/jwxt/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(b))
	HTMLStr := string(b)
	return getTokenFromHTML(HTMLStr)
}

func getTokenFromHTML(HTMLStr string) string {
	startIndex := strings.LastIndex(HTMLStr, "name=\"execution\" value=\"")
	endIndex := strings.LastIndex(HTMLStr, "\" /><input type=\"hidden\" name=\"_eventId\"")
	resStr := HTMLStr[startIndex+24 : endIndex]
	// fmt.Println(resStr)
	return resStr
}
