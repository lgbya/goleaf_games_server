# goleaf_games_server

### 说明
#### **根据go leaf框架写的多个悠闲游戏， 现有猜拳，井字棋**

### 配置环境
**linux + golang + mysql**

### 安装步骤

1. **创建配置文件**
```
mv bin/conf/server.json.example bin/conf/server.json
```
1. **创建数据库**
	mysql表放在 bin/db.sql

1. **修改server.json**
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

1. **执行以下命令生成server二进制文件**
```
go get -v 
go build server
```

1. **运行bin目录下的server**
```
	./server
```





