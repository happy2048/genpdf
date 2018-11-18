package main
import(
	"github.com/jessevdk/go-flags"
	"net/http"
	"fmt"
	"io"
	"bytes"
	"os"
	"log"
	"encoding/json"
	"io/ioutil"
)
type ReturnData struct {
	Msg string `json: "msg"`
	Data string `json: "data"`
	Code string `json: "code"`
}
type PdfClient struct {
	Content string 
	Type string
	Args string
}
type Options struct {
	Type string `short:"t" long:"type" description:"give the resource type(url: the url of html,body: a html file \nwith no <html>,<head> and <body> tags,complete: a complete html file)." default:"body"`
	Name string `short:"T" long:"title" description:"when type is body,you can give a title with match the content." default:""`
	Host string `short:"H" long:"host" description:"give the server ip which is running wkhtmltopdf." default:"127.0.0.1"`
	Port string `short:"P" long:"port" description:"give the server service listen port." default:"6660"`
	PdfArgs string `short:"a" long:"args" description:"give the wkhtmltopdf args,eg: '--outline --disable-internal-links'." default:""`
	Doc bool `short:"d" long:"wkhtmltopdf-args" description:"print the wkhtmltopdf args."`

}
func main() {
	opt,args := NewOptions()
	opt.Check(args)
	PostGeneratePdfReq(opt.Name,args[0],opt.PdfArgs,opt.Type,opt.Host,opt.Port,args[1])
}
func PostGeneratePdfReq(name,con,args,typ,server,port,out string) {
	url := "http://" + server + ":"+ port + "/generate"
	data := make(map[string]string)
	if typ == "body" || typ == "complete" {
		tmpdata,err := ioutil.ReadFile(con)
		con = string(tmpdata)
		if err != nil {
			log.Printf("read file error,reason: %s\n",err.Error())
			return 
		}
	}
	data["name"] = name
	data["content"] = con
	data["type"] = typ
	data["args"] = args
	bytesData,err := json.Marshal(data)
	if err != nil {
		log.Printf("json marshal failed,reason: %s\n",err.Error())
		return
	}
	redata,err := Operate("POST",url,bytesData)
	if err != nil {
		log.Printf("http request failed,reason: %s\n",err.Error())
		return 
	}
	var parse ReturnData
	err = json.Unmarshal([]byte(redata),&parse)
	if err != nil {
		log.Printf("json unmarshal return data failed,reason: %s\n",err.Error())
		return
	}
	if parse.Code == "1000" {
		url = "http://" + server + ":"+ port + "/pdf/" + parse.Data
		res, err := http.Get(url)
		if err != nil {
			log.Printf("get pdf file failed,reason: %s\n",err.Error())
			return
		}
		f,err := os.Create(out)
		if err != nil {
			log.Printf("get pdf file failed,reason: %s\n",err.Error())
			return

		}
		io.Copy(f,res.Body)
	}else {
		log.Printf("get pdf file failed,reason: %s\n",parse.Msg)
	}
}

func GetOsEnv(env string) string {
    return os.Getenv(env)
}
func Operate(method,url string,data []byte) (string,error) {
    client := &http.Client{}
    var request *http.Request
    var err error
    if string(data) == "" {
        request,err = http.NewRequest(method,url,nil)
    }else {
        request,err = http.NewRequest(method,url,bytes.NewReader(data))
    }
    request.Header.Set("Connection", "keep-alive")
	if method == "POST" {
		request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	}
    response,err := client.Do(request)
    if err != nil {
        return "",err
    }
    if response.StatusCode == 200 {
        body,err := ioutil.ReadAll(response.Body)
        if err != nil {
         return "",err
        }
        return string(body),nil
    }
    return "",fmt.Errorf("%s","requst failure")
}
func NewOptions() (*Options,[]string) {
    var options Options
	pdata := flags.NewParser(&options, flags.Default)
	pdata.Usage = "[OPTIONS] INPUT [OUTPUT FILE]"
    args,err := pdata.Parse()
    if err != nil {
        if flagsErr, ok := err.(*flags.Error);ok && flagsErr.Type == flags.ErrHelp {
            os.Exit(0)
        }else {
			fmt.Println(err.Error())
        	os.Exit(1)
		}
    }
	if len(args) == 1 {
		args = append(args,"generate.pdf")
	}
    return &options,args
}


func (opt *Options) Check(args []string) {
    if len(os.Args) == 1 {
        fmt.Printf("Error: you should give some options,plese use -h or --help to get usage.\n")
        os.Exit(1)
    }
	if opt.Doc == true {
		data,err := Operate("GET","https://wkhtmltopdf.org/usage/wkhtmltopdf.txt",[]byte(""))
		if err != nil {
			fmt.Printf("Error: get the wkhtmltopdf args failed,reason: %s\n",err.Error())
			os.Exit(2)
		}
		fmt.Println(data)
		os.Exit(0)
	}
	if opt.Type != "complete" && opt.Type != "body" && opt.Type != "url" {
		fmt.Printf("Error: unknown type,its' value should be in [complete,body,url].\n")
		os.Exit(1)
	}
	if len(args) == 0 {
		fmt.Printf("Error: you should give the input resource (a html file or a url) and the name of output pdf file.\n")
		os.Exit(2)
	}
	if opt.Type == "complete" || opt.Type == "body" {
		exist := CheckFileExist(args[0])
		if exist == false {
			fmt.Printf("Error: the file %s does not exist.",args[0])
			os.Exit(3)
		}
	}	

}

func CheckFileExist(filename string) bool {
    var exist = true
    _,err := os.Stat(filename)
    if os.IsNotExist(err) {
        exist = false
    }
    return exist
}

