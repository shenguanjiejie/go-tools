package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"
)

// Slog 带行号输出
func Slog(a ...interface{}) {
	// 获取上层调用者PC，文件名，所在行
	pc, _, codeLine, ok := runtime.Caller(1)

	if !ok {
		fmt.Print(a...)
	}

	// 根据PC获取函数名
	format, slice := formatWithValues(pc, codeLine, slogFormat(a), a...)
	fmt.Printf(format, slice...)
}

// Slogln 带行号输出
func Slogln(a ...interface{}) {
	// 获取上层调用者PC，文件名，所在行
	pc, _, codeLine, ok := runtime.Caller(1)

	if !ok {
		fmt.Println(a...)
	}

	// 根据PC获取函数名
	format, slice := formatWithValues(pc, codeLine, slogFormat(a), a...)
	fmt.Printf(format+"\n", slice...)
}

// Slogf 带行号格式输出
func Slogf(format string, a ...interface{}) {
	// 获取上层调用者PC，文件名，所在行
	pc, _, codeLine, ok := runtime.Caller(1)

	if !ok {
		fmt.Printf(format, a...)
	}

	// 根据PC获取函数名
	finalFormat, slice := formatWithValues(pc, codeLine, format, a...)
	fmt.Printf(finalFormat, slice...)
}

func slogFormat(a ...interface{}) string {
	formatStr := ""
	dataArr, ok := a[0].([]interface{})
	if !ok {
		formatStr += "%v, "
	}

	for _, data := range dataArr {
		switch data.(type) {
		case nil, bool, int, int64, string, error:
			formatStr += "%v, "
		default:
			formatStr += "%#v, "
		}
	}
	return formatStr
}

func formatWithValues(pc uintptr, codeLine int, format string, a ...interface{}) (string, []interface{}) {
	funName := runtime.FuncForPC(pc).Name()
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	slice := []interface{}{timeStr, funName, codeLine}
	slice = append(slice, a...)
	return "%s--%s--第%d行--: " + format, slice
}

// RemoveDuplicateElementArray 去重
func RemoveDuplicateElementArray(sourceArray []string) []string {
	result := make([]string, 0, len(sourceArray))
	temp := map[string]int{}
	for _, item := range sourceArray {
		if _, ok := temp[item]; !ok {
			temp[item] = 1
			result = append(result, item)
		}
	}
	return result
}

// TimeCost @brief：耗时统计函数
func TimeCost(signs ...string) func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		slogFormat("%v time cost = %v\n", signs, tc)
	}
}

// WaitHandle 阻塞型协程队列, 所有参数必传才执行
func WaitHandle(channel chan interface{}, goCount int, waitingFor func(), asyncHandle func(channelObj interface{})) {
	if channel == nil || asyncHandle == nil || waitingFor == nil {
		return
	}

	waitGroup := new(sync.WaitGroup)

	for i := 0; i < goCount; i++ {
		waitGroup.Add(1)
		go func() {
			for {
				obj, ok := <-channel
				if !ok {
					break
				}

				asyncHandle(obj)
			}
			waitGroup.Done()
		}()
	}
	waitingFor()
	close(channel)
	waitGroup.Wait()
}

// LoadJSON 加载json文件
func LoadJSON(path string) []byte {
	jsonFile, err := os.Open(path)
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			Slogln("read json err:", err)
		}
	}()
	if err != nil {
		Slogln(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}

// SaveJSON 生成json文件
func SaveJSON(path string, data []byte) {
	jsonFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		Slogln(err)
	}
	defer jsonFile.Close()
	_, err = jsonFile.Write(data)
	if err != nil {
		Slogln(err)
	}
}

func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}
