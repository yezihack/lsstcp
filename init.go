package serverRoom

import (
	"flag"
	"path"

	"github.com/ThreeKing2018/goutil/golog"
	"github.com/ThreeKing2018/goutil/golog/conf"
)

func Init() {
	flag.Parse()

	//打印版本并退出
	if Arg.Getver() {
		printVersion()
	}

	golog.SetLogger(
		golog.ZAPLOG,
		conf.WithLogType(conf.LogJsontype),
		conf.WithLogLevel(conf.DebugLevel),
		conf.WithFilename(path.Join(Arg.logdir, ServiceName+".log")),
	)
}
