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

~~~
go get github.com/sirupsen/logrus
go get github.com/streadway/amqp
~~~

### amqp 手册

[手册链接](https://godoc.org/github.com/streadway/amqp)

RMQ官网提供的教程：[https://www.rabbitmq.com/getstarted.html](https://www.rabbitmq.com/getstarted.html) 
go-amqp库函数手册：[https://godoc.org/github.com/streadway/amqp](https://godoc.org/github.com/streadway/amqp)

### 安装rabbitmq

~~~
docker pull rabbitmq:management
docker run -d --name rabbitmq -p 5671:5671 -p 5672:5672 -p 4369:4369 -p 25672:25672 -p 15671:15671 -p 15672:15672 rabbitmq:management
~~~

通过命令可以看出，一共映射了三个端口，简单说下这三个端口是干什么的。 
5672：连接生产者、消费者的端口。 
15672：WEB管理页面的端口。 
25672：分布式集群的端口。 
