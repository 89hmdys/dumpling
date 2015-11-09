package wechat

type Pay interface {
	Notify(s string) (*Notification, error)
	Pay(order Order) (map[string]string, error)
}
