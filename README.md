##genpdf:

本软件基于wkhtmltopdf，做成了web server形式。

使用方法：

###安装genpdf-server


直接用docker启动：

	[root@localhost ~]# docker run --name genpdf -p 6660:6660 happy365/genpdf

如果想监听在其他端口可以使用如下的启动命令：

	[root@localhost ~]# docker run  --name genpdf -p 6661:6661 -e PORT=6661 happy365/genpdf

需要说明的是，因为在生成pdf的过程中会产生大量的临时文件，所以需要定时清理这些文件，在宿主机（不是容器内）上执行如下命令：

	[root@localhost ~]# curl http://127.0.0.1:6660/deletefiles

如果需要定时清理，建议把这条命令写进/etc/crontab中：

	[root@localhost ~]# echo "1 1 * * * root curl http://127.0.0.1:6660/deletefiles" >> /etc/crontab

上面这条命令表示，每天晚上一点钟执行清理操作。

上面这条api只能在宿主机或者容器内以"127.0.0.1"或"localhost"为ip进行调用。


###安装客户端


	[root@localhost ~]# git clone https://github.com/happy2048/genpdf.git
	[root@localhost ~]# cp genpdf/genpdf /usr/bin
	[root@localhost ~]# genpdf -h
	Usage:
	  genpdf [OPTIONS] INPUT [OUTPUT FILE]

	Application Options:
	  -t, --type=             give the resource type(url: the url of html,body: a html file
							  with no <html>,<head> and <body> tags,complete: a complete html
							  file). (default: body)
	  -T, --title=            when type is body,you can give a title with match the content.
	  -H, --host=             give the server ip which is running wkhtmltopdf. (default:
							  127.0.0.1)
	  -P, --port=             give the server service listen port. (default: 6660)
	  -a, --args=             give the wkhtmltopdf args,eg: '--outline
							  --disable-internal-links'.
	  -d, --wkhtmltopdf-args  print the wkhtmltopdf args.

	Help Options:
	  -h, --help              Show this help message

下面对选项进行一些说明：

使用方式：

	genpdf [OPTIONS] INPUT [OUTPUT FILE]

	OPTIONS: 相关选项
	INPUT:资源，如果资源类型为url，那么这这里跟一个url链接，如果资源类型为body或者complete，那么跟一个html文件
	OUTPUT FILE: 输出的PDF文件，如果没指定，默认的会在当前目录下生成一个叫generate.pdf


-t: 指定资源类型，这里有三种类型：

 * url: 表示需要打印的内容是一个url地址，根据这个地址打印该网页。比如： http://www.google.com

 * complete: 表示一个完整的html文件。

 * body: 表示需要打印的内容是一个html文件，这个文件不是一个完整的html文件，只有<body></body>标签内的东西，例如：

	```html
		<div>
		<span><img src="test.png"><p>这是一个测试文件</p></span>
		</div>
	```

-T: 当资源类型为body时，可以指定文档的标题，这个标题会被打印在首页。

-H：server服务器的IP地址

-P：server服务器的端口

-a: 指定wkhtmltopdf的相关选型，具体信息可以到其官网查看，或者使用`genpdf -d`查看

-d: 打印wkhtmltopdf相关选项

###使用示例

1.打印www.baidu.com这个网页：

	[root@localhost ~]# genpdf  -t url www.baidu.com  test.pdf

2.打印资源类型为complete的html文件：

	[root@localhost ~]# genpdf -t complete  /mnt/test.html  test.pdf

3.打印资源类型为body的html文件：

	[root@localhost ~]# genpdf -t body /mnt/body.html test.pdf

###使用API调用：

可以使用http方式调用：

	curl -s -X POST \
	http://localhost:6660/generate \
	-H "content-type: application/json" \
	-d "{
		"name": "test",
		"type": "body",
		"content": "<p>测试当中</p>",
		"args": ""
	}"

调用成功后会会返回生成的pdf文件名，这个文件存放在服务器上，需要我们执行下载操作，文件假设为mytest.pdf：

	curl http:/localhosts:6660/pdf/mytest.pdf
