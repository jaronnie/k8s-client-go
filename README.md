# k8s-client-go example

k8s version: 1.21.2

提供两种方式连接到 k8s server

* connect.Connect()
* connect.DefaultConnect(), connect from .kube/config, more simple
## config

如果你使用第二种方式进行连接，下面的内容不必了解, 直接看 example 即可

* credential.crt
* credential.pem
* token.txt
* url.txt

使用如下的脚本获取相关的参数, **默认使用命令空间 `default`**

```shell
grep certificate-authority-data ~/.kube/config | cut -d" " -f 6 | base64 -D > config/credential.crt
openssl x509 -in config/credential.crt -out config/credential.pem
kubectl -n kube-system describe secret default| awk '$1=="token:"{print $2}' > config/token.txt
grep server  ~/.kube/config | cut -d" " -f 6 > config/url.txt
```

具体的一些分析请参照 [kubeconfig2credential](docs/kubeconfig2credential.md)

## example
### get

* get pod list
* get deployment list

```shell
cd get
go test -gcflags=-l -v .
```
## TODO
* deploy program
* ...
