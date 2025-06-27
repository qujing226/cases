package Context

import (
	"fmt"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestCreateCtx(t *testing.T) {
	ctx := context.Background()
	// todo表明未来context可能会从参数传递进来
	//ctx1 := Context.TODO()
	ctx = context.WithValue(ctx, "userId", 1024)
	func1(ctx)
}

func func1(ctx context.Context) {
	userId := ctx.Value("userId")

	fmt.Println("userId:", userId)
}

func TestCancel(t *testing.T) {
	ctx := context.Background()
	ctx1, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	//ctx2,cancel2 := Context.WithCancel(ctx1)
	//Go func() {
	//	time.Sleep(10 * time.Second)
	//}()
	select {
	case <-ctx1.Done():
		fmt.Println("ctx1 done")
	}

}
