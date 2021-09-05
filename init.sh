grep certificate-authority-data ~/.kube/config | cut -d" " -f 6 | base64 -D > config/credential.crt
openssl x509 -in config/credential.crt -out config/credential.pem
kubectl -n kube-system describe secret default| awk '$1=="token:"{print $2}' > config/token.txt
grep server  ~/.kube/config | cut -d" " -f 6 > config/url.txt
