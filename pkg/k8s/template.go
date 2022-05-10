package k8s

type AlertManagerConfig struct {
	Name   string
	Router []Router
}

type Router struct {
	Matchers       map[string]string
	RepeatInterval string
}
