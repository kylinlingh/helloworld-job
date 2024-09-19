
1. 使用低版本 go 来编译，提示错误：note: module requires Go 1.20
完整错误提示如下：
```shell
 ~/Code/project-sast/golang-sast/cmd/app   master ±✚  go build   
# golang.org/x/exp/slog
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/level.go:159:13: undefined: atomic.Int64
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/attr.go:20:19: undefined: StringValue
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/attr.go:68:19: undefined: GroupValue
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/handler.go:446:15: undefined: StringValue
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/json_handler.go:109:20: v.str undefined (type Value has no field or method str)
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/record.go:192:9: undefined: GroupValue
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/text_handler.go:99:20: v.str undefined (type Value has no field or method str)
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/value.go:87:7: undefined: stringptr
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/value.go:91:7: undefined: groupptr
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/value.go:173:10: undefined: StringValue
../../../../../../.gvm/pkgsets/go1.18.10/global/pkg/mod/golang.org/x/exp@v0.0.0-20230905200255-921286631fa9/slog/value.go:173:10: too many errors
note: module requires Go 1.20
```
解法：降低 go.mod 里的 viper 版本为：github.com/spf13/viper v1.16.0