package shutdown

import "sync"

// ShutdownManager 监听shutdown事件
// 在收到shutdown事件后，触发GSInterface的Shutdown流程
type ShutdownManager interface {
	GetName() string
	Start(gs GSInterface) error // 通知GSInterface开始Shutdown流程
	ShutdownStart() error       // 在执行Shutdown回调函数前执行,由GSInterface调用
	ShutdownFinish() error      // 在执行Shutdown回调函数后执行,由GSInterface调用
}

// GSInterface 管理Shutdown时的回调事件
type GSInterface interface {
	StartShutdown(sm ShutdownManager)
	ReportError(err error) // 在发生错误时调用,
	AddShutdownCallback(shutdownCallback ShutdownCallback)
}

type ShutdownCallback func(string) error

// OnShutdown 定义在Shutdown发生时需要执行的操作
func (f ShutdownCallback) OnShutdown(shutdownManager string) error {
	return f(shutdownManager) // 在这里执行调用函数
}

// ErrorHandler 可以传递错误处理函数来处理异步错误
type ErrorHandler interface {
	OnError(err error)
}

// 接口的隐式实现（duck typing）
// 方法可绑定到任何类型（包括函数类型）
// 函数作为一等公民的特性
type ErrorFunc func(err error)

// OnError 定义在发生错误时需要执行的操作; 函数类也可以定义方法,以此来实现接口.
func (f ErrorFunc) OnError(err error) {
	f(err)
}

type GracefulShutdown struct {
	managers     []ShutdownManager
	callbacks    []ShutdownCallback
	errorHandler ErrorHandler
}

func New() *GracefulShutdown {
	return &GracefulShutdown{
		managers:  make([]ShutdownManager, 0, 10), // length=0, cap=10
		callbacks: make([]ShutdownCallback, 0, 3),
	}
}

func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.managers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}
	return nil
}

func (gs *GracefulShutdown) SetErrorHandler(f ErrorHandler) {
	gs.errorHandler = f
}

// ReportError 向注册的error handler报告错误
func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}

func (gs *GracefulShutdown) AddShutdownManager(sm ShutdownManager) {
	gs.managers = append(gs.managers, sm)
}

// StartShutdown Shutdown流程
func (gs *GracefulShutdown) StartShutdown(sm ShutdownManager) {
	gs.ReportError(sm.ShutdownStart())

	// 回调函数可以并发执行的
	// 当需要跨协程同步状态时，优先使用闭包捕获变量
	wg := sync.WaitGroup{}
	for _, shutdownCallback := range gs.callbacks {
		wg.Add(1)
		go func(shutdownCallback ShutdownCallback) {
			defer wg.Done()

			gs.ReportError(shutdownCallback.OnShutdown(sm.GetName()))
		}(shutdownCallback) // 循环变量显式传参,避免循环中所有携程共享最后一个回调
	}

	wg.Wait()

	gs.ReportError(sm.ShutdownFinish()) // 调用sm的一些清理操作
}

func (gs *GracefulShutdown) AddShutdownCallback(f ShutdownCallback) {
	gs.callbacks = append(gs.callbacks, f)
}

var _ GSInterface = (*GracefulShutdown)(nil)
