package posixsignal

import (
	"syscall"
	"testing"
	"time"

	"github.com/Ranper/iam/pkg/shutdown"
)

type startShutdownFunc func(sm shutdown.ShutdownManager)

func (f startShutdownFunc) AddShutdownCallback(shutdownCallback shutdown.ShutdownCallback) {

}

func (f startShutdownFunc) ReportError(err error) {
}

// 开始执行shutdown流程的时候, 执行传入的函数
func (f startShutdownFunc) StartShutdown(sm shutdown.ShutdownManager) {
	f(sm)
}

var _ shutdown.GSInterface = (startShutdownFunc)(nil)

func waitSig(t *testing.T, c <-chan int) {
	select {
	case <-c:

	case <-time.After(10 * time.Second):
		t.Error("Timeout waiting for StartShutdown.")
	}
}

func TestStartShutdownCalledOnDefaultSignals(t *testing.T) {
	c := make(chan int, 10)

	psm := NewPosixSignalManager()

	// psm.Start 监听信号. 在执行shutdown流程的时候, 会调用传入的函数
	psm.Start(startShutdownFunc(func(sm shutdown.ShutdownManager) {
		c <- 1
	}))

	time.Sleep(time.Millisecond)

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	waitSig(t, c)

	// psm.Start 监听信号. 在执行shutdown流程的时候, 会调用传入的函数
	psm.Start(startShutdownFunc(func(sm shutdown.ShutdownManager) {
		c <- 1
	}))

	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	waitSig(t, c)
}

// 终端断开/配置重载

func TestStartShutdownOnCustomSignal(t *testing.T) {
	c := make(chan int, 10)

	psm := NewPosixSignalManager(syscall.SIGHUP)

	psm.Start(startShutdownFunc(func(sm shutdown.ShutdownManager) {
		c <- 1
	}))

	time.Sleep(time.Millisecond)

	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)

	waitSig(t, c)
}
