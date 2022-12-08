# 學習 Ent and Atlas 的 migration 操作

## 初始化

```sh
git init -b main

go mod init ent-atlas-migration
```

## 建立 User Model

```sh
# 建立User指令
go run -mod=mod entgo.io/ent/cmd/ent init User
```

```go
// ent/schema/user.go
// 新增欄位
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name").
            Default("unknown"),
        field.Int("age").
            Positive(),
    }
}
// ...
```

```sh
# 生成代碼
go generate ./ent
```

## 資料庫連線設定

```sh
# 安裝 mysql driver
go mod download github.com/go-sql-driver/mysql
```

```go
// main.go
// 主程式
package main

import (
	"context"
	"log"

	"ent-atlas-migration/ent"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 連線設定
	client, err := ent.Open("mysql", "root:password@tcp(127.0.0.1:3306)/ent_atlas_migration?parseTime=True")
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	defer client.Close()

	// 自動執行 migration 工具
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
```

```sh
# 套件依賴處理
go mod tidy
# 執行後,觀察資料庫表和欄位
go run main.go
```
