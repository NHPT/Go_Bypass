# Go_Bypass
可以过国内主流免杀
用法：服务器建立二个文件 key.txt 用来存放 密钥，Exp.txt用来存放AES加密的payload


go build -ldflags "-H windowsgui -w -s" Exp.go
//编译

原理及过程：
https://mp.weixin.qq.com/s/3R0sebWHIiHeDKqmRsoSwg
