# daily report

## 10.1,10.2

10.1

- 下载WSL，Ubuntu
- 体验dig
- 学习commit message 规范

10.2
上午

- 发现WSL遇到网络问题，尝试了更换镜像，更改clash设置打开服务模式，解决了下载中的网络问题，但在dig使用中再次出现网络问题，通过配置代理解决（网络问题此前也遇到过，都花费很长时间找问题,这次试了两个小时，吸取经验）
- WSL配置vscode  
有用的参考：<https://learn.microsoft.com/zh-cn/windows/wsl/tutorials/wsl-vscode>

下午

- 学习了在项目中使用CMake（后面发现可以用go）
- 学习DNS,试图理解项目要做些什么,一些有用的参考  

    >什么是DNS <https://aws.amazon.com/cn/route53/what-is-dns/>  
    互联网协议入门（一） <https://www.ruanyifeng.com/blog/2012/05/internet_protocol_suite_part_i.html>  
    互联网协议入门（二） <https://www.ruanyifeng.com/blog/2012/06/internet_protocol_suite_part_ii.html>  
    DNS 原理入门 <https://www.ruanyifeng.com/blog/2016/06/dns.html>  
    DNS 查询原理详解 <https://www.ruanyifeng.com/blog/2022/08/dns-query.html>

晚上

- 发现可以用go，发现再次遇到了建立不了工作区的问题，试图解决，好像发现了原因，但还没
- 写记录的时候遇到Error loading webview: Error: Could not register service worker问题  
    >有用的参考：<https://stackoverflow.com/questions/67698176/error-loading-webview-error-could-not-register-service-workers-typeerror-fai>

10.3

- 一次完整的http请求过程 <https://cloud.tencent.com/developer/article/1500463>
- 实现简单的服务器和客户端，完成相互测试
- 实现通过命令行控制，未成功，调试中