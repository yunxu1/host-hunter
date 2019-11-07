# host-hunter

探测只有绑定指定IP才能访问的域名,主要用于信息收集使用.

<p><b>Build：</b></p>

```
go build 
```

<p><b>Help:</b></p>

```shell
./hosthunter -h
  -h	help
  -hL string  #加载指定的host列表文件
    	load domain list file path
  -host string #设置host
    	target domain eg:xxx.com
  -iL string #加载指定的ip列表文件
    	load ip list file path
  -ip string #加载指定ip,支持CIDR解析,多个逗号分隔
    	input ip range ,eg:CIDR or ip or ip,ip,ip
  -output string #输出结果到指定文件,默认hostinfo.txt
    	output result to file (default "hostinfo.txt")
  -port string #指定端口,多个端口逗号分隔,支持1-65535范围解析,80/443默认加载无需指定
    	port eg:8080,8888,81 or 1-65535 or port
    	
  -redirect #是否开启30x跳转,默认关闭
    	follow 30x redirect
  -thread int # 设置线程数
    	thread default 5 (default 5)
  -timeout int #设置请求超时时间
    	timeout default 3 (default 3)
  -code string # 显示设置的http响应码,默认只显示结果为200的状态码，多个状态码逗号分隔
    	http status code filter options eg:200,201,500 or 200 (default "200")
```



##### 测试:

###### 使用所有参数扫描:

```shell
./hosthunter -hL domain.txt -iL iplist.txt -host w.xxx.com -port 88,8000,8080,9090-9200 -thread 15 -timeout 2
```



###### 仅加载文件扫描:

```shell
./hosthunter -hL domain.txt -iL iplist.txt -thread 15
```



###### 指定IP和host扫描:

```shell
./hosthunter -ip 192.168.1.0/24,192.168.2.2 -host w.xxx.com
```



###### 回显字段:

```
状态码 body长度 server x-power-by content-type 请求域名-IP 页面标题 
```



