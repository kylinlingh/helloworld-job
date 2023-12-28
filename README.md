# 使用指南

- 在 internal/server 目录上执行 go generate 命令，生成新的文件：zz_generated.options.go


# 设计架构

- 使用工厂模式构建数据模块，将业务端与数据端完全分离
- 使用工厂模式构建上报模块，可以将内容上报到 es，redis 等