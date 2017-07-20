package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	//	"time"
)

/*
var url = "http://ieqi.net/wp-content/themes/ieqi-wptheme/css/main.css"
var destDir = "F:\\kuaipan\\Projects\\webcopyer\\template"

var img_ext = []string{".jpg", ".gif", ".png", ".jpeg", ".JPG", ".GIF", ".PNG", ".JPEG", ".ico", ".ICO"}
var css_ext = []string{".css", ".less", ".CSS", ".LESS"}
var js_ext = []string{".js", ".JS"}
var other_ext = []string{".dll"}
var extArray = []string{".jpg", ".gif", ".png", ".jpeg", ".css", ".less", ".js", ".ico", ".JPG", ".GIF", ".PNG", ".JPEG", ".CSS", ".LESS", ".JS", ".ICO", ".dll"}

var css_resource_list []string
*/
type Config map[string]string

var (
	cfg Config

	current_url       string
	destDir           string
	css_resource_list []string

	img_ext   []string
	css_ext   []string
	js_ext    []string
	other_ext []string

	html_dir   string
	img_dir    string
	css_dir    string
	js_dir     string
	other_dir  string
	configFile *string
	help       *bool
	version    *bool
	extArray   []string

	method string
	arg1   string
	arg2   string

	ext string
)

func init() {

	configFile = flag.String("config", "config", "配置文件路径")
	help = flag.Bool("help", false, "查看帮助")
	version = flag.Bool("version", false, "show version")

	flag.Parse()
	if *help {
		fmt.Println()
		fmt.Println("复制模板：")
		fmt.Println("    webcopyer http://lyric.im/")
		fmt.Println("    webcopyer get http://lyric.im/")
		fmt.Println()
		fmt.Println("拷贝本地模板：")
		fmt.Println("    webcopyer getLocal F:\\template\\index.html  http://drizzlep.diandian.com/")
		fmt.Println()
		fmt.Println("拷贝html及其直接关联资源模板：")
		fmt.Println("    webcopyer getHtml http://lyric.im/")
		fmt.Println()
		fmt.Println("拷贝css及其直接关联资源模板：")
		fmt.Println("    webcopyer getCss http://lyric.im/styles/style.css")
		fmt.Println()
		fmt.Println("关于：")
		fmt.Println("    version： 1.0")
		fmt.Println("    author：  肖立群")
		fmt.Println("    email：   xlqstar@gmail.com")
		fmt.Println("    repo：    https://github.com/xlqstar/WebCopyer.git")
		fmt.Println()
		os.Exit(0)
	}
	if *version {
		fmt.Println()
		fmt.Println("关于：")
		fmt.Println("    version： 1.0")
		fmt.Println("    author：  肖立群")
		fmt.Println("    email：   xlqstar@gmail.com")
		fmt.Println("    repo：    https://github.com/xlqstar/WebCopyer.git")
		fmt.Println()
		os.Exit(0)
	}
	method = flag.Arg(0)
	arg1 = flag.Arg(1)
	arg2 = flag.Arg(2)

	url := flag.Arg(0)
	if method == "getLocal" {
		url = flag.Arg(2)
	}
	if method == "get" {
		url = flag.Arg(1)
	}
	log.SetFlags(log.Ltime)
	//或取配置文件

	cfg := Configure(*configFile)

	destDir = cfg.getStr("destDir")

	img_ext = cfg.getArray("img_ext")
	css_ext = cfg.getArray("css_ext")
	js_ext = cfg.getArray("js_ext")
	other_ext = cfg.getArray("other_ext")

	html_dir = fixResPath(cfg.getStr("html_dir"))
	img_dir = fixResPath(cfg.getStr("img_dir"))
	css_dir = fixResPath(cfg.getStr("css_dir"))
	js_dir = fixResPath(cfg.getStr("js_dir"))
	other_dir = fixResPath(cfg.getStr("other_dir"))

	extArray = arrayMerge(extArray, img_ext)
	extArray = arrayMerge(extArray, css_ext)
	extArray = arrayMerge(extArray, js_ext)
	extArray = arrayMerge(extArray, other_ext)

	destDir = checkAndMkDir(destDir, url)

}

func main() {

	// fmt.Println(get_destdir_and_filetype("http://www.chinaz.com/"))

	//getAll("http://drizzlep.diandian.com/", "")        //获取在线模版
	//getAll("F:\\kuaipan\\Projects\\webcopyer\\template\\index.html", "http://drizzlep.diandian.com/page/2/") //获取本地模版

	//完整测试
	// getHtml("http://www.chinaz.com/")

	//测试download
	// down_resource("https://ss.cnnic.cn/seallogo.dll?sn=e12020335020010628301467&size=3")

	//css测试 常用
	// getCss("http://img.chinaz.com/templates/chinaz/css/style.css?v=20121121")

	//css测试 import
	//getCss("http://download.csdn.net/css/download.css")

	//fix_url 测试
	// current_url = "http://ieqi.net/"
	// fmt.Println(fix_url("./sdfs/sdfa/wahaha.jpg"))

	if method == "get" {
		if arg1 == "" {
			log.Fatal("请正确输入参数")
		}
		getAll(arg1, "") //获取在线模版
	} else if method == "getLocal" || method == "getlocal" {
		if arg1 == "" || arg2 == "" {
			log.Fatal("请正确输入参数")
		}
		getAll(arg1, arg2) //获取本地模版
	} else if method == "getCss" || method == "getcss" {
		if arg1 == "" {
			log.Fatal("请正确输入参数")
		}
		getCss(arg1)
	} else if method == "getHtml" || method == "gethtml" {
		if arg1 == "" {
			log.Fatal("请正确输入参数")
		}
		getHtml(arg1)
	} else {
		if method == "" {
			log.Fatal("请正确输入参数")
		}
		getAll(method, "") //获取在线模版
	}

}

func getAll(file string, url string) {
	var html string
	ext = "html"
	if strings.HasPrefix(file, "http") || strings.HasPrefix(file, "https") {
		html = http_get(file)
		current_url = file
	} else {
		htmlbyte, _ := ioutil.ReadFile(file)
		html = string(htmlbyte)
		current_url = url
	}

	html = extruct_html_resource(html)
	ioutil.WriteFile(destDir+html_dir+get_true_filename(current_url), []byte(html), 0777) //转储html
	for k := range css_resource_list {
		current_url_tmp := current_url
		current_url = fix_url(css_resource_list[k])
		css := http_get(current_url)
		css = extruct_css_resource(css) //1、改引用路径 2、下载资源
		dest := getDestDir("css")
		write_erro := ioutil.WriteFile(destDir+dest+get_true_filename(current_url), []byte(css), 0777) //转储css
		if write_erro != nil {
			fmt.Println(write_erro)
		}
		current_url = current_url_tmp
	}

	log.Println("Done !!")
}

func getHtml(url string) {
	ext = "html"
	html := http_get(url)
	current_url = url
	html = extruct_html_resource(html)
	ioutil.WriteFile(destDir+html_dir+get_true_filename(url), []byte(html), 0777) //转储css
	log.Println("Done !!")
}

func getCss(url string) {
	ext = "css"
	css := http_get(url)
	current_url = url
	css = extruct_css_resource(css) //1、改引用路径 2、下载资源
	dest := getDestDir("css")
	ioutil.WriteFile(destDir+dest+get_true_filename(url), []byte(css), 0777) //转储css
	log.Println("Done !!")
}

//==============================================
//css提取
func extruct_css_resource(content string) string {
	re1, _ := regexp.Compile("background(-image)?\\s*:[^;})]*url\\s*\\((\\S+)?\\)") //src\s*=\s*(["|'])([\S\s]+?)\1
	content = re1.ReplaceAllStringFunc(content, css_change_path_1)

	re2, _ := regexp.Compile("@import.*url\\s*(\\(.*\\));") //src\s*=\s*(["|'])([\S\s]+?)\1
	content = re2.ReplaceAllStringFunc(content, append_import_css_resource)
	return content
}

func css_change_path_1(str string) string {
	re, _ := regexp.Compile("background(-image)?\\s*:[^;})]*url\\s*\\((\\S+)?\\)")
	old_path := re.FindStringSubmatch(str)[2]
	old_path = strings.Trim(old_path, "\"' ")
	new_path := css_change_path(str, old_path)
	return new_path
}

func append_import_css_resource(str string) string {
	re, _ := regexp.Compile("@import.*url\\s*\\((.*)\\);")
	import_url := re.FindStringSubmatch(str)[1]
	import_url = fix_url(strings.Trim(import_url, "\"' "))
	getCss(import_url)
	return str
}

//@import.*url\s*(\(.*\));
func css_change_path(str string, old_path string) string {

	dir, filetype := get_destdir_and_filetype(old_path)
	if dir != "unexcept" {
		down_resource(old_path, dir)
		true_filename := get_true_filename(old_path)
		relPath := getRelPath(css_dir, getDestDir(filetype))
		return strings.Replace(str, old_path, relPath+"/"+true_filename, -1)
	}
	return str

}

//================================================
//htmlt资源提取
func extruct_html_resource(content string) string {
	re1, _ := regexp.Compile("src\\s*=\\s*\"([\\S\\s]+?)\"")
	re2, _ := regexp.Compile("src\\s*=\\s*'([\\S\\s]+?)'")                          //src\s*=\s*(["|'])([\S\s]+?)\1
	re3, _ := regexp.Compile("<\\s*link([\\S\\s]+?)href\\s*=\\s*\"([\\S\\s]+?)\"")  //src\s*=\s*(["|'])([\S\s]+?)\1
	re4, _ := regexp.Compile("<\\s*link([\\S\\s]+?)href\\s*=\\s*'([\\S\\s]+?)'")    //src\s*=\s*(["|'])([\S\s]+?)\1
	re5, _ := regexp.Compile("background(-image)?\\s*:[^;})]*url\\s*\\((\\S+)?\\)") //src\s*=\s*(["|'])([\S\s]+?)\1
	content = re1.ReplaceAllStringFunc(content, html_change_path_1)
	content = re2.ReplaceAllStringFunc(content, html_change_path_2)
	content = re3.ReplaceAllStringFunc(content, html_change_path_3)
	content = re4.ReplaceAllStringFunc(content, html_change_path_4)
	content = re5.ReplaceAllStringFunc(content, html_css_change_path)
	return content
}

func html_change_path_1(str string) string {
	re, _ := regexp.Compile("src\\s*=\\s*\"([\\S\\s]+?)\"")
	old_path := re.FindStringSubmatch(str)[1]
	new_path := html_change_path(str, old_path)
	return new_path
}

func html_change_path_2(str string) string {
	re, _ := regexp.Compile("src\\s*=\\s*'([\\S\\s]+?)'")
	old_path := re.FindStringSubmatch(str)[1]
	new_path := html_change_path(str, old_path)
	return new_path
}
func html_change_path_3(str string) string {
	re, _ := regexp.Compile("<\\s*link([\\S\\s]+?)href\\s*=\\s*\"([\\S\\s]+?)\"")
	old_path := re.FindStringSubmatch(str)[2]
	new_path := html_change_path(str, old_path)
	return new_path
}

func html_change_path_4(str string) string {
	re, _ := regexp.Compile("<\\s*link([\\S\\s]+?)href\\s*=\\s*'([\\S\\s]+?)'")
	old_path := re.FindStringSubmatch(str)[2]
	new_path := html_change_path(str, old_path)
	return new_path
}

func html_change_path(str string, old_path string) string {
	dir, filetype := get_destdir_and_filetype(old_path)
	if dir != "unexcept" {
		down_resource(old_path, dir)
		if filetype == "css" {
			css_resource_list = append(css_resource_list, old_path)
		}
		true_filename := get_true_filename(old_path)
		resDir := getDestDir(filetype)
		relPath := getRelPath(html_dir, resDir)
		return strings.Replace(str, old_path, relPath+"/"+true_filename, -1)
	}
	return str
}

func html_css_change_path(str string) string {
	re, _ := regexp.Compile("background(-image)?\\s*:[^;})]*url\\s*\\((\\S+)?\\)")
	old_path := re.FindStringSubmatch(str)[2]
	old_path = strings.Trim(old_path, "\"' ")
	new_path := _html_css_change_path(str, old_path)
	return new_path
}

func _html_css_change_path(str string, old_path string) string {

	dir, filetype := get_destdir_and_filetype(old_path)
	if dir != "unexcept" {
		down_resource(old_path, dir)
		true_filename := get_true_filename(old_path)
		resDir := getDestDir(filetype)
		relPath := getRelPath(html_dir, resDir)
		return strings.Replace(str, old_path, relPath+"/"+true_filename, -1)
	}
	return str

}

//======================================

func http_get(url string) string {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}

func get_true_filename(url string) string {
	filename := filepath.Base(url)
	extname := filepath.Ext(url)
	for k := range extArray {
		if strings.Contains(extname, extArray[k]) {
			true_filename := strings.Replace(filename, extname, extArray[k], -1)
			return true_filename
		}
	}
	true_filename := strings.Split(filename, "?")[0]

	if len(true_filename) > 0 {
		if !(url == arg1 || url == method || url == arg2) {
			log.Println("发现未知类型文件：" + url)
			log.Println()
		} else {
			true_filename = true_filename + "." + ext
		}
		return true_filename
	}

	return filename
}

func get_destdir_and_filetype(url string) (string, string) {
	filetype := ""
	dir := destDir

	if in_array(filepath.Ext(url), img_ext) {
		filetype = "images"
	} else if in_array(filepath.Ext(url), css_ext) {
		filetype = "css"
	} else if in_array(filepath.Ext(url), js_ext) {
		filetype = "js"
	} else if in_array(filepath.Ext(url), other_ext) {
		filetype = "other"
	} else {
		return "unexcept", "unexcept"
	}

	dir = dir + "/theme/" + filetype + "/"
	return dir, filetype
}

func down_resource(url string, destDir string) {

	url = strings.Trim(url, " \t\n\r")
	filename := get_true_filename(url)
	fullfilename := destDir + filename
	fixed_url := fix_url(url) //矫正url（相对路径转换为绝对路径）

	log.Println("FROM: " + fixed_url)
	log.Println("TO:   " + fullfilename)

	if strings.HasPrefix(fixed_url, "http") || strings.HasPrefix(fixed_url, "https") {
		resp, err := http.Get(fixed_url)
		if err != nil {
			fmt.Println(">>>>>>>>>>>>>>>>>>以下资源获取出错，跳过ing<<<<<<<<<<<<<<<<<<<<<<")
			fmt.Println(err)
			fmt.Println("_______________________________________________________________")
			//os.Exit(0)
		} else {

			out, create_err := os.Create(fullfilename)
			if create_err != nil {
				fmt.Println(create_err)
				os.Exit(0)
			}
			_, copy_err := io.Copy(out, resp.Body)
			if copy_err != nil {
				fmt.Println(copy_err)
				os.Exit(0)
			}
			out.Close()
			resp.Body.Close()
		}
	}
	fmt.Println()
}

func in_array(v string, array []string) bool {
	for k := range array {
		if strings.Contains(strings.ToLower(v), strings.ToLower(array[k])) {
			return true
		}
	}
	return false
}

//http://www.baidu.com/wahaha/sdfds/index.html
func fix_url(url string) string {

	re1, _ := regexp.Compile("http[s]?://[^/]+")
	destrooturl := re1.FindString(current_url)

	//当url为：//wahaha/xiaoxixi/tupian.png
	if strings.HasPrefix(url, "//") {
		url = "http:" + url
	} else if strings.HasPrefix(url, "/") {
		// re1,_ := regexp.Compile("http[s]?://[^/]+")
		// destrooturl := re1.FindString(current_url)
		url = destrooturl + url
	}

	//当url为："../wahaha/xiaoxixi/tupian.png"、"./wahaha/xiaoxixi/tupian.png"、"wahaha/xiaoxixi/tupian.png"
	if !strings.HasPrefix(url, "/") && !strings.HasPrefix(url, "http") && !strings.HasPrefix(url, "https") {
		// current_url = strings.TrimSuffix(current_url, "/")
		if destrooturl == current_url {
			url = current_url + "/" + url
		} else {
			re2, _ := regexp.Compile("[^/]+?$")
			url = re2.ReplaceAllString(current_url, "") + url
		}

	}

	return url
}

func Configure(filePath string) Config {
	configByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("配置文件不存在！")
	}
	configMap := Config{}
	config := strings.Replace("\r\n", "\n", string(configByte), -1)
	configArray := strings.Split(config, "\n")
	for _, v := range configArray {
		v = Trim(v)
		vArray := strings.SplitN(v, ":", 2)
		if len(vArray) == 2 && !strings.HasPrefix(v, "#") {
			key := Trim(vArray[0])
			value := Trim(vArray[1])
			configMap[key] = value
		}

	}
	return configMap
}

func (cfg Config) getInt(key string) int {
	value, err := strconv.Atoi(cfg[key])
	if err != nil {
		log.Fatal(key + "值未填写或填写不正确，请确认为整数")
	}
	return value
}

func (cfg Config) getStr(key string) string {
	if cfg[key] == "" {
		log.Fatal(key + "值未填写或填写不正确")
	}
	return cfg[key]
}

func (cfg Config) getArray(key string) []string {
	return strings.Split(cfg[key], "|")
}

func Trim(s string) string {
	return strings.Trim(s, " \t\n\r")
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func arrayMerge(old []string, other []string) []string {
	for _, v := range other {
		old = append(old, v)
	}
	return old
}

func checkAndMkDir(destDir string, url string) string {
	//判断destDir是否存在 不存在的话创建 失败的话报错终止
	if !exist(destDir) {
		err := os.Mkdir(destDir, os.ModePerm)
		if err != nil {
			log.Fatal("无法创建目录！")
		}
	}

	destDir = destDir + "/" + get_true_filename(url)[0:len(get_true_filename(url))-1] + "/" // + strconv.Itoa(int(time.Now().Unix()))
	fmt.Println(destDir)
	if !exist(destDir) {
		err := os.Mkdir(destDir, os.ModePerm)
		if err != nil {
			log.Fatal("无法创建目录！")
		}
	}

	if !exist(destDir + "/theme/js") {
		err := os.MkdirAll(destDir+"/theme/js", os.ModePerm)
		if err != nil {
			log.Fatal("无法创建目录！")
		}
	}

	if !exist(destDir + "/theme/images") {
		err := os.MkdirAll(destDir+"/theme/images", os.ModePerm)
		if err != nil {
			log.Fatal("无法创建目录！")
		}
	}

	if !exist(destDir + "/theme/css") {
		err := os.MkdirAll(destDir+"/theme/css", os.ModePerm)
		if err != nil {
			log.Fatal("无法创建目录！")
		}
	}
	return destDir
}

func getDestDir(filetype string) string {

	path := ""
	if filetype == "images" {
		path = img_dir
	} else if filetype == "js" {
		path = js_dir
	} else if filetype == "css" {
		path = css_dir
	} else if filetype == "other" {
		path = other_dir
	} else if filetype == "html" {
		path = html_dir
	}

	return path
}

//检查修复各个资源存储路径
func fixResPath(path string) string {

	fixed := path

	if path == "" || path == "/" {
		fixed = "/"
	} else {
		if !strings.HasSuffix(path, "/") {
			fixed = fixed + "/"
		}
	}

	if !strings.HasPrefix(path, "/") {
		fixed = "/" + path
	}

	return fixed
}

func getRelPath(base string, target string) string {
	rel, _ := filepath.Rel(base, target)
	rel = strings.Replace(rel, "/\\", "/", -1)
	rel = strings.Replace(rel, "\\", "/", -1)
	return rel
}
