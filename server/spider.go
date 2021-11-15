package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"ssq-spider/dao"
	"ssq-spider/logger"
	"ssq-spider/model"
	"strconv"
	"strings"
	"time"
)

type SSQPageInfo struct {
	IssueNum string
	PageUrl  string
}

const (
	StartPage = "https://kaijiang.500.com/ssq.shtml"
)

var PageRe = regexp.MustCompile(`https://kaijiang.500.com/shtml/ssq/(\d+).shtml`)
var RedBallRe = regexp.MustCompile(`<li class="ball_red">(\d+)</li>`)
var BlueBallRe = regexp.MustCompile(`<li class="ball_blue">(\d+)</li>`)
var BallSortRe = regexp.MustCompile(`<td>\s*(\d+)\s*(\d+)\s*(\d+)\s*(\d+)\s*(\d+)\s*(\d+)</td>\s*</tr>\s*</table>`)

func RequestPage(url string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Logger.Error("new request error: ", err)
		return ""
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			logger.Logger.Info("Redirect: ", req)
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Error("do request error: ", err)
		return ""
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Error("read body error: ", err)
		return ""
	}

	return string(content)
}

func GetAllPage(maxIssueNum int) *[]SSQPageInfo {
	var maxIssueNumString string
	var pageInfoList []SSQPageInfo

	startPageContent := RequestPage(StartPage)

	pageUrlList := PageRe.FindAllStringSubmatch(startPageContent, -1)
	if len(pageUrlList) > 0 {
		maxIssueNumString = pageUrlList[0][1]
	}

	for _, pageUrl := range pageUrlList {
		issueNum, err := strconv.Atoi(pageUrl[1])
		if err != nil {
			logger.Logger.Error("strconv issueNum error: ", err)
			return nil
		}

		if issueNum > maxIssueNum {
			pageInfoList = append(pageInfoList, SSQPageInfo{
				IssueNum: pageUrl[1],
				PageUrl:  pageUrl[0],
			})
		}
	}

	if len(maxIssueNumString) > 0 {
		if err := dao.UpdateOneSysConfig("1", "maxIssueNum", maxIssueNumString); err != nil {
			logger.Logger.Error("update maxIssueNum error: ", err)
			return nil
		}
	}

	return &pageInfoList
}

func ParseOnePage(page *SSQPageInfo) (string, string, error) {
	pageContent := RequestPage(page.PageUrl)

	var redBall []string
	redBallMatches := RedBallRe.FindAllStringSubmatch(pageContent, -1)
	if len(redBallMatches) == 0 {
		return "", "", errors.New("no red ball match")
	}
	for _, redNum := range redBallMatches {
		redBall = append(redBall, redNum[1])
	}

	blueBallMatches := BlueBallRe.FindAllStringSubmatch(pageContent, -1)
	if len(blueBallMatches) == 0 {
		return "", "", errors.New("no blue ball match")
	}

	ballString := page.IssueNum + "+" + strings.Join(redBall, ",") + "+" + blueBallMatches[0][1]

	ballSortList := BallSortRe.FindAllStringSubmatch(pageContent, -1)
	if len(ballSortList) != 1 {
		return "", "", errors.New("no ball sort match")
	}

	ballSortString := ""
	for _, v := range ballSortList[0][1:] {
		ballSortString += v
		ballSortString += ","
	}

	return ballString, ballSortString[:len(ballSortString)-1], nil
}

func SSQSpider() {
	for {
		time.Sleep(time.Minute)

		maxIssueNumString, err := dao.GetOneSysConfig("1", "maxIssueNum")
		if err != nil {
			logger.Logger.Error("get maxIssueNum error: ", err)
			continue
		}

		maxIssueNum, err := strconv.Atoi(maxIssueNumString)
		if err != nil {
			logger.Logger.Error("atoi maxIssueNumString error: ", err)
			continue
		}

		var ssqOpenNumList []model.SSQOpenNumber

		pageInfoList := GetAllPage(maxIssueNum)
		for index, page := range *pageInfoList {
			fmt.Printf("get page %d, url: %s", index, page.PageUrl)
			ballString, ballSortString, err := ParseOnePage(&page)
			if err != nil {
				logger.Logger.Error("parse one page error: ", err, " url: ", page.PageUrl)
				continue
			}

			if len(ballString) == 0 || len(ballSortString) == 0 {
				continue
			}

			openNums := strings.Split(ballString, "+")
			if len(openNums) != 3 {
				continue
			}

			ssqOpenNumList = append(ssqOpenNumList, model.SSQOpenNumber{
				OpenNo:   openNums[0],
				RedNum:   openNums[1],
				BlueNum:  openNums[2],
				BallSort: ballSortString,
			})

			time.Sleep(time.Millisecond * 10)
		}

		// save to db
		if err := dao.BulkSaveSSQOpenNumbers(&ssqOpenNumList); err != nil {
			logger.Logger.Error("BulkSaveSSQOpenNumbers error: ", err)
			continue
		}

		logger.Logger.Debugf("save ssq open number to db, count: %d", len(ssqOpenNumList))
	}
}
