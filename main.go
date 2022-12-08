package main

import (
	"log"

	"ent-atlas-migration/ent"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 連線設定
	client, err := ent.Open("mysql", "root:password@tcp(localhost:3306)/ent_atlas_migration?parseTime=True")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	// 自動執行 migration 工具
	// if err := client.Schema.Create(context.Background()); err != nil {
	// 	log.Fatalf("failed creating schema resources: %v", err)
	// }
}
