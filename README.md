# Pholcus [![GitHub release](https://img.shields.io/github/release/hegeng1212/pholcus.svg?style=flat-square)](https://github.com/hegeng1212/pholcus/releases) [![report card](https://goreportcard.com/badge/github.com/hegeng1212/pholcus?style=flat-square)](http://goreportcard.com/report/hegeng1212/pholcus) [![github issues](https://img.shields.io/github/issues/hegeng1212/pholcus.svg?style=flat-square)](https://github.com/hegeng1212/pholcus/issues?q=is%3Aopen+is%3Aissue) [![github closed issues](https://img.shields.io/github/issues-closed-raw/hegeng1212/pholcus.svg?style=flat-square)](https://github.com/hegeng1212/pholcus/issues?q=is%3Aissue+is%3Aclosed) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/hegeng1212/pholcus)

Pholcus（幽灵蛛）是一款纯 Go 语言编写的支持分布式的高并发爬虫软件，仅用于编程学习与研究。

它支持单机、服务端、客户端三种运行模式，拥有Web、GUI、命令行三种操作界面；规则简单灵活、批量任务并发、输出方式丰富（mysql/mongodb/kafka/csv/excel等）；另外它还支持横纵向两种抓取模式，支持模拟登录和任务暂停、取消等一系列高级功能。

![image](https://github.com/hegeng1212/pholcus/raw/master/doc/icon.png)

## 免责声明

**本软件仅用于学术研究，使用者需遵守其所在地的相关法律法规，请勿用于非法用途！！
如在中国大陆频频爆出爬虫开发者涉诉与违规的[新闻](https://github.com/HiddenStrawberry/Crawler_Illegal_Cases_In_China)。<br>
<b>郑重声明</b>：因违法违规使用造成的一切后果，使用者自行承担！！**

## 爬虫原理

![image](https://github.com/hegeng1212/pholcus/raw/master/doc/module.png)

&nbsp;

![image](https://github.com/hegeng1212/pholcus/raw/master/doc/project.png)

&nbsp;

![image](https://github.com/hegeng1212/pholcus/raw/master/doc/distribute.png)


## 框架特点

- 为具备一定Go或JS编程基础的用户提供只需关注规则定制、功能完备的重量级爬虫工具；
- 支持单机、服务端、客户端三种运行模式；
- GUI(Windows)、Web、Cmd 三种操作界面，可通过参数控制打开方式；
- 支持状态控制，如暂停、恢复、停止等；
- 可控制采集量；
- 可控制并发协程数；
- 支持多采集任务并发执行；
- 支持代理IP列表，可控制更换频率；
- 支持采集过程随机停歇，模拟人工行为；
- 根据规则需求，提供自定义配置输入接口
- 有mysql、mongodb、kafka、csv、excel、原文件下载共五种输出方式；
- 支持分批输出，且每批数量可控；
- 支持静态Go和动态JS两种采集规则，支持横纵向两种抓取模式，且有大量Demo；
- 持久化成功记录，便于自动去重；
- 序列化失败请求，支持反序列化自动重载处理；
- 采用surfer高并发下载器，支持 GET/POST/HEAD 方法及 http/https 协议，同时支持固定UserAgent自动保存cookie与随机大量UserAgent禁用cookie两种模式，高度模拟浏览器行为，可实现模拟登录等功能；
- 服务器/客户端模式采用Teleport高并发SocketAPI框架，全双工长连接通信，内部数据传输格式为JSON。

## 下载安装

```
go get -u -v github.com/hegeng1212/pholcus
```

## 创建项目

```
package main

import (
    "github.com/hegeng1212/pholcus/exec"
    // _ "pholcus_lib_pte" // 同样你也可以自由添加自己的规则库
)

func main() {
    // 设置运行时默认操作界面，并开始运行
    // 运行软件前，可设置 -a_ui 参数为"web"、"gui"或"cmd"，指定本次运行的操作界面
    // 其中"gui"仅支持Windows系统
    exec.DefaultRun("web")
}
```
&nbsp;

## 编译运行
正常编译方法
```
cd {{replace your gopath}}/src/github.com/hegeng1212/pholcus
go install 或者 go build
```
Windows下隐藏cmd窗口的编译方法
```
cd {{replace your gopath}}/src/github.com/hegeng1212/pholcus
go install -ldflags="-H=windowsgui -linkmode=internal" 或者 go build -ldflags="-H=windowsgui -linkmode=internal"
```
查看可选参数: 
```
pholcus -h
```
![image](https://github.com/hegeng1212/pholcus/raw/master/doc/help.jpg)

&nbsp;

> *<font size="2">Web版操作界面截图如下：*

![image](https://github.com/hegeng1212/pholcus/raw/master/doc/webshow_1.png)

&nbsp;

> *<font size="2">GUI版操作界面之模式选择界面截图如下*

![image](https://github.com/hegeng1212/pholcus/raw/master/doc/guishow_0.jpg)

&nbsp;

> *<font size="2">Cmd版运行参数设置示例如下*

```
$ pholcus -_ui=cmd -a_mode=0 -c_spider=3,8 -a_outtype=csv -a_thread=20 -a_dockercap=5000 -a_pause=300
-a_proxyminute=0 -a_keyins="<pholcus><golang>" -a_limit=10 -a_success=true -a_failure=true
```

&nbsp;

*注意：*Mac下如使用代理IP功能，请务必获取root用户权限，否则无法通过`ping`获取可以代理！

&nbsp;

## 运行时目录文件

```
├─pholcus 软件
│
├─pholcus_pkg 运行时文件目录
│  ├─config.ini 配置文件
│  │
│  ├─proxy.lib 代理IP列表文件
│  │
│  ├─spiders 动态规则目录
│  │  └─xxx.pholcus.html 动态规则文件
│  │
│  ├─phantomjs 程序文件
│  │
│  ├─text_out 文本数据文件输出目录
│  │
│  ├─file_out 文件结果输出目录
│  │
│  ├─logs 日志目录
│  │
│  ├─history 历史记录目录
│  │
└─└─cache 临时缓存目录
```

&nbsp;

## 动态规则示例

特点：动态加载规则，无需重新编译软件，书写简单，添加自由，适用于轻量级的采集项目。
<br/>
xxx.pholcus.html
```
<Spider>
    <Name>HTML动态规则示例</Name>
    <Description>HTML动态规则示例 [Auto Page] [http://xxx.xxx.xxx]</Description>
    <Pausetime>300</Pausetime>
    <EnableLimit>false</EnableLimit>
    <EnableCookie>true</EnableCookie>
    <EnableKeyin>false</EnableKeyin>
    <NotDefaultField>false</NotDefaultField>
    <Namespace>
        <Script></Script>
    </Namespace>
    <SubNamespace>
        <Script></Script>
    </SubNamespace>
    <Root>
        <Script param="ctx">
        console.log("Root");
        ctx.JsAddQueue({
            Url: "http://xxx.xxx.xxx",
            Rule: "登录页"
        });
        </Script>
    </Root>
    <Rule name="登录页">
        <AidFunc>
            <Script param="ctx,aid">
            </Script>
        </AidFunc>
        <ParseFunc>
            <Script param="ctx">
            console.log(ctx.GetRuleName());
            ctx.JsAddQueue({
                Url: "http://xxx.xxx.xxx",
                Rule: "登录后",
                Method: "POST",
                PostData: "username=44444444@qq.com&amp;password=44444444&amp;login_btn=login_btn&amp;submit=login_btn"
            });
            </Script>
        </ParseFunc>
    </Rule>
    <Rule name="登录后">
        <ParseFunc>
            <Script param="ctx">
            console.log(ctx.GetRuleName());
            ctx.Output({
                "全部": ctx.GetText()
            });
            ctx.JsAddQueue({
                Url: "http://accounts.xxx.xxx/member",
                Rule: "个人中心",
                Header: {
                    "Referer": [ctx.GetUrl()]
                }
            });
            </Script>
        </ParseFunc>
    </Rule>
    <Rule name="个人中心">
        <ParseFunc>
            <Script param="ctx">
            console.log("个人中心: " + ctx.GetRuleName());
            ctx.Output({
                "全部": ctx.GetText()
            });
            </Script>
        </ParseFunc>
    </Rule>
</Spider>
```

## 静态规则示例

特点：随软件一同编译，定制性更强，效率更高，适用于重量级的采集项目。
<br/>
xxx.go

```
func init() {
    Spider{
        Name:        "静态规则示例",
        Description: "静态规则示例 [Auto Page] [http://xxx.xxx.xxx]",
        // Pausetime: 300,
        // Limit:   LIMIT,
        // Keyin:   KEYIN,
        EnableCookie:    true,
        NotDefaultField: false,
        Namespace:       nil,
        SubNamespace:    nil,
        RuleTree: &RuleTree{
            Root: func(ctx *Context) {
                ctx.AddQueue(&request.Request{Url: "http://xxx.xxx.xxx", Rule: "登录页"})
            },
            Trunk: map[string]*Rule{
                "登录页": {
                    ParseFunc: func(ctx *Context) {
                        ctx.AddQueue(&request.Request{
                            Url:      "http://xxx.xxx.xxx",
                            Rule:     "登录后",
                            Method:   "POST",
                            PostData: "username=123456@qq.com&password=123456&login_btn=login_btn&submit=login_btn",
                        })
                    },
                },
                "登录后": {
                    ParseFunc: func(ctx *Context) {
                        ctx.Output(map[string]interface{}{
                            "全部": ctx.GetText(),
                        })
                        ctx.AddQueue(&request.Request{
                            Url:    "http://accounts.xxx.xxx/member",
                            Rule:   "个人中心",
                            Header: http.Header{"Referer": []string{ctx.GetUrl()}},
                        })
                    },
                },
                "个人中心": {
                    ParseFunc: func(ctx *Context) {
                        ctx.Output(map[string]interface{}{
                            "全部": ctx.GetText(),
                        })
                    },
                },
            },
        },
    }.Register()
}
```

&nbsp;

## 代理IP

- 代理IP写在`/pholcus_pkg/proxy.lib`文件，格式如下，一行一个IP：

```
http://183.141.168.95:3128
https://60.13.146.92:8088
http://59.59.4.22:8090
https://180.119.78.78:8090
https://222.178.56.73:8118
http://115.228.57.254:3128
http://49.84.106.160:9000
```

- 在操作界面选择“代理IP更换频率”或命令行设置`-a_proxyminute`参数，进行使用

- *注意：*Mac下如使用代理IP功能，请务必获取root用户权限，否则无法通过ping获取可以代理！

&nbsp;

## FAQ

请求队列中，重复的URL是否会自动去重？
```
url默认情况下是去重的，但是可以通过设置Request.Reloadable=true忽略重复。
```

URL指向的页面内容若有更新，框架是否有判断的机制？
```
url页面内容的更新，框架无法直接支持判断，但是用户可以自己在规则中自定义支持。
```

请求成功是依据web头的状态码判断？
```
不是判断状态，而是判断服务器有无响应流返回。即，404页面同样属于成功。
```

请求失败后的重新请求机制？
```
每个url尝试下载指定次数之后，若依然失败，则将该请求追加到一个类似defer性质的特殊队列中。  
在当前任务正常结束后，将自动添加至下载队列，再次进行下载。如果依然有没下载成功的，则保存至失败历史记录。  
当下次执行该条爬虫规则时，可通过选择继承历史失败记录，把这些失败请求自动加入defer性质的特殊队列……（后面是重复步骤）
```
