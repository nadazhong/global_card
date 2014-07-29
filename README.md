global_card
===========

梦想做一款炉石一样的卡牌游戏(或者根本就不像呢) 

day 1 
完成文件框架搭建 测试开各个server是否正常
day 2 
构建自动编译文件 测试package
添加.gitignore 忽略本地编译产生文件
引用github中工程 
git clone https://github.com/jimlawless/cfg.git
day 3 
启动server和client连接测试,关于解包和封装包不在工程中
恶补net包
day 4
看了一个朋友的golang学习过程 自己补习了一遍
git clone https://github.com/SimonWaldherr/golang-examples.git
同时看了看
git clone https://github.com/mxk/go-sqlite.git
day 5
这两天还真热,没有网络没有空调,乡下农村非常热.还好可以静静的写一些代码.
day 6
服务器组为 Gate 负责client通信
GS为玩家活动服务器,可以理解在一个王国的时间.各个王国之间的玩家不能攻击,但是可以允许整体搬迁到另外一个王国.
CENTER为全球事件服务器,存储公会 全球排行榜一类的数据
