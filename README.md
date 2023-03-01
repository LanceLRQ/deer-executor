<h1 align="center">Deer Executor</h1>
<p align="center">一个基于Go语言实现的代码评测工具</p>

![自动构建](https://github.com/LanceLRQ/deer-executor/workflows/Go/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/LanceLRQ/deer-executor)](https://goreportcard.com/report/github.com/LanceLRQ/deer-executor)

## ✨ 特性
 - 以CLI方式运行，不需要OJ平台；
 - 支持多种编程语言程序的判定，你可以自行扩展更多的语言；
 - 支持将题目配置和数据打包，随处都能一行命令运行评测；
 - 支持评测结果压缩打包，便于OJ存档和回放；
 - 支持使用Testlib作为出题工具；
 - 基于Linux和Mac OS平台，如果你感兴趣可以帮我实现Windows平台的代码;-)。
 
## 🔨 开发进度

 - ✅ 接入CLI
 - ✅ 多语言支持
 - ✅ 多评测方式
 - ✅ 完整评测流程支持
 - ✅ 评测结果打包并签名
 - ✅ 题目配置和数据打包功能
 - ✅ 兼容WeJudge3.x的数据结构
 - ✅ Testlib 支持
 - ✅ 编写文档
 - ✅ 评测日志
 - ✅ 题目数据包支持ZIP
 - 🔲 Windows评测支持
 - 🔲 Deer-Executor GUI
 - 🔲 安全沙箱
 
## 📦 文档

文档托管在Github Wiki上，[点击访问](https://github.com/LanceLRQ/deer-executor/wiki)

## 待解决

go>=1.17，不能正常进行Fork，原因未知...请使用go1.16编译。

## 🤝 鸣谢

感谢开源项目[Loco's runner](https://github.com/dojiong/Lo-runner) 为本程序提供了黑盒评测的实现思路。

感谢我的同学Wolf Zheng和Tosh Qiu提出的交互式评测的需求和基本工作流程的描述。

感谢以下组织对WeJudge项目的支持（排名不分先后）：

* [北京师范大学(珠海校区)](http://www.bnuz.edu.cn)

* [珠海市计算机学会](http://www.zhcomputer.org.cn/)

* [信息技术学院](http://itc.bnuz.edu.cn)

* [北师珠ACM协会](http://acm.bnuz.edu.cn)

感谢WeJudge团队每一位成员对项目的支持和付出！

感谢以下博客、开源项目等为本项目提供参考学习的资料。（不分顺序）

* [JanBox的小站](https://boxjan.com/)
* [VOJ](https://github.com/hzxie/voj)
* [QDUOJ-Judger](https://github.com/QingdaoU/Judger)

等等


## 🔗 相关链接

📃 我的博客：[https://www.lanrongqi.com](https://www.lanrongqi.com)

《从零开始的代码评测系统设计与实践》序列

1. [序](https://www.lanrongqi.com/2020/07/online-judge-development-0/)
2. [判题机篇-进程和输入输出](https://www.lanrongqi.com/2020/07/online-judge-development-1/)
3. [判题机篇-资源占用与限制](https://www.lanrongqi.com/2020/08/online-judge-development-2/)
4. [判题机篇-运行结果处理](https://www.lanrongqi.com/2020/08/online-judge-development-3/)
4. [判题机篇-特殊评测](https://www.lanrongqi.com/2020/08/online-judge-development-4/)

🖥️ WeJudge：

[https://oj.bnuz.edu.cn](https://oj.bnuz.edu.cn)

**本项目基于GPLv3协议开源，欢迎各位开发者以非商业目的使用和开发本程序，使用过程中如果遇到什么问题，请发Issue一起讨论哦！**

**如果你正在使用本判题机的开发OJ网站，欢迎通过ISSUE告知，我会将链接挂在这里哦！**
