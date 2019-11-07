package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)


func IPParse(ipstr string) []string{
	var iplist []string
	//支持逗号分隔
	ips:=strings.Split(ipstr,",")
	for _,ip_temp :=range ips{

		ipaddr:=net.ParseIP(ip_temp)
		//标准IP校验
		if ipaddr!=nil{
			iplist=append(iplist, ipaddr.String())
		}

		//CIDR格式解析
		ipaddr,ipnet,err:=net.ParseCIDR(ip_temp)

		if err==nil{
			for ip := ipaddr.Mask(ipnet.Mask); ipnet.Contains(ip); incip(ip) {
				iplist=append(iplist,ip.String())
			}
		}

	}

	return iplist
}

func incip(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}


func PortParse(ports string)([]int,error){
	var err error
	var portlist []int
	//有,号优先,分割处理
	portsItem:=strings.Split(ports,",")
	for _,port :=range portsItem{

		//1-65535解析
		if strings.Contains(port,"-"){
			port:=strings.Split(port,"-")

			portA,err:=strconv.Atoi(port[0])
			if err!=nil{
				return nil,err
			}
			portB,err:=strconv.Atoi(port[1])

			if err!=nil{
				return nil,err
			}

			for i:=portA;i<=portB;i++{
				//默认会访问的端口443 80,不需要写进列表
				if i!=443 && i!=80{
					portlist=append(portlist,i)
				}
			}

		}else{
			portNumber,err:=strconv.Atoi(port)
			if err!=nil{

				return portlist,err
			}

			//默认会访问的端口443 80,不需要写进列表
			if portNumber!=443 && portNumber!=80{
				portlist=append(portlist,portNumber)
			}
		}
	}



	return portlist,err

}

//读取IP或host列表
func ListReader(path string,opr string) []string{
	var list []string
	fi,err:=os.Open(path)

	if err!=nil{
		log.Printf("Error: %s\n", err)
		return nil
	}
	defer fi.Close()
	br:=bufio.NewReader(fi)

	for{
		line,_,err:=br.ReadLine()

		if err==io.EOF{
			break
		}
		line_str:=strings.TrimSpace(string(line))
		if opr=="ip"{ //如果是ip行,就解析ip行地址,支持每行写法是:ip,ip,ip,CIDR或者单个ip
			list=append(list,IPParse(line_str)...)
		}else if opr=="host"{
			list=append(list,line_str)
		}
	}

	return list

}
// Slice string去重
func SliceStringUnique(s []string) []string {
	result := []string{}
	m := make(map[string]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}
// Slice int去重
func SliceIntUnique(s []int) []int {
	result := []int{}
	m := make(map[int]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}