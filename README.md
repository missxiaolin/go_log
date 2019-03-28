# go_log

### 架构图

<img src="http://missxiaolin.com/log._20190326png.png">


### nginx 配置


~~~
server{
	listen  80;
	server_name  www.miss-log.com;

	access_log  /Users/web/go/log/miss-log.log;

	location = /dig {
		empty_gif;
		error_page 405 =200 $request_uri;
	}


}
~~~

### 日志格式

~~~
log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';
~~~