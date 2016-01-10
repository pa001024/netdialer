Netdialer
=========

The dialer for the fucking shanxun network in Zhejiang, written in go. 

### for developer

get this by following command

> go get -v -u github.com/pa001024/netdialer/netdialer 

简介
----

本程序分GUI版本和CLI(命令行)版本

以下是使用说明

GUI版本使用说明
---------------

（有BUG请在 issues 里说）

因为是 GUI 其实没啥可说的。

就说几个 tips: （亲，看完再下载）

### 关于模式

内建几种路由器适配。（目前还没有 TP 真是抱歉 233）

__注意：要使用路由模式，你需要先将网线插到路由器WAN口 然后进路由器后台
将WAN口设置为 DHCP 模式（也可能叫动态IP）。__

然后打开软件，选择对应的模式并且改好路由器的地址密码（注意密码不是 WiFi 密码而是后台管理密码）等参数之后就可以直接用了。

如果没有自己的路由器的模式，那么就去路由界面自行获取IP直接填到模式内即可。


命令行版本使用说明
----------------------

参数跟GUI基本一样

不过因为可以命令行使用，所以方便到任何平台运行，也可设置开机启动、计划任务等等，总之想怎么玩怎么玩。

#### 开始连接

日常使用方法如下：

如果是本地拨号：
> netdialer -u xx -p xx

如果是路由器：
> netdialer -ip hiwifi -ra xx -ru xx -rp xx -u xx -p xx

#### 断开连接

ps.断开连接可以不要闪讯账号密码。

如果是本地拨号：
> netdialer -d

如果是路由器：
> netdialer -ip hiwifi -ra xx -ru xx -rp xx


#### 支持自己开发路由器适配

只要`-ip stdin`并用管道传入WAN口IP即可。

> 比如 `yourrouter -u xxx -p xxxx | netdialer -u xxx -p xxx -ip stdin`

#### 依然可以使用内建的几种路由支持：

> 比如 `-ip hiwifi` 然后填写`-ra [路由地址] -ru [路由用户名] -rp [路由密码]`参数即可

#### 关于路由模式2

（也就是 `-r` 参数）

这个是通过WAN口直接拨号的，本程序没有提供相应的适配器（除了TP），因为生成的用户名每五秒就会变，所以实用性不高。

如果你需要使用，可以自己编写适配器，可拨到本地或路由。

只需开启 `-r` 模式，程序就会通过stdout输出两行文本。第一行是经过URL编码的用户名，第二行是密码。只需要将PPPoE的账号密码设置成这个然后立刻拨号即可（一定要快，因为有效期只有五秒）。

##### 目前只有一个TP的拨号程序，使用方法如下：

> netdialer -u xx -p xx -r | tplink -p xx -a 192.168.1.1



下载
----

度云 [http://pan.baidu.com/s/1dDQiXpf](http://pan.baidu.com/s/1dDQiXpf)
