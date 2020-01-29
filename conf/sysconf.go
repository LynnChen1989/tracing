package conf

// 系统全局定义
const ()

var (
	// 日志配置
	LogPath       = "/tmp/log.log" //系统路径
	LogLevel      = "debug"        //日志级别
	LogMaxSize    = 128            //每个日志文件保存的最大尺寸 单位：M
	LogMaxBackups = 10             //日志文件最多保存多少个备份
	LogMaxAge     = 10             //文件最多保存多少天
	LogCompress   = true           //是否压缩
	// 系统配置
	SysName     = "jaeger-ex"      //系统名
	JaegerAgent = "127.0.0.1:6831" //Jaeger地址
	TraceOpen   = 1                //开启调用链追踪
)
