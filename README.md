### gateway

```
功能简介
根据平台配置相应的协议信息及点表,部署后采集数据上报
```

### 项目结构
```
├── app
│   ├── external 服务部署在服务器 对web提供接口
│   └── interior docker方式部署在网关 只接收来自服务器的请求
├── common 通用方法
├── deploy docker 部署 
```

### 注意事项
```
goctl 参数类型暂时没有实现 map[string]interface{}
```

### 环境安装
```
  golang go1.18 版本
  安装包下载地址
  http://mirrors.nju.edu.cn/golang/
  
  goctl 版本1.3.5
  安装命令  GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/go-zero/tools/goctl@v1.3.5
  
  goctl-go-compact
  安装命令  GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/goctl-go-compact@latest
  
  goctl-swagger
  git clone https://github.com/zeromicro/goctl-swagger.git
  修改 go.mod 与 goctl 版本一致
  github.com/zeromicro/go-zero/tools/goctl v1.3.5
  go build -o goctl-swagger main.go
  mv goctl-swagger $GOPATH/bin
  
  项目tpl配置
  goctl template clean 清除
  cd $HOME/.goctl
  goctl template init 初始化``
```


### 项目构建

```
1. 在 app/ 下面微服务目录 例如 gateway/app/external 目录执行，生成api目录结构
goctl api go -api ./desc/gateway.api -dir ./
```

### model生成
```
goctl model mysql ddl -src ./model/ne_device_data_access_config.sql -dir ./model
goctl model mysql ddl -src ./model/ne_device_data_access_item.sql -dir ./model
goctl model mysql ddl -src ./model/fe_edge_device.sql -dir ./model
goctl model mysql ddl -src ./model/ne_device.sql -dir ./model

ne_device 设备表
fe_edge_device  网关表
ne_device_data_access_config 协议配置表 
ne_device_data_access_item 协议点表

```

### 文档生成

```
goctl api plugin -plugin goctl-swagger="swagger -filename gateway.json -basepath /api" -api ./app/external/api/desc/gateway.api -dir ./doc
```

### 代码监测工具
```
1. goimports -l -w ./  代码 格式化工具
go install golang.org/x/tools/cmd/goimports@latest 安装命令

2. go vet ./...  静态代码分析工具，用于检测代码中的潜在问题

3. errcheck ./... 检查未处理的err 
go install github.com/kisielk/errcheck@latest 安装命令

4.gocyclo . # 当前项目所有方法的复杂度拍讯
  gocyclo main.go # 单个文件复杂度排序
  gocyclo -top 10 src/ # src 文件下最复杂的10个方法
  gocyclo -over 25 src # src 目录下复杂度大于 25 的方法
  gocyclo -avg . # 当前项目的平均复杂度
  gocyclo -top 20 -ignore "_test|pb.go|Godeps|vendor/" . # 忽略部分文件

  go install github.com/fzipp/gocyclo/cmd/gocyclo@latest 安装命令
  复杂度说明
  1-10：简单函数，易于理解和维护。
  10-15：中等复杂度，仍然可以接受，但可能需要注意。
  15-20：较高的复杂度，建议重构以降低复杂度。
  20+：非常高的复杂度，通常需要重构。
  建议不要超过20
```

### docker 部署
```
打包镜像
Dockerfile 应在项目第一层级

docker build -f Dockerfile -t gateway:v1 .

运行
docker run -id gateway:v1 /bin/bash

开机自启
docker update --restart always 容器ID或名称

进入
docker exec -it d7bda8428193 /bin/bash

docker-compose up -d 后台启动

指定文件运行
docker-compose -f docker-compose.yml up -d 

查看容器log
docker-compose logs -f protocol
protocol 为 <service-name>与docker-compose.yml 对应

tag
docker tag gateway:v1 xunjikeji.com.cn:5000/gateway:v1

推送到私有库
docker push xunjikeji.com.cn:5000/gateway:v1

拉取镜像 
docker pull xunjikeji.com.cn:5000/gateway:v1

私有库镜像查看地址
https://xunjikeji.com.cn:5000/v2/_catalog
```
### 项目部署
```
直接在服务器运行编译后的文件即可
例如
/home/gowork/src/gateway/app/interior/api/gateway -f /home/gowork/src/gateway/app/interior/api/etc/pro/apis.yaml
如需开机自启动将此命令添加在 /etc/rc.local 文件下即可
注意 必须使用绝对路径

### 问题记录
1.一个设备下有两个协议 其中一个协议部署成功一个协议失败 设备是部署成功还是失败
2.删除网关、删除设备、是否要停掉网关里的采集协议
3.一个协议下发成功后第二次下发失败是否改为为下发状态
4. 只有下发未下发是否需要其他状态、设备状态是否要修改
5.设备侧删除协议点部署是否删除网关内该设备协议点表
##########
1. 只存网关一份
是否有必要将协议点表相关的配置存储两份建议直接存储在网关
物模型选择性上报web

2.只存服务器一份
理论上如果网关访问不到服务器数据也无法上报
缺点 网关断电重启后如果拿不到服务器数据无法采集数据并保存在本地

TODO 判断采集非法值问题
```