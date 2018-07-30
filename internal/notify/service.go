package notify

type ServiceProvider interface {
	Send(args map[string]interface{}) error
}

var Service ServiceProvider
