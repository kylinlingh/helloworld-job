# 功能特性
- [x] 快速添加 command
- [x] 数据上报
- [x] 快速编译
- [x] Dockerfile
- [x] 统一的错误码
- [ ] 分布式链路跟踪
- [ ] 认证鉴权

# 使用指南

1. 修改配置文件
2. 修改 internal/entity/taskrecord.go，定义好要写入数据库的结构
3. 执行 docs 下的 sql 语句


## 配置文件


## 生成辅助代码
- 在项目根目录执行命令，生成辅助代码（zz_generated.options.go)[注意：文件 config/config.go 里有 go:generate 命令，一旦修改了Config之后发现日志没有打印对应的配置项，查看 go:generate 命令是否完全指定了所有要打印的 struct]
```shell
# 在目录/tools/codegen下编译 codegen 工具
go build

# 将编译后的产物 codegen 放入到bin目录中
mv codegen $GOPATH/bin/codegen

# 安装依赖项目
go get github.com/ecordell/optgen@v0.0.9

# 自动生成代码，如果在执行本命令时出错：c.DebugMap undefined (type *ScanJobParamConfig has no field or method DebugMap)
# 需要先注释掉这行代码后再执行命令：log.Ctx(ctx).Info().Fields(helpers.Flatten(c.DebugMap())).Msg("configuration as: ")
go generate ./...

```

## 启动运行

需要配置好必须的环境变量
```shell
TASK_ID=1;LOG_DIR=./temp-artifact/logs helloworld try 
```
成功运行结束后可以看到： Done running try job.
另一个更复杂的任务启动命令
```shell
TASK_ID=1;LOG_DIR=./temp-artifact/logs helloworld scan 
```
