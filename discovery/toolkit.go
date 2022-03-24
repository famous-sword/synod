package discovery

import "synod/conf"

func forPublish(name string) string {
	return "/synod/" + name + conf.String("app.id")
}

func forSubscribe(name string) string {
	return "/synod/" + name
}
