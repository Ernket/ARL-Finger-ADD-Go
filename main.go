package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

/*
删除所有指纹这一个操作其实是从finger_num开始的
添加指纹是从make_file开始
*/

// 定义与 YAML 文件匹配的结构体
type ARLConfig struct {
	ARLConfig struct {
		URL      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Threads  int    `yaml:"threads"`
	} `yaml:"arl_config"`
}

/*
var (
	loginURL       = flag.String("url", "", "URL地址")
	loginName      = flag.String("username", "", "用户名")
	loginPassword  = flag.String("password", "", "密码")
	thread_num     = flag.Int("thread", 10, "线程数")
	del_all_finger = flag.Bool("n", false, "删除所有指纹")
)
*/

var (
	//searchName      = flag.String("s", "", "搜索的项目名")
	del_all_finger  = flag.Bool("d", false, "删除所有指纹")
	add_file_finger = flag.Bool("a", false, "添加finger.json文件中的指纹")
)

// 登录需要用到的头信息
var headers = map[string]string{
	"Accept":     "application/json, text/plain, */*",
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36",
	"Connection": "close",
	//"Token":           token,
	"Accept-Encoding": "gzip, deflate",
	"Accept-Language": "zh-CN,zh;q=0.9",
	"Content-Type":    "application/json; charset=UTF-8",
}
var client = createClient()

// 自定义帮助信息
func customUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-d|-a]\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "选项:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = customUsage
	flag.Parse()

	/*
		flag.Usage = customUsage
		flag.Parse()
		loginURL := *loginURL
		loginName := *loginName
		loginPassword := *loginPassword

		// 检查必要的参数是否已提供
		if loginURL == "" || loginName == "" || loginPassword == "" {
			flag.Usage() // 如果参数缺失，显示帮助信息
			fmt.Println("All of -url, -username, and -password are required.")
			return
		}
	*/
	yamlData, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}
	// 解析 YAML 文件
	var config ARLConfig
	err = yaml.Unmarshal(yamlData, &config)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
	loginURL := config.ARLConfig.URL
	fmt.Printf("API url: %s\n", loginURL)
	loginName := config.ARLConfig.Username
	loginPassword := config.ARLConfig.Password
	thread_num := config.ARLConfig.Threads

	// 登录
	token, err := login(loginURL, loginName, loginPassword)
	if err != nil {
		fmt.Println("Login failed:", err)
		return
	}
	fmt.Println("[+] Login Success!!")

	// 登录成功后的token写到等会要用的头部
	headers["Token"] = token

	// 这部分的功能是删除所有指纹
	if *del_all_finger {
		finger_num(loginURL)
		return
	}

	if *add_file_finger {
		make_file(loginURL, thread_num)
	}

}

func addFinger(name, rule, url string) {
	url = fmt.Sprintf("%s/api/fingerprint/", url)
	data := map[string]string{"name": name, "human_rule": rule}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	// 创建一个HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataJSON))
	if err != nil {
		fmt.Println("请求创建失败:", err)
		return
	}
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()
	// 打印响应内容
	if resp.StatusCode == 200 {
		fmt.Printf("Add: [+] %s\n", dataJSON)
	} else {
		fmt.Printf("请求失败，状态码：%d\n", resp.StatusCode)
	}
}
func delFinger(url string, allIDs []string) {
	// 创建一个map来存储JSON对象
	var jsonMap = make(map[string]interface{})
	// 创建_id字段的切片
	jsonMap["_id"] = allIDs
	DelData, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("转换JSON数据失败: ", err)
		return
	}

	url = fmt.Sprintf("%s/api/fingerprint/delete/", url)
	// 创建一个HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(DelData))
	if err != nil {
		fmt.Errorf("请求创建失败:", err)
		return
	}
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("请求失败:", err)
	}
	defer resp.Body.Close()
	//resp, err = http.Post(url, "application/json", bytes.NewBuffer(DelData))

	if resp.StatusCode != http.StatusOK {
		fmt.Println("请求接口删除指纹失败:", resp.Status)
		return
	}
	fmt.Println("[+] 所有指纹已删除")
}

func finger_num(url string) {
	page_num := 1
	// 存储id
	var allIDs []string
	for {
		one_item, check_for, err := finger_id(url, page_num)
		if err != nil {
			fmt.Println("发生报错: ", err)
			return
		}
		for _, item := range one_item {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				return

			}
			id, ok := itemMap["_id"].(string)
			if !ok {
				return
			}
			allIDs = append(allIDs, id)
		}
		if !check_for {
			break
		}

		page_num++
	}
	fmt.Println("当前获取的id数：", len(allIDs))
	delFinger(url, allIDs)
}

func finger_id(url string, page_num int) ([]interface{}, bool, error) {
	page_size := 500
	url = fmt.Sprintf("%sapi/fingerprint/?page=%d&size=%d&order_name=update_date", url, page_num, page_size)
	// 创建一个HTTP请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, true, fmt.Errorf("请求创建失败:", err)
	}
	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, true, fmt.Errorf("请求失败:", err)
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, fmt.Errorf("读取响应失败:", err)
	}
	// 解析JSON到map
	var result map[string]interface{}
	//fmt.Printf("%s", body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, true, fmt.Errorf("JSON解析失败:", err)
	}
	//finger_num := result["total"].(float64)
	//fmt.Println("当前指纹数量: ", int(finger_num))

	// 获取items数组
	items, ok := result["items"].([]interface{})
	if !ok {
		return nil, true, fmt.Errorf("Error asserting items field as []interface{}:")
	}
	idCount := len(items)
	fmt.Println("获取到了", idCount, "条指纹")
	if idCount < page_size {
		return items, false, nil
	}
	return items, true, nil
}

func login(url, username, password string) (string, error) {
	loginData := map[string]string{"username": username, "password": password}
	loginDataJSON, err := json.Marshal(loginData)
	if err != nil {
		return "", fmt.Errorf("Error marshaling JSON: %v", err)
	}
	loginURL := fmt.Sprintf("%sapi/user/login", url)
	// 发起POST请求
	resp, err := client.Post(loginURL, "application/json", bytes.NewBuffer(loginDataJSON))
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}
	// 检查resp是否为nil，不检查好像也没问题
	if resp == nil {
		return "", fmt.Errorf("Response is nil")
	}
	// 检查resp.Body是否为nil
	if resp.Body == nil {
		return "", fmt.Errorf("Response body is nil")
	}
	// 执行defer语句，确保在函数返回之前关闭响应体
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Login failed: %s", body)
	}

	var loginResp map[string]interface{}
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		return "", fmt.Errorf("Error parsing JSON response: %v", err)
	}

	token, ok := loginResp["data"].(map[string]interface{})["token"].(string)
	if !ok {
		return "", fmt.Errorf("Token not found in response")
	}
	// 登录成功后的token返回
	return token, nil
}

// 定义一个函数来创建并返回一个HTTP客户端
func createClient() *http.Client {
	//proxyURL, err := url.Parse("http://127.0.0.1:8080")
	//if err != nil {
	//	fmt.Println("设置代理出错:", err)
	//	return nil
	//}
	jar, _ := cookiejar.New(nil) // 忽略错误处理
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 关闭证书认证
		//	Proxy:           http.ProxyURL(proxyURL),
	}
	return &http.Client{
		Jar:       jar,
		Transport: tr,
	}
}
func make_file(loginURL string, thread_num int) {
	// 创建信号量
	semaphore := make(chan struct{}, thread_num)
	var wg sync.WaitGroup
	// 读取JSON文件并解析内容
	file, err := os.ReadFile("./finger.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	// 解析JSON文件
	var loadDict map[string]interface{}
	err = json.Unmarshal(file, &loadDict)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	// 根据JSON中的规则添加指纹
	for _, finger := range loadDict["fingerprint"].([]interface{}) {
		wg.Add(1) // 增加等待的线程数
		go func(finger interface{}) {
			semaphore <- struct{}{} // 获取信号量
			defer func() {
				<-semaphore // 释放
				wg.Done()   // 完成一个线程
			}()
			// 处理finger.json中的数据
			fingerMap := finger.(map[string]interface{})
			name := fingerMap["cms"].(string)
			method := fingerMap["method"].(string)
			location := fingerMap["location"].(string)
			keywordInterface := fingerMap["keyword"].([]interface{})
			keywordSlice := make([]string, len(keywordInterface))
			for i, v := range keywordInterface {
				keywordSlice[i] = v.(string)
			}
			var rule string
			if method == "keyword" {

				if location == "body" {
					rule = fmt.Sprintf("body=\"%s\"", strings.Join(keywordSlice, "\",\""))
				} else if location == "title" {
					rule = fmt.Sprintf("title=\"%s\"", strings.Join(keywordSlice, "\",\""))
				} else if location == "header" {
					rule = fmt.Sprintf("header=\"%s\"", strings.Join(keywordSlice, "\",\""))
				}

			} else if method == "icon_hash" {
				rule = fmt.Sprintf("icon_hash=\"%s\"", strings.Join(keywordSlice, "\",\""))
			}

			addFinger(name, rule, loginURL) // 调用addFinger函数写入指纹到ARL
		}(finger)
	}
	wg.Wait() // 等待所有线程完成
}
