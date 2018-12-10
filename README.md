### wallet世界观
以下是wallet项目的介绍，在各个package下有进一步的md文档介绍

### 币种测试环境客户端安装
1. btc  => bitcoin-core 修改配置文件
2. omni => omni-core 修改配置文件（和btc公用一个目录，所以，如果用omni和btc公用omni-core代码需要微调，目前分开）
3. ltc   => litecoin-core
4. eth   => eth私链

### 目录结构

```
+-- bccoin
|   +-- distribute    //coin信息初始化json文件
|   +-- coin.go       //币种实例化，算术运算等主要操作，未进缓存的不能用
|   +-- coin_cache.go //币种缓存信息，新增币种会代币会把代币的单位，精度等内容缓存起来
|
+-- bcconfig  //配置文件相关，代币配置各主链独立
| 
+-- bccore    //全局唯一的常量(重点：币种的全局，定律的那种)
|
+-- bctrans  
|   +-- client       //一个主链一个客户端
|   +-- clientseries //客户端系别，可以理解为币种的客户端变量，client需要从这里实例化
|   +-- token        //代币token添加的地方，添加后会缓存起来
|   +-- trans.go     //逻辑总入口，第三方对接调用这里就够了
|
|
+-- signature       //离线签名工具包，这里为纯逻辑计算
|
+-- daemon   //后台程序，各币种爬块，目前使用到client的部分直接调用了"bctrans/client"中的实力，需要注意的是，这里之后解耦的话改为引用此package后自己实例化就好了，千万要有这个概念，因为目前只要一个进程启动，所以这样做
|
+-- db   //数据库相关
|
+-- errors //全局error
|
+-- mock  //单元测试模拟数据
|
+-- pipeline //用于解耦daemon和bctrans的通信通道，以后业务扩大，需要分离的话，这里替换为队列的概念就行了
|
+-- util  //工具包，纯工具
```


### 设计概念
1. 币种概念有2个枚举
    - enum int 
    >举例：BTC:1,ETH:2, ERC2=ETH=2；用于账户公用等抽象

    - enum string 
    >举例： BTC:BTC,ETH=ETH,ERC20=ERC20;这里用来区分各种币种，ERC20也抽象为了一种币种，代币间具体区分还会加上一个合约token；之后有其他的主链的合约币也是如此设计
2. 多币种的抽象中，有主链币，代币，寄生币，fork币等等，所以总结后以代币客户端类型，单例客户端2中概念
    - 客户端类型：如比特币类型，以太坊类型，比特币类型上可以跑比特币主链，USDT
    - 客户端单例：从客户端类型中实例化，这里就把eth和erc20分为2个单例
    
3. 发起交易等操作与后台扫块解耦，以pipeline通信，方便以后架构升级
```
graph TD
A[交易部分]
A-->B
B-->A
B[队列通信]
B-->C
C-->B
C[后台爬块]
```
>概念上可大可小，小的时候队列通信就是一个channel，大了可以是一个mq工具，后台爬块小了可以是一个线程，大了可以是多个反代的微服务



4. 单元测试，所有理念的设计必须在编写完成的同时附带上添加逻辑的整个单元测试，一个方便后期维护，同时可以鉴别设计的可行性


### 如何在boxwallet上继续多币种的开发
#### 1. bccoin 目录
>1. 如果是主链币，需要在 ./dicstribute/coin_info.json中添加基本信息,参考以下json结构

```
[
  {
    "coinType": 1,
    "token": "",
    "symbol": "BTC",
    "decimals": 8,
    "name": "比特币"
  },
  {
    "coinType": 2,
    "token": "",
    "symbol": "ETH",
    "decimals": 18,
    "name": "以太坊"
  },
  {
    "coinType": 3,
    "token": "",
    "symbol": "LTC",
    "decimals": 18,
    "name": "莱特币"
  }
]

```
> 如果是eth的erc20代币，因为币种的精度是自定义的，这里在订阅token的同时会自动添加到kv数据中，所以不需要在bccoin目录改任何东西，如果是其他的币种的代币，参照erc20的写法


#### 2. bckey 目录

>1. boxwallet全局采用[hd钱包的概念](http://book.8btc.com/books/1/master_bitcoin/_book/4/4.html)，所有币种主钥共用一把，子钥按照深度衍生，可定制结构如下


```
m/bccore.BloclChainType(币种类比)/自定义层级([]uint32)/num(uint32)

1：m:"主钥"

2：bccore.BloclChainType：wallet中的枚举，如：btc这里设定的是1，这个是代码定死的

3：自定义层级：这个是一个衍生深度定制，可以有业务端定义后传进来，是一个uint32的数组

4：num:最后一次衍生的下标，这里也是系统自己计算的，如果当前深度是 m/1/2/3/4，那么下一个衍生出来的只能是m/1/2/3/5

```

> 2. bckey/distribute 包下是不同币种的addres转换方法，如果需要添加其他主链币需要在这里加该币种的address转换方法，最后在bckey/address.go中工厂模式添加


#### 3. bctrans 目录

> 1. 该目录下以此有4个子目录，分别为：clientseries、client、token、txutil

- clientseries:客户端系列目录，抽象了客户端的种类，如果有币种使用的同一套源码开发，如btc,usdt,ltc或者eth,erc20(这里把erc20也理解成了一个币种)，所以这里就需要btc系和eth系的客户端类别，如果要加xlm币种则需要再加一个系别，以太坊经典的话应该要使用geth fork出来的库了，或许不能直接用eth的类别了
- client: 客户端实例化，原始目录中有btc、eth、erc20，可以看见，eth和erc20都是从clientseries中的eth_series中实例化出来的
- token: 代币目录，这里主要是功能是添加代币，token_law.go中定义了GetTokenInfo的函数接口，boxwallet中，如果要对代币进行算术运算，必须先得调用此方法，此方法会把相关代币的币种信息缓存到bccoin中，没订阅的币种在boxwallet中计算是不能成功的；另外对比token地址也做了缓存，在token_cache.go中可以看到，具体作用是，爬块的时候对于代币的解析会用到地址的概念（要注意的是：如果后台爬块和client解耦后，这块依赖也需要解耦，这是很容易的，所以该文件要保持洁净）
- txutil: 这里是一个tx的工具文件，主要用于tx的离线签名，所以这里需要保持不引用客户端或者数据库等其它介质，必须单纯依靠类库和cpu就完成运算，这也是boxwallet的难点之一
- trans 这里是对client的工厂化封装，这里有aop切面化编程的概念，对于全局的交易数据监控拦截等可以在这里操作，所以这里包含了对于pipeline包的数据传输（pipeline在后面会有介绍）


#### 4. daemon 目录

> 这里主要是数据爬块，概念上这个是纯后台的东西，daemon.go中定义了爬块的主流程和接入币种需要实现的方法，这里也会使用到pipeline包

daemon接入新币种：

1. 添加daemon_xxx.go文件，参照daemon_btc.go和daemon_eth.go的实现
2. daemon.go/initDaemons() 中添加初始化
3. 地址相关，因为爬块需要检测指定地址，所以这里引用了地址的管理包，目前直接依赖于bckey/address_cache.go,可以看到，如果要添加新的币种，需要在这里初始化的时候加载相关币种，所有的初始化依赖于之前的kv数据库的存储

#### 5. pipeline 目录
> 这里只有一个文件，作用是解耦client和daemon之间的通信，是业务更加流程化后便于后期的拓展，这在最上面介绍过为什么如此做，这里不再介绍了

#### 6. mock 目录
> 这里是本地单元测试的mock数据，同时也包含了如何启动boxwallet的概念，以及全局化了各个package的单元测试，变相的实现了全局测试，可以更好的优化项目结构


### pipeline 结构相关
概念图见document文件夹，后期放到github后替换为图片预览

### 注意事项
1. 初始化第一步，必须先SaveMasterKey,所有的地址衍生都依赖于masterPublicKey,这里千万别把privateKey导进来了，因为公私钥衍生都支持，这里一定要注意
2. 所有的地址阻塞只阻塞系统内部衍生的地址
3. 所有的代币的计算必须要getTokenInfo后让token的概念性数据缓存到系统内部后才能增删改查，不然不支持，因为代币多了去了，随便查随便用容易导致各种bug,有规则的使用才能无忧无虑
### 总结

最简单的接入方法
1. 看到目录下有  xxx_btc以及xxx_eth等文件的，直接新建 xxx_新币种，复制黏贴，替换内部方法
2. 跑xxx_test单元测试，只要修改到能够运行
3. 如果不能运行，看一下上面的介绍，继续 1和2步操作，当前如果新加的币种能直接用以上的代码，那就更简单了


