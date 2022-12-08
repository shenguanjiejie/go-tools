package tools

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

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
		Logln(CallerLevel(1), "%v time cost = %v\n", signs, tc)
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

// LoadFile 加载文件
func LoadFile(path string) []byte {
	file, err := os.Open(path)
	defer func() {
		err = file.Close()
		if err != nil {
			Logln("read err:", err)
		}
	}()
	if err != nil {
		Logln(err)
	}
	byteValue, _ := ioutil.ReadAll(file)
	return byteValue
}

// SaveFile 生成文件
func SaveFile(path string, data []byte) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		Logln(err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		Logln(err)
	}
}

// MD5 md5
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

var isOnline = true

// InternetCheck 网络检测是否在线
func InternetCheck(logLevel LogLevel, switchAction ...func(isOnline bool)) bool {
	res := new(http.Response)
	err := HttpGet("http://connect.rom.miui.com/generate_204", nil, res, &HttpConfig{Log: logLevel})

	if err != nil || res.StatusCode != 204 {
		Logln(Condition(err != nil), err)
		if isOnline {
			isOnline = false
			if len(switchAction) > 0 {
				go switchAction[0](false)
			}
		}
		return false
	}

	if !isOnline {
		isOnline = true
		if len(switchAction) > 0 {
			go switchAction[0](true)
		}
	}
	return true
}
