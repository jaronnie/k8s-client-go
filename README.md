# k8s-client-go example

## config

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

具体的一些分析请参照 [kubeconfig2credential](docs/kubeconfig2credential.md)## example

## example
### get

* get pod list
* get deployment list

```shell
cd get
go test -gcflags=-l -v .
```
