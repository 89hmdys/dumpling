# dumpling

## 起因
*   嗯，微信并没有提供Golang的SDK，╮(╯▽╰)╭，只好自己刀耕火种了，目前仅支持公众号支付和APP支付。

### 例子
1.公众号支付
    
    /*
      创建用于微信公众号的client
    */
    client:=wechat.NewSubscriptionClient("appId", "appKey", "mchId", "notifyUrl")
    /*
      使用OrderBuilder创建支付参数
    */
    orderBuilder := &wechat.OrderBuilder{}
	order := orderBuilder.BuildForSubscription("订单编号","商品描述","用户OPENID","商品价格")
    /*
      创建预支付订单，将预支付订单返回前台
    */
	prepareOrder, err := client.Pay(order)
	if err != nil {
		fmt.Println(err)
	}


2.APP支付
    
    /*
      创建用于微信公众号的client
    */
    client:=wechat.NewAppClient("appId", "appKey", "mchId", "notifyUrl")
    /*
      使用OrderBuilder创建支付参数
    */
    orderBuilder := &wechat.OrderBuilder{}
	order := orderBuilder.BuildForApp("订单编号","商品描述","用户OPENID","商品价格")
    /*
      创建预支付订单，将预支付订单返回前台
    */
	prepareOrder, err := client.Pay(order)
	if err != nil {
		fmt.Println(err)
	}

3.处理微信异步通知

    /*
      微信的通知消息是放在request的content内的xml字符串，
      从content中取出来，作为参数传入Notify方法就OK了，支付成功会返回
      type Receipt struct {
	      OrderNo    string
	      TradeNo    string
	      TotalPrice string
      }这样一个struct的实例，如果失败，返回error
    */
    receipt,err:=client.Notify(notifyMessge)
    if err != nil {
		fmt.Println(err)
	}