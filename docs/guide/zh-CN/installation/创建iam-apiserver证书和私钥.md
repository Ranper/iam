# 创建iam-apiserver证书和私钥

## 创建 iam-apiserver 证书和私钥

创建证书签名请求：

``` bash
$ cd $HOME/user00/work
$ source $HOME/user00/work/environment.sh
$ cat > iam-csr.json <<EOF
{
  "CN": "iam-apiserver",
  "hosts": [
    "127.0.0.1",
    "${IAM_APISERVER_HOST}"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "iam",
      "OU": "user00"
    }
  ]
}
EOF
```

+ hosts 字段指定授权使用该证书的 **IP 和域名列表**，这里列出了 iam-apiserver 节点 IP.

生成证书和私钥：

``` bash
$ cfssl gencert -ca=$HOME/user00/work/ca.pem \
  -ca-key=$HOME/user00/work/ca-key.pem \
  -config=$HOME/user00/work/ca-config.json \
  -profile=iam iam-csr.json | cfssljson -bare iam 
$ ls iam*pem
iam-key.pem  iam.pem
```

1. cfssl gencert: 调用 cfssl 工具生成证书

2. -ca: 指定 CA 证书文件路径

3. -ca-key: 指定 CA 私钥文件路径

4. -config: 指定证书配置文件路径

5. -profile=iam: 使用配置文件中名为 "iam" 的配置模板

6. iam-csr.json: 证书签名请求(CSR)文件

7. | cfssljson -bare iam: 将输出通过管道传递给 cfssljson 工具，生成的文件前缀为 "iam"

执行后会生成两个文件：

iam.pem: 证书文件
iam-key.pem: 私钥文件
这个命令的作用是基于已有的 CA 证书和配置，为 iam-apiserver 服务生成新的 TLS 证书和私钥，用于服务间的安全通信。