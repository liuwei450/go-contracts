package cycle

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"go-contracts/util"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 定义长期运行服务的生命周期接口
type Service interface {
	// Start 启动服务（阻塞方法，需在 goroutine 中运行）
	Start(ctx context.Context) error
	// Stop 停止服务（可选，用于主动触发清理）
	Stop(ctx context.Context) error
	Stopped() bool
}

// 服务启动函数类型（业务逻辑入口）
type ServiceStartFunc func(ctx *cli.Context, close context.CancelCauseFunc) (Service, error)

func LifecycleCmd(startFn ServiceStartFunc) cli.ActionFunc {
	return func(c *cli.Context) error {
		// 1. 创建根上下文（控制服务生命周期）
		ctx, cancel := context.WithCancelCause(context.Background())
		defer cancel(nil)

		// 2. 启动业务服务（调用用户定义的 startFn）
		service, err := startFn(c, cancel)
		if err != nil {
			util.Log.Crit("服务启动失败", "error", err)
			return err
		}

		// 3. 在 goroutine 中运行服务（避免阻塞信号监听）
		go func() {
			if err := service.Start(ctx); err != nil && err != context.Canceled {
				util.Log.Crit("服务运行异常退出", "error", err)
				cancel(nil) // 服务异常时主动触发退出
			}
		}()

		// 4. 监听系统退出信号（SIGINT/SIGTERM）
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// 5. 等待退出信号或服务主动退出
		util.Log.Info("服务启动成功，等待退出信号（Ctrl+C 或 SIGTERM）")
		select {
		case sig := <-sigChan:
			util.Log.Info(fmt.Sprintf("收到退出信号: %s", sig.String()))
		case <-ctx.Done():
			util.Log.Info("服务主动请求退出")
		}

		// 6. 优雅关停流程
		util.Log.Info("开始优雅关闭服务...")
		cancel(nil) // 通知服务退出

		// 7. 调用服务 Stop 方法清理资源
		if err := service.Stop(ctx); err != nil {
			util.Log.Warn("服务停止失败", "error", err)
		}

		// 8. 等待清理完成（超时控制）
		if err := waitForShutdown(10 * time.Second); err != nil {
			util.Log.Warn("优雅关闭超时", "error", err)
		}

		util.Log.Info("服务已完全关闭")
		return nil
	}

}

// 等待服务清理完成（超时控制）
func waitForShutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return ctx.Err() // 返回超时或取消错误
	}
}
