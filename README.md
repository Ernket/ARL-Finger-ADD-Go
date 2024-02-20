# 前言
在github上看到某个师傅，这种脚本都要单独拿来卖钱属实是没想到，刚好最近在写新的自动化框架，拿go练练手，所以这个项目只是写来先“测试”用的<br>
有能力的师傅们也可以自己修修改改，毕竟25也能吃顿好的了，不至于浪费钱
![](https://github.com/Ernket/ARL-Finger-ADD-Go/blob/48087cc2de0d65fa72e6a2d81beeeed329140f66/png/1.png)

# 环境
- go版本: 1.21.6
- ARL版本: 2.5.5

# 用法
加了个删除所有指纹的操作，方便重置<br>
`finger.json`文件和执行程序在同目录即可
```
Usage: main -url="https://x.x.x.x" -username="xxx" -password="xxxx" [-thread=10|-n="true"]
选项:
  -n	删除所有指纹
  -password string
    	密码
  -thread int
    	线程数 (default 10)
  -url string
    	URL地址
  -username string
    	用户名

```

# 结果
我在自己搭建的arl中运行，结果是`7685`条<br>
![](https://github.com/Ernket/ARL-Finger-ADD-Go/blob/48087cc2de0d65fa72e6a2d81beeeed329140f66/png/2.png)

# 参考项目
https://github.com/Funsiooo/chunsou  (finger.json文件)<br>
https://github.com/loecho-sec/ARL-Finger-ADD
