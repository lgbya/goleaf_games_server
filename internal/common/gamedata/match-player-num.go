package gamedata

// 文件名必须和此结构体名称相同（大小写敏感）
// 结构体的一个实例映射 recordfile 中的一行
type MatchPlayerNum struct {
	// 将第一列按 int 类型解析
	// "index" 表明在此列上建立唯一索引
	GameId int "index"
	// 将第二列解析为长度为 4 的整型数组
	PlayerNum int
}

var RfMatchNum = readRf(MatchPlayerNum{})

func GetMatchNum(gameId int) int {
	// 按索引查找
	// 获取 Test.txt 中 Id 为 1 的那一行
	r := RfMatchNum.Index(gameId)

	if r != nil {
		return r.(*MatchPlayerNum).PlayerNum
	}
	return 0
}
