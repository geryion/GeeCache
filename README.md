# GeeCache

1、LRU缓存淘汰机制

2、单机并发缓存

3、HTTP服务端

4、一致性哈希

5、分布式节点

6、防止缓存击穿

7、使用Protobuf

ubuntu 安装搭建Protobuf环境

1）下载并解压protobuf安装包到/usr/local下

wget https://github.com/protocolbuffers/protobuf/releases/download/v3.11.2/protoc-3.11.2-linux-x86_64.zip

sudo 7z x protoc-3.11.2-linux-x86_64.zip -o/usr/local

2）将/usr/local/bin加入到PATH中

PS：protoc --version (能够看到当前protobuf版本表示安装正确)

3）下载转换为golang的工具包

go get -u github.com/golang/protobuf/protoc-gen-go

4）将$GOPATH/bin加入到PATH中

5）创建本地简单的示例

6）在当前目录下执行protoc --go_out=.*.proto

7）可以在当前目录下看到新生成的.go文件
