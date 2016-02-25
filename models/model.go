package models

// Tables 将生成表的models集
func Tables() []interface{} {
	return []interface{}{
		&OAuth{},
		&LocalAuth{},
		&CodeAuth{},
		&CodeLog{},
		&User{},
		&SigninLog{},
	}
}

// CreateTables 创建表
func CreateTables() {
	db.CreateTable(Tables())
}

// AutoMigrate 自动更新表
func AutoMigrate() {
	for _, value := range Tables() {
		db.AutoMigrate(value)
	}
}
