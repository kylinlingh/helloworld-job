package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"syscall"
	"time"
)

var stdErrFileHandler *os.File

// RewriteStderrFile 把程序运行时的标准错误替换成日志文件，Go在panic的时候它还是往标准错误里写，只不过这里把标准错误的文件描述符换成了日志文件的描述符
func RewriteStderrFile(logDir string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	// 获取当前时间
	currentTime := time.Now()
	// 格式化时间为字符串，作为文件名
	fileName := currentTime.Format("20060102_150405") // 根据需要的格式进行调整
	stdErrFile := path.Join(logDir, fmt.Sprintf("panic-%s.log", fileName))
	file, err := os.OpenFile(stdErrFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	stdErrFileHandler = file //把文件句柄保存到全局变量，避免被GC回收

	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		fmt.Println(err)
		return err
	}
	// 内存回收前关闭文件描述符
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})

	return nil
}
