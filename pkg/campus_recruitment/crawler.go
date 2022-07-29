package campus_recruitment

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
	"xsky_crawler/pkg/consts"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
	"k8s.io/klog"
)

type PositionInfo struct {
	Name        string `json:"name"`
	Base        string `json:"base"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

const (
	// Default total number of pages
	defaultCurrent = 1
	defaultLimit   = 100
	// HTML Tag Rules
	scratchSelector = "#bd > section > section > main > div > div > div.content__bb7170 > div.rightBlock.rightBlock__bb7170 > div.borderContainer__bb7170 > div.listItems__bb7170"
	UserAgent       = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
	param           = `document.querySelector("body")`
)

type SelectorData struct {
	Type     int
	Selector string
}

const (
	PositionTypeName = iota
	PositionTypeBaseAndType
	PositionTypeDescription
)

var selectorData = []*SelectorData{
	{PositionTypeName, ".positionItem-title-text"},
	{PositionTypeBaseAndType, ".subTitle__bb7170,.positionItem-subTitle"},
	{PositionTypeDescription, ".jobDesc__bb7170,.positionItem-jobDesc"}}

func Crawler(url string,current,limit int) error {
	start := time.Now()
	klog.Infof("start campus recruitment: %v \n",start.Format(consts.DefaultTimeFormat))
	requestURL := fmt.Sprintf("%s/?current=%d&limit=%d", url, defaultCurrent, defaultLimit)
	content, err := getHTMLContent(requestURL, scratchSelector, param)
	if err != nil {
		return err
	}
	res, err := collectPositionInfo(content, selectorData...)
	err = writeFile(res)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	end := time.Now()
	klog.Infof("end campus recruitment: %v, execute crawler duration(unit milliseconds)： %d\n",end.Format(consts.DefaultTimeFormat),end.Sub(start).Milliseconds())
	return err
}

// getHTMLContent Get the HTML content by URL、selector、sel
func getHTMLContent(url string, selector string, sel interface{}) (string, error) {
	options := []chromedp.ExecAllocatorOption{
		// disable image
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(UserAgent),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	c, _ := chromedp.NewExecAllocator(context.Background(), options...)
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	_ = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 30*time.Second)
	defer cancel()
	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector),
		chromedp.OuterHTML(sel, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		panic(err)
		return "", err
	}
	return htmlContent, nil
}

func collectPositionInfo(htmlContent string, selectorData ...*SelectorData) ([]*PositionInfo, error) {
	// create dom reader
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, err
	}
	var res []*PositionInfo
	response := make(chan *PositionInfo, 20)
	wgResponse := &sync.WaitGroup{}
	wg := &sync.WaitGroup{}
	go func() {
		wgResponse.Add(1)
		for rc := range response {
			res = append(res, rc)
		}
		wgResponse.Done()
	}()
	selection := dom.Find(joinSelectorData(selectorData...))
	for index := 0; index < len(selection.Nodes); index += 3 {
		wg.Add(1)
		go func(index int) {
			tempPositionInfo := &PositionInfo{}
			tempPositionInfo.Name = (&goquery.Selection{Nodes: []*html.Node{selection.Nodes[index]}}).Text()
			runes := []rune((&goquery.Selection{Nodes: []*html.Node{selection.Nodes[index+1]}}).Text())
			tempPositionInfo.Base, tempPositionInfo.Type = string(runes[:2]), string(runes[2:])
			tempPositionInfo.Description = (&goquery.Selection{Nodes: []*html.Node{selection.Nodes[index+2]}}).Text()
			response <- tempPositionInfo
			wg.Done()
		}(index)
	}
	wg.Wait()
	close(response)
	wgResponse.Wait()
	return res, nil
}

// joinSelectorData selector
func joinSelectorData(selectorData ...*SelectorData) string {
	var strSlice []string
	for _, v := range selectorData {
		strSlice = append(strSlice, v.Selector)
	}
	temp := strings.Join(strSlice, ",")
	return temp[:len(temp)-1]
}

func writeFile(positionInfo []*PositionInfo) error {
	output, err := json.MarshalIndent(&struct {
		Total        int             `json:"total"`
		PositionInfo []*PositionInfo `json:"position_info"`
	}{len(positionInfo), positionInfo}, "", "\t\t")
	if err != nil {
		klog.Errorln(err)
		return err
	}
	err = ioutil.WriteFile("xsky-crawler.json", output, 0644)
	if err != nil {
		klog.Errorln(err)
		return err
	}
	return nil
}
