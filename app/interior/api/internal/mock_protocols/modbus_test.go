package mock

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"

	//"xunjikeji.com.cn/gateway/app/interior/api/internal/mock"
	"testing"
)

func Test_MockProtocol_Start(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockProtocol := NewMockProtocol(ctrl)

	mockProtocol.EXPECT().Start(ctx).Do(func() {
		fmt.Println("hello")
	})

	mockProtocol.Start(ctx)

}
