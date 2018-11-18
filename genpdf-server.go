package main
import(
	"os/exec"
	"math/rand"
	"path"
	"time"
	"bytes"
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"os"
	"io/ioutil"
	"log"
)
type HtmlInfo struct {
	Type string `json: "type"`
	Content string `json: "content"`
	Args string `json: "args"`
	Name string `json: "name"`
}
type FileInfo struct {
	FileName string `json:"filename"`
	Content  string `json: "content"`
}
type ReturnData struct {
	Msg string `json: "msg"`
	Data string `json: "data"`
    Code  string `json: "code"`
}
func main() {
	if GetOsEnv("X_SERVER_CMD") == "" {
		os.Setenv("X_SERVER_CMD","Xvfb :42 -screen 0 1024x768x24 -ac +extension GLX +render -noreset")
	}
	if GetOsEnv("TMP_PATH") == "" {
		os.Setenv("TMP_PATH","/tmp/pdf")
	}
	if GetOsEnv("TEMPLATE_FILE") == "" {
		os.Setenv("TEMPLATE_FILE","/usr/local/html-template/template_file.txt")
	}
	if GetOsEnv("TEMPLATE_FILE_WITH_NAME") == "" {
		os.Setenv("TEMPLATE_FILE_WITH_NAME","/usr/local/html-template/template_file_with_name.txt")
	}
	if GetOsEnv("PORT") == "" {
		os.Setenv("PORT","6660")
	}
	if  state,_ := PathExists(GetOsEnv("TMP_PATH"));state == false {
		os.Mkdir(GetOsEnv("TMP_PATH"),0644)
	}
	go RunCmd(GetOsEnv("X_SERVER_CMD"))
	http.HandleFunc("/generate",HandleHtml)
	http.HandleFunc("/deletefiles",HandleDeleteTmpFiles)
	http.Handle("/pdf/", http.StripPrefix("/pdf/", http.FileServer(http.Dir(GetOsEnv("TMP_PATH") + "/"))))
	log.Printf("service is starting")
	err := http.ListenAndServe(":" + GetOsEnv("PORT"),nil)
	if err != nil {
		log.Printf("Error: ",err.Error())
		return 
	}
}
		
func HandleDeleteTmpFiles(w http.ResponseWriter,req *http.Request) {
	if req.Method == "GET" {
		if req.Host != "localhost:" + GetOsEnv("PORT") && req.Host != "127.0.0.1:" + GetOsEnv("PORT") {
			ReturnValue(w,"1100","","you have not privilege to execute this option.")
			return
			
		}
		_,_,err := RunCmd("rm -rf " + GetOsEnv("TMP_PATH") + "/*")
		if err != nil {
			ReturnValue(w,"1100","",err.Error())
			return
		}
		ReturnValue(w,"1000","","delete tmp files succeed.")
		return
	}else {
		ReturnValue(w,"1100","","invalid method")
		return
	}
}
func HandleHtml(w http.ResponseWriter,req *http.Request)  {
	if req.Method == "POST" {
		body,_ := ioutil.ReadAll(req.Body)
		var htmlInfo HtmlInfo
		if err := json.Unmarshal(body,&htmlInfo);err == nil {
			pdf,err := CreatePdf(htmlInfo.Name,htmlInfo.Type,htmlInfo.Content,htmlInfo.Args)
			if err != nil {
				ReturnValue(w,"1100","",pdf)
				return
			}
			ReturnValue(w,"1000",pdf,"")
		}else {
			ReturnValue(w,"1100","",err.Error())
			return 
		}
	}else {
		ReturnValue(w,"1100","","http method is invalid,you should use POST.")
		return
	}
}
func CreatePdf(name,typ,content,args string) (string,error) {
	fid := GetRandomString(10)
	html := fid + ".html"
	pdf := fid + ".pdf"
	tmpPath := GetOsEnv("TMP_PATH")
	var cmdStr string
	if tmpPath == "" {
		return "",fmt.Errorf("%s","tmp path is not found.")
	}
	if typ == "body" || typ == "complete"  {
		if typ == "body" {
			if name == "" {
				data,err := ioutil.ReadFile(GetOsEnv("TEMPLATE_FILE"))
				if err != nil {
					return "",err
				}
				content = fmt.Sprintf(string(data),content)
			}else {
				data,err := ioutil.ReadFile(GetOsEnv("TEMPLATE_FILE_WITH_NAME"))
				if err != nil {
					return "",err
				}
				content = fmt.Sprintf(string(data),name,content)
			}
		}
		err := ioutil.WriteFile(path.Join(tmpPath,html),[]byte(content),os.ModeAppend)
		if err != nil {
			return "",err
		}	
		cmdStr = "wkhtmltopdf.sh " + args + " " + path.Join(tmpPath,html) + " " + path.Join(tmpPath,pdf) 
	}else if typ == "url" {
		cmdStr = "wkhtmltopdf.sh " + args + " " + content + " " + path.Join(tmpPath,pdf)
	}else {
		return "",fmt.Errorf("unknown type %s,the type should be in [complete,body,url].",typ)
	}
	out,_,err := RunCmd(cmdStr)
	log.Printf("\n%s\n",out)
	if err != nil {
		return out,err
	}
	return pdf,err
}
func GetOsEnv(env string) string {
    return os.Getenv(env)
}
func ReturnValue(w http.ResponseWriter,code string,data string,msg string) {
    redata,err := json.Marshal(&ReturnData{Code: code,Msg: msg,Data: data})
    if err != nil {
		log.Printf("%s",err.Error())
        return
    }
    io.WriteString(w,string(redata))
}
func RunCmd(cmdStr string) (string,string,error) {
	cmd := exec.Command("/bin/sh","-c",cmdStr)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	return out.String(),errOut.String(),err
}

func GetRandomString(length int64) string{
   str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
   bytes := []byte(str)
   result := []byte{}
   r := rand.New(rand.NewSource(time.Now().UnixNano()))
   for i := int64(0); i < length; i++ {
      result = append(result, bytes[r.Intn(len(bytes))])
   }
   return string(result)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
