package wechat_test

import (
	"dumpling/wechat"
	"fmt"
	"time"
)

func init() {
	orderBuilder := &wechat.OrderBuilder{}
	order := orderBuilder.SetAttach("21312312").SetExpireTime(time.Now()).BuildForApp("DDD", "321", 1122)

	//	wechat.NewSubscriptionClient()

	fmt.Println(order)
}
