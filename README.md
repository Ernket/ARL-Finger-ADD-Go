# 前言
在github上看到某个师傅，这种脚本都要单独拿来卖钱属实是没想到，刚好最近在写新的自动化框架，拿go练练手，所以这个项目只是写来先“测试”用的<br>
有能力的师傅们也可以自己修修改改，毕竟25也能吃顿好的了，不至于浪费钱
![](https://github.com/Ernket/ARL-Finger-ADD-Go/blob/48087cc2de0d65fa72e6a2d81beeeed329140f66/png/1.png)

# 环境
- go版本: 1.21.6
- ARL版本: 2.5.5

# 用法
加了个删除所有指纹的操作，方便重置<br>
`finger.json`文件和执行程序在同目录即可，同时新增了`config.yaml`配置文件，方便执行
添加指纹用到的线程数也放在了配置文件里
```
Usage: main [-d|-a]
选项:
  -a	添加finger.json文件中的指纹
  -d	删除所有指纹
  -s string
    	查询的任务名称

```


# 结果
我在自己搭建的arl中运行，结果是`12568`条<br>
![](https://github.com/Ernket/ARL-Finger-ADD-Go/blob/main/png/2.png)

# 更新记录

- 2024.3.5
<br>将获取地址和账号密码的方式修改成了读取配置文件的方式<br>
- 2024.5.31
<br>1.修改了判断成功的逻辑，之前是判断请求的状态码，发现目前版本并不可行，改为判断json里的code<br>
2.更改了添加指纹的逻辑，keyword存在多个值的时候并不能用,分割的方式来添加，但是可以重名，所以同个名称会出现不同的规则<br>
3.修复了一个[bug](https://github.com/Ernket/ARL-Finger-ADD-Go/issues/2)<br>
4.更新了-s参数，用来导出任务<br>
5.为了避免每次使用脚本会退出登录（不允许重复登录），增加了api_key的方式来请求，当apikey存在的时候，默认先用key，如果为空则使用账号密码登录<br>

# 参考项目
https://github.com/Funsiooo/chunsou  (finger.json文件)<br>
https://github.com/loecho-sec/ARL-Finger-ADD
