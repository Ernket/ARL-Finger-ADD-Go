# 前言
在github上看到某个师傅，这种脚本都要单独拿来卖钱属实是没想到，刚好最近在写新的自动化框架，拿go练练手，所以这个项目只是写来先“测试”用的<br>
有能力的师傅们也可以自己修修改改，毕竟25也能吃顿好的了，不至于浪费钱<br>
（PS：各位师傅在用的时候如果有任何bug或者建议都欢迎反馈）
![](https://github.com/Ernket/ARL-Finger-ADD-Go/blob/48087cc2de0d65fa72e6a2d81beeeed329140f66/png/1.png)

# 环境
- go版本: 1.21.6
- ARL版本: 2.5.5

# 用法
加了个删除所有指纹的操作，方便重置<br>
`finger.json`文件和执行程序在同目录即可，同时新增了`config.yaml`配置文件，方便执行
添加指纹用到的线程数也放在了配置文件里
```
Usage: main [-d|-a|-s]
选项:
  -a	添加finger.json文件中的指纹
  -d	删除所有指纹
  -o    使用旧版本的指纹添加逻辑
  -s string
    	查询的任务名称

```

# 指纹文件优化
在指纹识别的时候发现一个问题，那就是有的指纹误报率极高，发现的逻辑我自己猜想，是因为原先工具会去对keyword里的去组合对比<br>
看到新版的arl中其实是支持&&的，但是为了旧版本的用户使用，决定还是留着old_finger.json，新版则使用finger.json<br>


# 结果
我在自己搭建的arl中运行，结果是`9264`条<br>
![](https://github.com/Ernket/ARL-Finger-ADD-Go/blob/main/png/4.png)

# 更新记录

- 2024.3.5
<br>将获取地址和账号密码的方式修改成了读取配置文件的方式<br>
- 2024.5.31
<br>1.修改了判断成功的逻辑，之前是判断请求的状态码，发现目前版本并不可行，改为判断json里的code<br>
2.更改了添加指纹的逻辑，keyword存在多个值的时候并不能用,分割的方式来添加，但是可以重名，所以同个名称会出现不同的规则<br>
3.修复了一个[bug](https://github.com/Ernket/ARL-Finger-ADD-Go/issues/2)<br>
4.更新了-s参数，用来导出任务<br>
5.为了避免每次使用脚本会退出登录（不允许重复登录），增加了api_key的方式来请求，当apikey存在的时候，默认先用key，如果为空则使用账号密码登录<br>
- 2024.10.12
<br>删除了部分无特征的规则，不然误报率极高
- 2024.11.1
<br>版本大更新，原本大改arl的指纹识别模块，结果惊喜的发现新版本(我用2.6.2)已经支持了&&等一系列符号
```
# 定义操作符
equals = CaselessLiteral("=")
contains = CaselessLiteral("==")
not_contains = CaselessLiteral("!=")
and_op = CaselessLiteral("&&")
or_op = CaselessLiteral("||")
not_op = CaselessLiteral("!")
```
<br>所以直接改指纹添加的逻辑即可，现在会去自动用&&去拼接(考虑到会有匹配到body，但是icon不存在的情况，所以同名规则拼合的写完后还是给我删了)，写一块去匹配，效果如下
![](https://github.com/Ernket/ARL-Finger-ADD-Go/blob/main/png/3.png)
<br>但是方便旧版本的用户也使用，我提供了-o参数，在使用 -a -o的时候，程序会使用旧版本的添加逻辑
<br>就是得麻烦手动改名，把old_finger.json改成finger.json，不然误报率很高

# TODO
增加参数，指定finger.json文件

# 参考项目
https://github.com/Funsiooo/chunsou  (finger.json文件)<br>
https://github.com/loecho-sec/ARL-Finger-ADD
