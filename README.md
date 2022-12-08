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
