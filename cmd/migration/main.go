package main

/*
現在は、go-migrateを使っているので、このファイル(gormのAuto migrate)は不要
*/

// import (
// 	"github.com/fuku01/test-v2-api/config"
// 	"github.com/fuku01/test-v2-api/pkg/domain/model"
// )

// func main() {

// 	db, err := config.NewDatabase()
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = db.AutoMigrate(&model.User{}, &model.Todo{})
// 	if err != nil {
// 		panic(err)
// 	}
// }
