# 电报种子提醒小助手

## 简介
---
电报提醒小助手是一个将种子推送到Telegram（电报）上的机器人。

![201909071942226b4d5844c9bf015a502431068b6d1b95.jpg](https://i.loli.net/2019/09/08/lseUHuR8ZYWKMVv.jpg)

每15分钟读取一次种子区页面，找出最新发布的，最新置顶的，以及优惠信息变动的种子，以上图形式发送到电报。可以远程连接到Transmission一键下载，或下载种子到指定目录

由于PT的特殊属性，不建议在同一IP上登录多个账号，以及种子与个人账号相关，所以此机器人只能自己搭建。

## 搭建流程
---
提示：最好有NAS、VPS等可以24小时运行的机器，可能用到虚拟机或docker，程序以Golang编写，并编译成二进制文件运行。
1、首先在电报内通过@BotFather创建一个bot，获取token
2、通过@BotFather设置bot，依次点击Edit Bot--->Edit Commands，添加/get 命令。
3、通过https://api.telegram.org/bot{token}/setWebhook?url={webhook_url} 地址来添加webhook，其中{token}替换成第一步获取的token，{webhook_url}换成自己的url，建议地址为:https://example.com/{token}。
（以上三步可以参考任何电报bot教程，也可在群内找我私聊，昵称Link）
4、将bot绑定的URL转发到9388端口
（家里路由需要有DDNS，或固定的外网IP）
5、修改config.toml,根据自己需求设定
6、运行程序，输入验证口令，输入验证码，输入/get命令开始获取种子。

### 更新信息
目前程序暂定版本为0.0.1beta，暂时支持站点为三个，MoeCat、铂金家、猫站，程序更新将唯一更新到本帖。遇到任何问题或不知如何操作可开启issues。

2019.9.13 更新：更新Pter下载链接，支持了ourbits和hdsky，支持设置cookie，支持CHD，支持QB远程下载，提示信息优化，修复一些bug。
再次提示：目前为beta版本，可能无法正确运行
