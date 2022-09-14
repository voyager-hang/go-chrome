package go_chrome

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
)

func getCtX(windowName, url string, opt ...GoChromeOptions) (context.Context, context.CancelFunc, *GoChromeOptions, string) {
	ctxInit := context.WithValue(context.Background(), "windowsName", windowName)
	runOpt := GoChromeOptions{
		ChromeRunCommand: ChromeRunCommand,
	}
	if len(opt) > 0 {
		runOpt = opt[0]
		runOpt.ChromeRunCommand = ChromeRunCommand
		chromeExecPath := getChromeExecPath(&runOpt)
		runOpt.chromeExecPath = chromeExecPath
		runOpt.ChromeRunCommand["headless"] = runOpt.CliModule
		gh := getHttp(&runOpt)
		url = gh.GetUrl(url, runOpt.UseHttpServer)
		if runOpt.AppModule || strings.HasPrefix(url, "data:text/html") {
			runOpt.ChromeRunCommand["app"] = url
		}
		if runOpt.WindowWidth > 0 && runOpt.WindowHeight > 0 {
			runOpt.ChromeRunCommand["window-size"] = fmt.Sprintf("%d,%d", runOpt.WindowWidth, runOpt.WindowHeight)
		}
		if runOpt.WindowPositionWidth > 0 && runOpt.WindowPositionHeight > 0 {
			runOpt.ChromeRunCommand["window-position"] = fmt.Sprintf("%d,%d", runOpt.WindowPositionWidth, runOpt.WindowPositionHeight)
		}
	} else {
		gh := getHttp(&runOpt)
		url = gh.GetUrl(url, runOpt.UseHttpServer)
		chromeExecPath := getChromeExecPath(&runOpt)
		runOpt.chromeExecPath = chromeExecPath
	}
	execAllocatorOption := []chromedp.ExecAllocatorOption{}
	if runOpt.chromeExecPath != "" {
		execAllocatorOption = append(execAllocatorOption, chromedp.ExecPath(runOpt.chromeExecPath))
	}
	for k, v := range runOpt.ChromeRunCommand {
		execAllocatorOption = append(execAllocatorOption, chromedp.Flag(k, v))
	}
	ctx, _ := chromedp.NewExecAllocator(
		ctxInit,
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			execAllocatorOption...,
		)...,
	)
	newCtx, cancel := chromedp.NewContext(
		ctx,
		// 设置日志方法
		chromedp.WithLogf(log.Printf),
	)
	return newCtx, cancel, &runOpt, url
}

func getChromeExecPath(opt *GoChromeOptions) string {
	chromeExecPath := ""
	packConfig := getPackageConfig()
	if packConfig != nil && packConfig.ChromeExecPath != "" {
		// test
		return packConfig.ChromeExecPath
	}
	if opt.RestoreAssets != nil {
		path, err := GetBrowserPath(opt)
		if err == nil {
			return path + "\\chrome-win\\chrome.exe"
		}
	}
	return chromeExecPath
}
