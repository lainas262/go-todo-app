package hub

type Router struct {
	rules map[string]EventHandler
}

func NewEventRouter() *Router {
	return &Router{rules: make(map[string]EventHandler)}
}

// register route and event handler
func (r *Router) Register(eventName string, handler EventHandler) {
	r.rules[eventName] = handler
}
