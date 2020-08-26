# goleaf_games_server

## 一，说明
#### 根据游戏服go leaf框架 + 客户端node + vue 写的多个小游戏集合， 现有猜拳，井字棋
#### 线上地址 http://139.159.155.172:8080/#/
#### 客户端git地址:https://github.com/lgbya/goleaf_games_client

## 二，配置环境
**linux + golang + mysql**

## 三，安装步骤

1. **创建配置文件**
```
mv bin/conf/server.json.example bin/conf/server.json
```

2. **创建数据库**
	mysql表放在 bin/db.sql

3. **修改server.json**
```
	{
		"LogLevel": "debug" //日志等级,
		"LogPath": "",//日志目录
		"WSAddr": "xxx.xxx.xxx.xxx:xxx"//websocket 如127.0.0.1:3666,
		"MaxConnNum": 20000,
		"Md5Key" : "gjwagawrfhwealu3131", //token和用户密码加密key
		"SqlSrv" : {
			"Username" : "root",    //mysql 用户名
			"Password" : "root",    //mysql 密码
			"Host" : "127.0.0.1",   //mysql 连接地址
			"Database": "puzzle_games",//mysql 对应的游戏库
			"MaxIdle": 15,        
			"MaxOpen": 10
		}
	}
```

4. **执行
以下命令生成server二进制文件**
```
go get -v 
go install server
```

5. **运行bin目录下的server**
```
./server
```





