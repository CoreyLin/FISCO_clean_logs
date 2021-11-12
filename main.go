package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func runLogCron() {
	logCron := cron.New(cron.WithSeconds())
	_, err := logCron.AddFunc("@every 1m", cleanNodesLogs)
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Start cleaning log cron job")
	logCron.Start()
	defer logCron.Stop()
	// 永远阻塞，让函数不退出
	select {}
}

func cleanNodesLogs() {
	// os.Args[1]就是FISCO nodes所在的path
	// 列出当前path下所有的目录
	fileInfos, err := ioutil.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range fileInfos {
		if !f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), "node") {
			cleanOldLog(os.Args[1], f.Name())
		}
	}
}

func cleanOldLog(nodePath string, nodeName string) {
	// 获取node log所在的路径
	logPath := strings.TrimRight(nodePath, string(os.PathSeparator)) + string(os.PathSeparator) + nodeName + string(os.PathSeparator) + "log"
	// 列出node log路径中所有的log，并且去掉最近的5条log
	fileInfos, err := ioutil.ReadDir(logPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 把文件夹都过滤掉，只剩下文件
	var logFileInfos []os.FileInfo
	for _, f := range fileInfos {
		if f.IsDir() {
			continue
		}
		logFileInfos = append(logFileInfos, f)
	}

	// 如果logFileInfos的长度小于等于5，就不删除任何log，直接退出
	if len(logFileInfos) <= 5 {
		return
	}

	// 把logFileInfos按时间排序，时间新的排在前面，即按时间倒序排列
	logFileInfos = sortByTime(logFileInfos)

	// 去掉logFileInfos里修改时间最近的5个文件
	for i := 5; i < len(logFileInfos); i++ {
		fmt.Println("remove ", logFileInfos[i].Name())
		err := os.Remove(logPath + string(os.PathSeparator) + logFileInfos[i].Name())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func sortByTime(infos []os.FileInfo) []os.FileInfo {
	sort.Slice(infos, func(i, j int) bool {
		flag := false
		if infos[i].ModTime().After(infos[j].ModTime()) {
			flag = true
		}
		return flag
	})
	return infos
}

func main() {
	// 定时删除log文件是一个定时任务，程序需一直在后台运行，每隔一段时间执行一次
	// 程序执行需要输入一个参数，即FISCO nodes所在的目录
	runLogCron()

	// FISCO nodes目录下的所有节点都需要删除，每个节点不删完，如果把正在写的log删掉，那么log就不能继续写下去了，保留5个最新的log文件不删，其他的全部删除

}
