package wechat

func NewAppClient(appId, appKey, mchId, notifyUrl string) Pay {
	return &app{wechat: wechat{
		AppId:     appId,
		AppKey:    appKey,
		MchId:     mchId,
		NotifyUrl: notifyUrl}}
}

func NewSubscriptionClient(appId, appKey, mchId, notifyUrl string) Pay {
	return &subscription{wechat: wechat{
		AppId:     appId,
		AppKey:    appKey,
		MchId:     mchId,
		NotifyUrl: notifyUrl}}
}
