package ctrl

import (
	"../util"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func init(){
	os.MkdirAll("./mnt", os.ModePerm)
}

func Upload(w http.ResponseWriter, r *http.Request){
	UploadLocal(w,r)
}

//1.存储位置 ./mnt，需要确保已经创建好
//2.url格式 /mnt/xxxx.png 需要确保网络能够访问/mnt/
func UploadLocal(w http.ResponseWriter, r *http.Request){
	//获得上传的源文件s
	srcfile, head, err := r.FormFile("file")
	if err!=nil{
		util.RespFail(w, err.Error())
	}
	//创建一个新文件d
	suffix := ".png"
	//如果前端文件名称包含后缀，如果前端指定了filetype
	ofilename := head.Filename
	tmp := strings.Split(ofilename,".")
	if len(tmp)>1{
		suffix = "." + tmp[len(tmp)-1]
	}
	//如果前端制定了filetype
	//formdata.append("filetype",".png")
	filetype := r.FormValue("filetype")
	if len(filetype)>0{
		suffix = filetype
	}
	//time.Now().Unix()时间戳
	filename := fmt.Sprintf("%d%04d%s", time.Now().Unix(),rand.Int31(), suffix)
	//先创建新文件
	dstfile,err := os.Create("./mnt/"+filename)
	if err!=nil{
		util.RespFail(w, err.Error())
		return
	}
	//将源文件内容copy到新文件
	_,err = io.Copy(dstfile, srcfile)
	if err!=nil{
		util.RespFail(w, err.Error())
		return
	}
	//将新文件路径转换为url地址
	url := "/mnt/"+filename
	//响应到前端
	util.RespOK(w,url,"")

}
