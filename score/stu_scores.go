package score

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"cqu/decode"
)

const (
	LoginScoreUrl        = "http://202.202.1.176:8080/_data/index_login.aspx"
	ContentType          = "application/x-www-form-urlencoded"
	__VIEWSTATE          = "dDw1OTgzNjYzMjM7dDw7bDxpPDE%2BO2k8Mz47aTw1Pjs%2BO2w8dDxwPGw8VGV4dDs%2BO2w86YeN5bqG5aSn5a2mOz4%2BOzs%2BO3Q8cDxsPFRleHQ7PjtsPFw8c2NyaXB0IHR5cGU9InRleHQvamF2YXNjcmlwdCJcPgpcPCEtLQpmdW5jdGlvbiBvcGVuV2luTG9nKHRoZVVSTCx3LGgpewp2YXIgVGZvcm0scmV0U3RyXDsKZXZhbCgiVGZvcm09J3dpZHRoPSIrdysiLGhlaWdodD0iK2grIixzY3JvbGxiYXJzPW5vLHJlc2l6YWJsZT1ubyciKVw7CnBvcD13aW5kb3cub3Blbih0aGVVUkwsJ3dpbktQVCcsVGZvcm0pXDsgLy9wb3AubW92ZVRvKDAsNzUpXDsKZXZhbCgiVGZvcm09J2RpYWxvZ1dpZHRoOiIrdysicHhcO2RpYWxvZ0hlaWdodDoiK2grInB4XDtzdGF0dXM6bm9cO3Njcm9sbGJhcnM9bm9cO2hlbHA6bm8nIilcOwppZih0eXBlb2YocmV0U3RyKSE9J3VuZGVmaW5lZCcpIGFsZXJ0KHJldFN0cilcOwp9CmZ1bmN0aW9uIHNob3dMYXkoZGl2SWQpewp2YXIgb2JqRGl2ID0gZXZhbChkaXZJZClcOwppZiAob2JqRGl2LnN0eWxlLmRpc3BsYXk9PSJub25lIikKe29iakRpdi5zdHlsZS5kaXNwbGF5PSIiXDt9CmVsc2V7b2JqRGl2LnN0eWxlLmRpc3BsYXk9Im5vbmUiXDt9Cn0KZnVuY3Rpb24gc2VsVHllTmFtZSgpewogIGRvY3VtZW50LmFsbC50eXBlTmFtZS52YWx1ZT1kb2N1bWVudC5hbGwuU2VsX1R5cGUub3B0aW9uc1tkb2N1bWVudC5hbGwuU2VsX1R5cGUuc2VsZWN0ZWRJbmRleF0udGV4dFw7Cn0KZnVuY3Rpb24gd2luZG93Lm9ubG9hZCgpewoJdmFyIHNQQz13aW5kb3cubmF2aWdhdG9yLnVzZXJBZ2VudCt3aW5kb3cubmF2aWdhdG9yLmNwdUNsYXNzK3dpbmRvdy5uYXZpZ2F0b3IuYXBwTWlub3JWZXJzaW9uKycgU046TlVMTCdcOwp0cnl7ZG9jdW1lbnQuYWxsLnBjSW5mby52YWx1ZT1zUENcO31jYXRjaChlcnIpe30KdHJ5e2RvY3VtZW50LmFsbC50eHRfZHNkc2RzZGpramtqYy5mb2N1cygpXDt9Y2F0Y2goZXJyKXt9CnRyeXtkb2N1bWVudC5hbGwudHlwZU5hbWUudmFsdWU9ZG9jdW1lbnQuYWxsLlNlbF9UeXBlLm9wdGlvbnNbZG9jdW1lbnQuYWxsLlNlbF9UeXBlLnNlbGVjdGVkSW5kZXhdLnRleHRcO31jYXRjaChlcnIpe30KfQpmdW5jdGlvbiBvcGVuV2luRGlhbG9nKHVybCxzY3IsdyxoKQp7CnZhciBUZm9ybVw7CmV2YWwoIlRmb3JtPSdkaWFsb2dXaWR0aDoiK3crInB4XDtkaWFsb2dIZWlnaHQ6IitoKyJweFw7c3RhdHVzOiIrc2NyKyJcO3Njcm9sbGJhcnM9bm9cO2hlbHA6bm8nIilcOwp3aW5kb3cuc2hvd01vZGFsRGlhbG9nKHVybCwxLFRmb3JtKVw7Cn0KZnVuY3Rpb24gb3Blbldpbih0aGVVUkwpewp2YXIgVGZvcm0sdyxoXDsKdHJ5ewoJdz13aW5kb3cuc2NyZWVuLndpZHRoLTEwXDsKfWNhdGNoKGUpe30KdHJ5ewpoPXdpbmRvdy5zY3JlZW4uaGVpZ2h0LTMwXDsKfWNhdGNoKGUpe30KdHJ5e2V2YWwoIlRmb3JtPSd3aWR0aD0iK3crIixoZWlnaHQ9IitoKyIsc2Nyb2xsYmFycz1ubyxzdGF0dXM9bm8scmVzaXphYmxlPXllcyciKVw7CnBvcD1wYXJlbnQud2luZG93Lm9wZW4odGhlVVJMLCcnLFRmb3JtKVw7CnBvcC5tb3ZlVG8oMCwwKVw7CnBhcmVudC5vcGVuZXI9bnVsbFw7CnBhcmVudC5jbG9zZSgpXDt9Y2F0Y2goZSl7fQp9CmZ1bmN0aW9uIGNoYW5nZVZhbGlkYXRlQ29kZShPYmopewp2YXIgZHQgPSBuZXcgRGF0ZSgpXDsKT2JqLnNyYz0iLi4vc3lzL1ZhbGlkYXRlQ29kZS5hc3B4P3Q9IitkdC5nZXRNaWxsaXNlY29uZHMoKVw7Cn0KXFwtLVw%2BClw8L3NjcmlwdFw%2BOz4%2BOzs%2BO3Q8O2w8aTwxPjs%2BO2w8dDw7bDxpPDA%2BOz47bDx0PHA8bDxUZXh0Oz47bDxcPG9wdGlvbiB2YWx1ZT0nU1RVJyB1c3JJRD0n5a2m5Y%2B3J1w%2B5a2m55SfXDwvb3B0aW9uXD4KXDxvcHRpb24gdmFsdWU9J1RFQScgdXNySUQ9J%2BW4kOWPtydcPuaVmeW4iFw8L29wdGlvblw%2BClw8b3B0aW9uIHZhbHVlPSdTWVMnIHVzcklEPSfluJDlj7cnXD7nrqHnkIbkurrlkZhcPC9vcHRpb25cPgpcPG9wdGlvbiB2YWx1ZT0nQURNJyB1c3JJRD0n5biQ5Y%2B3J1w%2B6Zeo5oi357u05oqk5ZGYXDwvb3B0aW9uXD4KOz4%2BOzs%2BOz4%2BOz4%2BOz4%2BOz7p2B9lkx%2BYq%2Fjf62i%2BiqicmZx%2Fxg%3D%3D"
	__VIEWSTATEGENERATOR = "CAA0A5A7"
	Sel_Type             = "STU"
	FetchScoreUrl        = "http://202.202.1.176:8080/xscj/Stu_MyScore_rpt.aspx"
)

type ScoreFetchSection struct {
	sel_xn   string //学年,当前学年的前一个年份，比如是2016-2017学年，那就是2016
	sel_xq   string //学期 ,0 代表第一学期，1代表第二学期
	SJ       string
	SelXNXQ  string
	zfx_flag string
	zxf      string
}

type StuInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Semester string `json:"semester"`
}

type Score struct {
	Lesson          string `json:"lesson"`            //课程
	Credit          string `json:"credit"`            //学分
	LessonType      string `json:"lesson_type"`       //类别
	Evaluation      string `json:"evaluation"`        //考核方式
	CourseCharacter string `json:"course_character"`  //修读性质
	Score           string `json:"score"`             //成绩
	MinorRepairMark string `json:"minor_repair_mark"` //辅修标识
}

type Responce struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	StuInfo *StuInfo `json:"stu_info"`
	Scores  []Score  `json:"scores"`
}

func StuScore(id, pass, year, semester string) (string, int, error) {
	pass = CquMd5Encrypted(id, pass)
	htmlResult, err := scoresHttpRequst(id, pass, year, semester)
	if err != nil {
		log.Printf("StuScore : httpRequst error=[%v]", err)
		return "", 0, err
	}
	result, num, err := analysisHtml(htmlResult)
	if err != nil {
		log.Printf("StuScore : analysis error=[%v]", err)
		return "", 0, nil
	}
	return result, num, nil
}

func StuCookie(id, pass string) (string, error) {

	client := http.Client{}
	bodyBuf := new(bytes.Buffer)
	bodyBuf.WriteString("__VIEWSTATE=" + __VIEWSTATE)
	bodyBuf.WriteString("&__VIEWSTATEGENERATOR=" + __VIEWSTATEGENERATOR)
	bodyBuf.WriteString("&Sel_Type=" + Sel_Type)
	bodyBuf.WriteString("&txt_dsdsdsdjkjkjc=" + id)
	bodyBuf.WriteString("&txt_dsdfdfgfouyy=" + "hw465520")
	bodyBuf.WriteString("&txt_ysdsdsdskgf=")
	bodyBuf.WriteString("&pcInfo=")
	bodyBuf.WriteString("&typeName=")
	bodyBuf.WriteString("&aerererdsdxcxdfgfg=")
	bodyBuf.WriteString("&efdfdfuuyyuuckjg=" + pass)

	req, err := http.NewRequest("POST", LoginScoreUrl, bodyBuf)
	if err != nil {
		log.Printf("StuCookie : newRequest error=[%v]", err)
		return "", err
	}

	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("StuCookie : client do error=[%v]", err)
		return "", err
	}

	defer resp.Body.Close()

	bufCookie := new(bytes.Buffer)
	var i, lenght = 0, len(resp.Cookies())
	for _, cookie := range resp.Cookies() { //todo
		strCookie := cookie.String()
		strs := strings.Split(strCookie, ";")
		strCookie = strs[0]
		if i != lenght-1 {
			bufCookie.WriteString(strCookie + "; ")
		} else {
			bufCookie.WriteString(strCookie)
		}
		i++
	}
	return bufCookie.String(), nil
}

func CreateErrorResponce(status, message string) (string, error) {
	resp := Responce{}
	resp.Scores = make([]Score, 0, 1)
	resp.StuInfo = new(StuInfo)
	resp.Status = status
	resp.Message = message
	result, err := json.Marshal(resp)
	if err != nil {
		log.Printf("score : CreateErrorResponce : json error=[%v]", err)
		return "", err
	}
	return string(result), nil
}

func scoresHttpRequst(id, pass, year, semester string) (string, error) {
	client := &http.Client{}

	bufBody := new(bytes.Buffer)

	cookie, err := StuCookie(id, pass)
	if err != nil {
		log.Printf("scoresHttpRequst : cookie error=[%v]", err)
		return "", err
	}

	scoreSection := ScoreFetchSection{
		sel_xn:   year,     //学年
		sel_xq:   semester, //学期，0表示第一学期，1表示第二学期
		SJ:       "0",
		SelXNXQ:  "2",
		zfx_flag: "0",
		zxf:      "0",
	}
	bufBody.WriteString(scoreSection.String())

	req, err := http.NewRequest("POST", FetchScoreUrl, bufBody)
	if err != nil {
		log.Printf("scoresHttpRequst : new post requst error=[%v]", err)
		return "", err
	}
	req.Header.Set("Content-Type", ContentType)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("scoresHttpRequst : client do error=[%v]", err)
		return "", err
	}

	defer resp.Body.Close()

	bufBodyResp := new(bytes.Buffer)
	_, err = io.Copy(bufBodyResp, resp.Body)
	if err != nil {
		log.Printf("scoresHttpRequst : resp body copy error=[%v]", err)
		return "", err
	}

	const Duration = 1 * time.Second
	end := time.Now().Add(Duration)
	try := 0
	var result string = ""
	for try = 0; time.Now().Before(end); try++ {
		err = nil
		result, err = decode.Decode(bufBodyResp.String())
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Printf("scoresHttpRequst : copy error,duraion 1s invoke %d time.error=[%v],repsBody=%s", try, err, bufBodyResp.String())
		return "", err
	}
	return result, nil
}

func analysisHtml(htmlResult string) (string, int, error) {
	strReader := strings.NewReader(htmlResult)
	rootNode, err := html.Parse(strReader)
	if err != nil {
		log.Printf("analysisHtml : parse error=[%v]", err)
		return "", 0, err
	}
	stuInfoScore := new(Responce)
	stuInfoScore.StuInfo = new(StuInfo)
	stuInfoScore.Scores = make([]Score, 0, 30)
	num := visitScores(stuInfoScore, rootNode, 0)
	if len(stuInfoScore.Scores) == 0 {
		stuInfoScore.Status = "400"
		stuInfoScore.Message = "id or pass error."
	} else {
		stuInfoScore.Status = "200"
		stuInfoScore.Message = "success."
	}
	result, err := json.Marshal(stuInfoScore)
	if err != nil {
		log.Printf("score : analysisHtml : json error=[%v]", err)
		return "", 0, err
	}
	return string(result), num - 3, nil
}

func visitScores(resp *Responce, n *html.Node, num int) int {
	if n.Type == html.ElementNode && n.Data == "tr" {
		if num == 0 {
			infoValNode := n.FirstChild.FirstChild
			if infoValNode != nil {
				infoVal := infoValNode.Data
				setStuInfo(resp, infoVal)
			}
			num++
			return num //num = 1
		}

		if num == 1 {
			semesterValNode := n.FirstChild.FirstChild
			if semesterValNode != nil {
				resp.StuInfo.Semester = semesterValNode.Data
			}
			num++
			return num //num = 2
		}
		if num == 2 { //标题
			num++
			return num //num = 3
		}
		var tdNum = 0
		var trVal = Score{}
		for td := n.FirstChild; td != nil; td = td.NextSibling {
			for val := td.FirstChild; val != nil; val = val.NextSibling {
				if val.Type == html.TextNode {
					switch tdNum {
					case 1:
						trVal.Lesson = val.Data
					case 2:
						trVal.Credit = val.Data
					case 3:
						trVal.LessonType = val.Data
					case 4:
						trVal.Evaluation = val.Data
					case 5:
						trVal.CourseCharacter = val.Data
					case 6:
						trVal.Score = val.Data
					case 7:
						trVal.MinorRepairMark = val.Data
					}
				}
			}
			tdNum++
		}
		resp.Scores = append(resp.Scores, trVal)
		num++
		return num //num > 2
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		num = visitScores(resp, c, num)
	}
	return num
}

func setStuInfo(resp *Responce, infoVal string) {
	strs := strings.Split(infoVal, "   ")
	if len(strs) == 2 {
		id := strings.TrimSpace(strings.Split(strs[0], "：")[1])
		name := strings.TrimSpace(strings.Split(strs[1], "：")[1])
		resp.StuInfo.Id = id
		resp.StuInfo.Name = name
	}
}

func (section *ScoreFetchSection) String() string {
	return "sel_xn=" + section.sel_xn + "&sel_xq=" + section.sel_xq + "&SJ=" +
		section.SJ + "&SelXNXQ=" + section.SelXNXQ + "&zfx_flag=" +
		section.zfx_flag + "&zxf=" + section.zxf
}
