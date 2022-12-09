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
# 產生代碼
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

# 設定 Atlas

```go
// main.go
// 不要啓用

	// 自動執行 migration 工具
	// if err := client.Schema.Create(context.Background()); err != nil {
	// 	log.Fatalf("failed creating schema resources: %v", err)
	// }
```

- 若是用 entgql 則改用 entc 方式

```diff
# ent/generate.go
# 填加 --feature sql/versioned-migration

- //go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
+ //go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration ./schema
```

```go
// ent/migrate/main.go
// 執行 migration 工具

//go:build ignore

package main

import (
    "context"
    "log"
    "os"

    "<project>/ent/migrate"

    atlas "ariga.io/atlas/sql/migrate"
    "entgo.io/ent/dialect"
    "entgo.io/ent/dialect/sql/schema"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    ctx := context.Background()
    // Create a local migration directory able to understand Atlas migration file format for replay.
    dir, err := atlas.NewLocalDir("ent/migrate/migrations")
    if err != nil {
        log.Fatalf("failed creating atlas migration directory: %v", err)
    }
    // Migrate diff options.
    opts := []schema.MigrateOption{
        schema.WithDir(dir),                         // provide migration directory
        schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
        schema.WithDialect(dialect.MySQL),           // Ent dialect to use
        schema.WithFormatter(atlas.DefaultFormatter),
    }
    if len(os.Args) != 2 {
        log.Fatalln("migration name is required. Use: 'go run -mod=mod ent/migrate/main.go <name>'")
    }
    // Generate migrations using Atlas support for MySQL (note the Ent dialect option passed above).
    err = migrate.NamedDiff(ctx, "mysql://root:password@localhost:3306/test", os.Args[1], opts...)
    if err != nil {
        log.Fatalf("failed generating migration file: %v", err)
    }
}
```

```sh
# 建立資料夾
mkdir ent/migrate/migrations

# 為 ent/generate.go 填加的 --feature sql/versioned-migration 產生代碼
go generate ./ent

# 會檢查 ent 的 schema, 必需有做 go generate ./ent
# 產生 *_create_users.sql, atlas.sum 來記錄 migration
go run -mod=mod ent/migrate/main.go create_users
```

# lint 代碼

```diff
# ent/schema/user.go
# 填加身高欄位

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
		field.Int("age").
			Positive(),
+		field.Float("height").
+			Positive(),
	}
}
```

```sh
# 產生 height 欄位的代碼
go generate ./ent

# 產生 migrate 的 *_users_add_height.sql 記錄
go run -mod=mod ent/migrate/main.go users_add_height

# *.sql 的執行檢查
# --latest 對最新的N遷移文件運行分析
go run -mod=mod ariga.io/atlas/cmd/atlas@master migrate lint \
  --dev-url="mysql://root:password@localhost:3306/test" \
  --dir="file://ent/migrate/migrations" \
  --latest=1

# 或

atlas migrate lint \
  --dev-url="mysql://root:password@localhost:3306/test" \
  --dir="file://ent/migrate/migrations" \
  --latest=1

# 顯示結果
*_users_add_height.sql: data dependent changes detected:

    L2: Adding a non-nullable "double" column "height" on table "users" without a default value implicitly sets existing rows with 0

# 修正
# 刪除 *_users_add_height.sql
# 還原 atlas.sum
```

```diff
# ent/schema/user.go
# 預設值為0

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
		field.Int("age").
			Positive(),
		field.Float("height").
			Positive().
+     Default(0),
	}
}
```

```sh
# 重做
go run -mod=mod ent/migrate/main.go users_add_height

# 驗證
atlas migrate lint \
  --dev-url="mysql://root:password@localhost:3306/test" \
  --dir="file://ent/migrate/migrations" \
  --latest=1
```

## Apply Migrations

```sh
# 同意提交,並修改
# 會在資料庫建立 atlas_schema_revisions 表

atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url mysql://root:password@localhost:3306/ent_atlas_migration
```

## 修改欄位名

```diff
# ent/schema/user.go
-	field.String("name").
+	field.String("nickname").
```

```sh
# 方法一
# 檢查 ent/schema/user.go, 產生代碼
go generate ./ent

# 檢查產生代碼後再產生 *_users_column_rename.sql
go run -mod=mod ent/migrate/main.go users_column_rename
```

```diff
# *_users_column_rename.sql
# 手動更改語法,原語法是新增欄位

- ALTER TABLE `users` ADD COLUMN `nickname` varchar(255) NOT NULL DEFAULT 'unknown';
+ ALTER TABLE `users` RENAME COLUMN `name` TO `nickname`;
```

```sh
# 更新 atlas.sum 碼
atlas migrate hash \
  --dir "file://ent/migrate/migrations"

# lint
# apply
# 完成
```

```sh
# 方法二
atlas migrate new users_column_rename \
  --dir "file://ent/migrate/migrations"
```

```go
// *_users_column_rename.sql
// 寫更換欄位名字的sql語法

ALTER TABLE `users` RENAME COLUMN `name` TO `nickname`;
```

```sh
# 更新 atlas.sum 碼
atlas migrate hash \
  --dir "file://ent/migrate/migrations"

# lint
# apply
# 完成
```

## 驗證

```sh
#
atlas migrate validate --dir file://ent/migrate/migrations
```
