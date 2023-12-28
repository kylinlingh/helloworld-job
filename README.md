# 使用指南

- 在项目根目录执行命令，生成辅助代码（zz_generated.options.go)
```shell
go generate ./...
```

# 设计架构

- 使用工厂模式构建数据模块，将业务端与数据端完全分离
- 使用工厂模式构建上报模块，可以将内容上报到 es，redis 等