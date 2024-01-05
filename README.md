# 使用指南

- 在项目根目录执行命令，生成辅助代码（zz_generated.options.go)[注意：文件 config/config.go 里有 go:generate 命令，一旦修改了Config之后发现日志没有打印对应的配置项，查看 go:generate 命令是否完全指定了所有要打印的 struct]
```shell
go generate ./...
```

# 设计架构

- 使用工厂模式构建数据模块，将业务端与数据端完全分离
- 使用工厂模式构建上报模块，可以将内容上报到 es，redis 等

数据上报的条件：
- 批量投递：buffer 已满，马上投递
- 超时投递：只要 buffer 有数据，并且到达了投递时间，也会马上投递
