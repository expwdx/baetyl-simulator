# Baetyl节点模拟器

## 构建可执行文件

- 构建
```shell
go build -o output/bsctl
```
or
```shell
make
```

- 拷贝可执行文件到 `~/go/bin`或者`/usr/bin`目录

```shell
cp output/bsctl ~/go/bin
```
如果未将GOBIN加入到全局环境变量，可拷贝到`/usr/bin`目录

## 配置
配置文件路径: `etc/conf.yaml`，采用相对路径存放。示例如下：

```yaml
test:
  planTime:
  nodeCount: 100
  nodeStartNo: 0

user:
  deploy:
    interval: "120s"
  read:
    interval: "30s"

engine:
  report:
    interval: "10s"
  desire:
    interval: "5s"

mock:
  namespace: "baetyl-cloud"
  nodeNamePrefix: "stress-test"
  nodeCount: 100
  nodeStartNo: 0
  nodeLabels:
    type: "vhost"
    use: "stress-test"
  appName: "mysql-simulator"

cloud:
  admin:
    schema: "http"
    host: "127.0.0.1:9004"
    apiVer: "v1"
    timeout: "30s"
  init:
    schema: "https"
    host: "127.0.0.1:9003"
    apiVer: "v1"
    timeout: "30s"
  sync:
    schema: "https"
    host: "127.0.0.1:9005"
    apiVer: "v1"
    timeout: "30s"

template:
  path: "scripts/templates/"

logger:
  filename: logs/run.log
  level: debug
  encoding: json
  compress: false
  maxAge: 15
  maxSize: 5
  maxBackups: 15
  encodeTime: 2006-01-02 15:04:05.555
  enableKafka: false
```

## 使用
- 初始化节点数据，爬取节点ssl证书到本地

```shell
bsctl init
```

- 启动模拟器

```shell
bsctl run
```

- 清理mock数据(删除init时创建的节点)

```shell
bsctl clean
```
