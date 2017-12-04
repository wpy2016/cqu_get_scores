package rank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
	"cqu/decode"
)

const (
	LOGIN_URL    = "http://oldjw.cqu.edu.cn:8088/login.asp"
	NEW_SCORES   = "http://oldjw.cqu.edu.cn:8088/score/sel_score/new_score_sel.asp"
	ALL_SCORES   = "http://oldjw.cqu.edu.cn:8088/score/sel_score/sum_score_sel.asp"
	CONTENT_TYPE = "application/x-www-form-urlencoded"
)

type Lession struct {
	Code        string `json:"code"`     //课程编码
	Name        string `json:"name"`     //课程名称
	Scores      string `json:"scores"`   //成绩
	Credit      string `json:"credit"`   //学分
	Elective    string `json:"elective"` //选修
	LessionType string `json:"type"`     //类别
	Teacher     string `json:"teacher"`  //教师
	Exam        string `json:"exam"`     //考别
	Remark      string `json:"remark"`   //备注
	Time        string `json:"time"`     //时间
}

type StuInfo struct {
	Userid   string `json:"userid"`   //学号
	username string `json:"username"` //姓名
	Major    string `json:"major"`    //专业
	Gpa      string `json:"gpa"`      //GPA
}

type Responce struct {
	Status      string    `json:"status"`      //状态码
	Message     string    `json:"message"`     //状态信息
	Stuinfo     *StuInfo  `json:"stuinfo"`     // 用户信息
	Stulessions []Lession `json:"stulessions"` //课程信息
}

func StuCookie(username, password string) (string, error) {
	client := &http.Client{}
	post := new(bytes.Buffer)
	post.WriteString("username=" + username)
	post.WriteString("&password=" + password)
	post.WriteString("&select1=" + "#")
	req, err := http.NewRequest("POST", LOGIN_URL, post)
	if err != nil {
		log.Printf("rank : httpPost : http request error=[%v]", err)
		return "", err
	}
	req.Header.Set("Content-Type", CONTENT_TYPE)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-cn")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("rank : httpPost : client do error=[%v]", err)
		return "", err
	}
	defer resp.Body.Close()
	cookies := new(bytes.Buffer)
	cookiesLenght := len(resp.Cookies()) - 1
	for pos, cookie := range resp.Cookies() {
		if pos != cookiesLenght {
			cookies.WriteString(cookie.Name + "=" + cookie.Value + ";")
		} else {
			cookies.WriteString(cookie.Name + "=" + cookie.Value)
		}
	}
	return cookies.String(), nil
}

/**
0 表示最新成绩
1 表示全部成绩
*/
func StuScores(username, password string, t int) (string, error) {
	var result string
	var err error
	if t == 0 {
		result, err = scoresHttpRequst(username, password, NEW_SCORES)
		if err != nil {
			log.Printf("rank : StuScores error=[%v]", err)
			return "", err
		}
		result, err = analysisHtml(result, visitScores)
		return result, nil
	}
	if t == 1 {
		result, err = scoresHttpRequst(username, password, ALL_SCORES)
		if err != nil {
			log.Printf("rank : StuScores error=[%v]", err)
			return "", err
		}
		result, err = analysisHtml(result, visitScores)
		return result, nil
	}
	return "", fmt.Errorf("rank : StuScores : out of requst type ,0 and 1")
}

func CreateResponce(status, message string) (string, error) {
	resp := new(Responce)
	resp.Stuinfo = new(StuInfo)
	resp.Stulessions = make([]Lession, 0, 1)
	resp.Status = status
	resp.Message = message
	result, err := json.Marshal(resp)
	if err != nil {
		log.Printf("rank : CreateResponce : json error=[%v]", err)
		return "", err
	}
	return string(result), nil
}

func scoresHttpRequst(id, pass, url string) (string, error) {
	cookis, err := StuCookie(id, pass)
	if err != nil {
		log.Printf("rank : HttpNewScores : cookie error=[%v]", err)
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("rank : HttpNewScores : http request error=[%v]", err)
		return "", err
	}
	req.Header.Set("Cookie", cookis)
	req.Header.Set("Accept-Language", "zh-CN")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("rank : HttpNewScores : client do error=[%v]", err)
		return "", err
	}
	defer resp.Body.Close()
	bufBody := new(bytes.Buffer)
	_, err = io.Copy(bufBody, resp.Body)
	if err != nil {
		log.Printf("rank : HttpNewScores : copy error=[%v]", err)
		return "", err
	}
	bodyUtf, err := decode.Decode(bufBody.String())
	if err != nil {
		log.Printf("rank : HttpNewScores : decode error=[%v]", err)
		return "", err
	}
	return bodyUtf, nil
}

func analysisHtml(bodyUtf string, operate func(*Responce, *html.Node) bool) (string, error) {
	bodyReady := strings.NewReader(bodyUtf)
	rootNode, err := html.Parse(bodyReady)
	if err != nil {
		log.Printf("rank : analysisHtml : html parse error=[%v]", err)
		return "", err
	}
	resp := new(Responce)
	resp.Stuinfo = new(StuInfo)
	resp.Stulessions = make([]Lession, 0, 1)
	operate(resp, rootNode)
	if len(resp.Stulessions) == 0 {
		resp.Message = "username or password error."
		resp.Status = "404"
	} else {
		resp.Message = "success."
		resp.Status = "200"
	}
	result, err := json.Marshal(resp)
	if err != nil {
		log.Printf("rank : analysisHtml : json marshal error=[%v]", err)
		return "", err
	}
	return string(result), nil
}

func visitScores(resp *Responce, n *html.Node) bool {
	if n.Type == html.ElementNode && n.Data == "center" {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && child.Data == "table" { // 找到table标签
				for bodyChild := child.FirstChild; bodyChild != nil; bodyChild = bodyChild.NextSibling {
					if bodyChild.Type == html.ElementNode && bodyChild.Data == "tbody" { //找到tbody标签
						for trChild := bodyChild.FirstChild; trChild != nil; trChild = trChild.NextSibling {
							if trChild.Type == html.ElementNode && trChild.Data == "tr" { // 找到tr
								for tdChild := trChild.FirstChild; tdChild != nil; tdChild = tdChild.NextSibling {
									if tdChild.Type == html.ElementNode && tdChild.Data == "td" { //找到td
										for valNode := tdChild.FirstChild; valNode != nil; valNode = valNode.NextSibling {
											if valNode.Type == html.ElementNode && valNode.Data == "p" {
												setStuInfo(resp, valNode)
											}
											if valNode.Type == html.ElementNode && valNode.Data == "table" {
												setLessions(resp, valNode)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return true
	}
	isStop := false
	for c := n.FirstChild; c != nil && !isStop; c = c.NextSibling {
		isStop = visitScores(resp, c)
	}
	return isStop
}

func setStuInfo(resp *Responce, n *html.Node) {
	pos := 0
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			switch pos {
			case 0:
				resp.Stuinfo.Userid = strings.TrimSpace(child.Data)
			case 1:
				resp.Stuinfo.username = strings.TrimSpace(child.Data)
			case 2:
				resp.Stuinfo.Major = strings.TrimSpace(child.Data)
			case 3:
				resp.Stuinfo.Gpa = strings.TrimSpace(child.Data)
			}
			pos++
		}
	}
}

func setLessions(resp *Responce, n *html.Node) {
	for tbody := n.FirstChild; tbody != nil; tbody = tbody.NextSibling {
		if tbody.Type == html.ElementNode && tbody.Data == "tbody" {
			count := -1
			for trchild := tbody.FirstChild; trchild != nil; trchild = trchild.NextSibling { //行节点
				if trchild.Type == html.ElementNode && trchild.Data == "tr" {
					count++
					if count == 0 {
						continue
					}
					lession := Lession{}
					tdnum := 0
					for trChild := trchild.FirstChild; trChild != nil; trChild = trChild.NextSibling { // 列节点遍历
						if trChild.Type == html.ElementNode && trChild.Data == "td" {
							for valNode := trChild.FirstChild; valNode != nil; valNode = valNode.NextSibling { //单个列节点
								if valNode.Type == html.TextNode {
									switch tdnum {
									case 0:
										break
									case 1:
										lession.Code = strings.TrimSpace(valNode.Data)
									case 2:
										lession.Name = strings.TrimSpace(valNode.Data)
									case 3:
										lession.Scores = strings.TrimSpace(valNode.Data)
									case 4:
										lession.Credit = strings.TrimSpace(valNode.Data)
									case 5:
										lession.Elective = strings.TrimSpace(valNode.Data)
									case 6:
										lession.LessionType = strings.TrimSpace(valNode.Data)
									case 7:
										lession.Teacher = strings.TrimSpace(valNode.Data)
									case 8:
										lession.Exam = strings.TrimSpace(valNode.Data)
									case 9:
										lession.Remark = strings.TrimSpace(valNode.Data)
									case 10:
										lession.Time = strings.TrimSpace(valNode.Data)
									}
								}
							}
							tdnum++
						}
					}
					resp.Stulessions = append(resp.Stulessions, lession)
				}
			}
		}
	}
}
