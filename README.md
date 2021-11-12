# FISCO_clean_logs
FISCO节点日志定期自动清理工具

编译、构建二进制可执行文件。如果是在windows上编译linux的二进制可执行文件，需先执行

```
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
```

然后执行
```
go build
```
构建好二进制文件后，在部署了FISCO节点的服务器上，执行
```
nohup ./FISCO_clean_logs /home/admin/fisco_mysql/nodes/127.0.0.1 &
```
只有一个参数，比如/home/admin/fisco_mysql/nodes/127.0.0.1，代表nodes所在的目录路径。这个工具会每隔1分钟自动清理所有node的log，只保留最近的5个log文件，其他log文件全部删除。这样就不用再担心FISCO log文件膨胀占用磁盘空间过大了。
