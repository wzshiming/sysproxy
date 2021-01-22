package sysproxy

import (
	"context"
	"fmt"
	"syscall"

	"github.com/getlantern/sysproxy"
	"github.com/wzshiming/notify"
)

func SysProxy(ctx context.Context, proxy string) error {
	err := sysproxy.EnsureHelperToolPresent("sysproxy-cmd", "Input your password and save the world!", "")
	if err != nil {
		return fmt.Errorf("EnsureHelperToolPresent: %w", err)
	}
	off, err := sysproxy.On(proxy)
	if err != nil {
		if off != nil {
			off()
		}
		return fmt.Errorf("Set proxy: %w", err)
	}
	defer off()
	ctx, cancel := context.WithCancel(ctx)
	cancelNotify := notify.Once(cancel, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()
	cancelNotify()
	return nil
}
