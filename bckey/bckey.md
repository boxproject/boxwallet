### 概念
 bckey维护的是hd账户，通过master key衍生出子钥，具体衍生方式参见与keyuitl.go

### 目录介绍
1. distribute  各币种hdpubkey to address的实现
2. address.go   ./distribute目录下各币种toAddress的工厂路由
3. prvkey.go 是一个简单的初始化master key的概念
4. key.go     key的具体方法实现
5. keytuil.go  key的db管理与衍生



### 如何拓展
恭喜你，这里如果要接其他币种，理论上只要实现./distribute下的toAddress方法就可以了，然后在address.go的工厂方法中添加新实现的方法



### 其他项目如何直接复用bckey
这是一个幸运的拓展项，这里的bccoin是完全解耦的，以下是bccoin的必须启动函数，参数db使用的是util工具包的数据库接口定义，prefix前缀动态传入
```
//实例化
func InitKeyUtil(db util.Database, pfk_key, pfk_key_count []byte, net bccore.Net) *KeyUtil {
	if keyUtil != nil {
		return keyUtil
	}
	keyUtil = &KeyUtil{db: db, Pfk_Key: pfk_key, Pfk_key_Count: pfk_key_count, net: net}
	initAddressMemCache(net) //初始化地址
	return keyUtil
}
```