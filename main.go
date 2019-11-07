package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var(
	schemas[2] string=[2] string{"http","https"}
	hL string
	iL string
	ip string
	host string
	h bool
	thread int
	addressList []string
	hostList []string
	threadCh chan bool
	tasklist []HostIP
	timeout int
	port string
	ports =[]int{80,443} //默认端口访问
	redirect bool
	flog *os.File
	outfile string
	wg sync.WaitGroup
	code string
)

func init(){

	log.SetOutput(ioutil.Discard)
	flag.BoolVar(&h, "h", false, "help")
	flag.StringVar(&hL,"hL","","load domain list file path")
	flag.StringVar(&iL,"iL","","load domain list file path")
	flag.StringVar(&port,"port","","port eg:8080,8888,81 or 1-65535 or port")
	flag.StringVar(&ip,"ip","","input ip range ,eg:CIDR or ip or ip,ip,ip")
	flag.StringVar(&host,"host","","target domain eg:xxx.com")
	flag.IntVar(&thread,"thread",5,"thread default 5")
	flag.IntVar(&timeout,"timeout",3,"timeout default 3")
	flag.BoolVar(&redirect, "redirect", false, "follow 30x redirect")
	flag.StringVar(&outfile,"output","hostinfo.txt","output result to file")
	flag.StringVar(&code,"code","200","http status code filter options eg:200,201,500 or 200")
}


type HostIP struct {
	Schema string
	Address string
	Host string
	Port int
}

func HostVerify(hostip HostIP) {

	defer func() {
		<-threadCh
		wg.Done()
	}()
	var target string
	var titleinf string=""

	tr := &http.Transport{
		Dial: (&net.Dialer{

			//解决连接过多出现err:too many open file.
			// https://colobu.com/2016/07/01/the-complete-guide-to-golang-net-http-timeouts/
			// http://craigwickesser.com/2015/01/golang-http-to-many-open-files/
			Timeout:   time.Duration(timeout) * time.Second,
			Deadline:  time.Now().Add(time.Duration(timeout) * time.Second),
			KeepAlive: time.Duration(timeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout:time.Duration(timeout)* time.Second,
		//忽略证书校验
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}


	client:=&http.Client{
		Timeout:time.Duration(timeout)*time.Second,
		Transport:tr,
	}
	//关闭redirect
	if !redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	if hostip.Port==443{
		log.Printf("%s,%d\n",hostip.Schema,port)
		target=fmt.Sprintf("https://%s",hostip.Address)
	}else if hostip.Port==80{
		log.Printf("%s,%d",hostip.Schema,port)
		target=fmt.Sprintf("http://%s",hostip.Address)
	}else{

		target=fmt.Sprintf("%s://%s:%d",hostip.Schema,hostip.Address,hostip.Port)
	}
	log.Printf("请求URL:%s\n",target)
	req,err:=http.NewRequest(http.MethodGet,target,nil)
	if err!=nil{
		log.Printf("error:%s\n",err)
		return
	}

	req.Host=hostip.Host
	req.Header.Set("Connection","close") // 解决 too many open files
	resp,err:=client.Do(req)

	if err!=nil{
		log.Printf("error:%s\n",err)
		return
	}

	defer resp.Body.Close()

	banner:=resp.Header.Get("Server")
	//contentlength := resp.Header.Get("Content-Length")
	x_power_by:=resp.Header.Get("X-Powered-By")
	contenttype:=""
	pair := strings.SplitN(resp.Header.Get("Content-Type"), ";", 2)
	if len(pair) == 2 {
		contenttype= pair[0]
	}

	body,err:=ioutil.ReadAll(resp.Body)

	if err!=nil{
		log.Printf("error:%s\n",err)
		return
	}

	re:=regexp.MustCompile("<title>(.*)<\\/title>")

	title:=re.FindSubmatch(body)

	if len(title)>0{
		titleinf=string(title[len(title)-1])

		if titleinf ==""{
			titleinf="空标题"
			log.Printf("%s\n",title)
		}
	}

	urlout:=fmt.Sprintf("%s://%s:%d-%s",hostip.Schema,hostip.Host,hostip.Port,hostip.Address)
	//out:=fmt.Sprintf("%-70s%-20s%-63s%-9s%-10s%-5s%-5s\n",urlout,status,titleinf,banner,x_power_by,contentlength,contenttype)
	out:=fmt.Sprintf("%-5d %-15d %-60s %-35s %-15s %-50s %s\n",resp.StatusCode, len(body) ,banner, x_power_by, contenttype,urlout, titleinf)
	var d *color.Color

	codeStr:=strconv.Itoa(resp.StatusCode)
	if strings.Contains(code,codeStr){

		if strings.HasPrefix(codeStr,"2"){
			d=color.New(color.FgHiGreen,color.Bold)
		}else if strings.HasPrefix(codeStr,"3"){
			d=color.New(color.FgCyan,color.Bold)
		}else{
			d=color.New(color.FgHiRed,color.Bold)
		}
		d.Printf(out)
		flog.WriteString(out)
	}



}
//所有协议 ip 和 host组合成任务列表
func MakeTask(hostList []string , addressList []string){
	for _,host :=range hostList {
		for _,addr:=range addressList{
			for _,schema:=range schemas{
				for _,port :=range ports{

					tasklist=append(tasklist, HostIP{schema,addr,host,port})
				}
			}
		}
	}

}
//多线程扫描
func HostScan(){
	for _,task :=range tasklist{

		threadCh<-true
		wg.Add(1)
		go HostVerify(task)

	}
	wg.Wait()
}

func main() {

	flag.Parse()

	if h{
		flag.Usage()
	}

	flog, _ = os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC,0755)//打印输出到文件

	//加载host列表
	if hL!=""{
		hostList=SliceStringUnique(append(hostList,ListReader(hL,"host")...))

	}

	//加载adderss列表
	if iL!=""{
		addressList=SliceStringUnique(append(addressList,ListReader(iL,"ip")...))

	}
	//加载用户输入的IP段或者IP
	if ip!=""{
		addressList=SliceStringUnique(append(addressList,IPParse(ip)...))
	}

	if port!=""{
		portNumbers,err:=PortParse(port)
		if err!=nil{

			fmt.Println(err)
			return
		}
		ports=SliceIntUnique(append(ports,portNumbers...))
	}
	//加载用户输入的host
	if host!=""{
		hostList=SliceStringUnique(append(hostList,host))
	}




	log.Printf("%v\n",addressList)
	log.Printf("%v\n",hostList)
	log.Printf("%v\n",ports)

	MakeTask(hostList,addressList) //加载任务列表
	//设置线程数
	threadCh=make(chan bool,thread)
	HostScan();
}
