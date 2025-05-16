package posixsignal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Ranper/iam/pkg/shutdown"
)

const Name = "PosixSignalManager"

type PosixSignalManager struct {
	signals []os.Signal
}

func NewPosixSignalManager(sig ...os.Signal) shutdown.ShutdownManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = os.Interrupt    // Ctrl+C
		sig[1] = syscall.SIGTERM // kill <PID>.  SIGKILL -> kill -9 <PID>
	}

	sm := &PosixSignalManager{
		signals: sig,
	}

	return sm
}

func (p *PosixSignalManager) GetName() string {
	return Name
}

// Start 开始监听Shutdown信号
func (p *PosixSignalManager) Start(gs shutdown.GSInterface) error {
	// 在程序启动的时候调用, 不能阻塞启动流程,所以放在go routine中
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, p.signals...) // 当发生这些信号时, 会写入c中

		<-c // 然后解除阻塞

		gs.StartShutdown(p) // 启动Shutdown流程
	}()
	return nil
}

// ShutdownStart does nothing.
func (p *PosixSignalManager) ShutdownStart() error {
	return nil
}

// ShutdownFinish exits the app with os.Exit(0).
func (p *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)

	return nil
}

// var _ PosixSignalManager = (shutdown.ShutdownManager)(nil)
var _ shutdown.ShutdownManager = (*PosixSignalManager)(nil)

/*
第一次执行go run时：
	shell 启动一个go run进程（PID 471142）
	go run会先编译你的程序，然后启动编译后的可执行文件
	编译后的可执行文件作为子进程（PID 471367）运行
	父进程go run可能会等待子进程结束，或者在子进程启动后退出
第一次发送 SIGTERM（kill 471142）时：
	你杀死的是go run命令本身（父进程）
	但你的应用程序（子进程）仍在运行，因为子进程没有收到 SIGTERM 信号
	这就是为什么你看到 "Terminated"，但程序仍在运行
第二次发送 SIGTERM（kill 471367）时：
	这次你杀死了真正的应用程序进程
	你的信号处理代码收到 SIGTERM，触发优雅停止逻辑
*/
