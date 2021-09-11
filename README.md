# S3Uploader (for MINIO)

```
git clone https://github.com/harryzhu/s3uploader
cd s3uploader
go build --ldflags "-w -s"
```

S3Uploader 是一个面向后端为 `MINIO` 服务 *S3* 接口提供文件 上传（fput）、下载（fget）、状态查看（stat）的小工具，不提供删除和批量删除功能 

`./s3uploader fput｜fget｜stat --bucket=test --object=s3/object/path_or_name.png --file=local/path/of/your/file.png ` 

## 全局参数：
`--prefix=` ：	指定环境变量中的前缀，用`--prefix`可以支持大量不同环境的配置，比如用`alpha`前缀指定测试环境，用`sand`前缀表示开发环境

`--bucket=` ：	指定后端 MINIO 的 S3接口的 `bucket` 名称，无默认值，需要明确指定

`--object=` ：	指定存储在S3中的对象名称

`--file=` ：	指定本地文件路径

`--debug` ：	指定日志级别为debug级别，显示更多日志，默认为`false`，仅出现warn、error、fatal级别的日志时才显示

`--logfile=` ：	指定日志保存路径


### `fput`：上传本地文件到S3中保存

`--mime=` ：	指定文件到类型，可选参数

### `fget`： 下载S3中的对象文件保存到本地

### `tag`：  给S3中的对象打标签，最多打64个标签，键值对格式 `key:value` 用`:`分隔，多个键值对用`,`分隔，最长128字节（128个ascii字符或者42个中文字符）
`--kv="project:BeChangedByTheWorld,owner:harryzhu,year:2020"`  
标签将分开显示为 `project:bechangedbytheworld`,`owner:harryzhu`,`year:2020`  

### `stat`： 不下载文件仅查看S3对象的元数据，该命令下，如果 `--file=` 指定了路径，则会把默认显示在console中的信息保存到指定的文件中

## 用法：

* 设置环境变量(可以设置在`/etc/profile`中也可以在脚本中指定），默认前缀`S3UPLOADER_`:

``` 
	export S3UPLOADER_USERNAME=YOUR_API_KEY
	export S3UPLOADER_PASSWORD=YOUR_API_SECRET
	export S3UPLOADER_ENDPOINT=YOUR_S3.DOMAIN.COM
	export S3UPLOADER_USESSL=false
	export S3UPLOADER_LOGFILE=/Users/harryzhu/logs/s3uploader.log
```

* 多环境支持，例如	--prefix=ALPHA，那么ALPHA环境的前缀就是：`ALPHA_S3UPLOADER_`

```
	export ALPHA_S3UPLOADER_USERNAME=YOUR_API_KEY
	export ALPHA_S3UPLOADER_PASSWORD=YOUR_API_SECRET
	export ALPHA_S3UPLOADER_ENDPOINT=YOUR_ALPHA_S3.DOMAIN.COM
	export ALPHA_S3UPLOADER_USESSL=false
	export ALPHA_S3UPLOADER_LOGFILE=/Users/harryzhu/logs/s3uploader_alpha.log
```

* `fput`上传本地文件到S3:  
 `./s3uploader fput --bucket=test --object=s3/object/path_or_name.png --file=local/path/of/your/file.png --mime="image/png"`

* `fget`下载S3中的文件到本地:  
 `./s3uploader fget --bucket=test --object=s3/object/path_or_name.png --file=local/path/of/your/file.png`

* `stat`查看S3中的对象的元数据:  
 `./s3uploader stat --bucket=test --object=s3/object/path_or_name.png`

* `stat`命令默认打印结果在终端，如果想要保存查询到的结果到文件，用 `--file=`指定即可，文件后缀名必须为`.json`，如果不是，会自动添加`.json`  
 `./s3uploader stat --bucket=test --object=s3/object/path_or_name.png --file=local/path/for/saving/result.json`

* `tag`给S3中的对象添加标签，`--kv=`指定，格式k:v为一对，使用`,`分割多个不同的标签对,键、值都将被转换为小写:  
 `./s3uploader stat --bucket=test --object=s3/object/path_or_name.png --kv="project:BeChangedByTheWorld,owner:harryzhu,year:2020"`

* 使用`ALPHA`环境:  
 `./s3uploader fput --prefix=ALPHA --bucket=test --object=s3/object/path_or_name.png --file=local/path/of/your/file.png --mime="image/png"`


* 使用`ALPHA`环境，且显示更多日志:  
 `./s3uploader fget --prefix=ALPHA --debug --bucket=test --object=s3/object/path_or_name.png --file=local/path/of/your/file.png` 

