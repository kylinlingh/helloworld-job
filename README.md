# 使用指南

1. 修改配置文件
2. 修改 internal/entity/taskrecord.go，定义好要写入数据库的结构
3. 执行 docs 下的 sql 语句


## 配置文件


## 生成辅助代码
- 在项目根目录执行命令，生成辅助代码（zz_generated.options.go)[注意：文件 config/config.go 里有 go:generate 命令，一旦修改了Config之后发现日志没有打印对应的配置项，查看 go:generate 命令是否完全指定了所有要打印的 struct]
```shell
go generate ./...
```

# 架构
