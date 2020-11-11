package server

type customHealthChecker struct {
	fn func() error
}

func (chc *customHealthChecker) Status() (interface{}, error) {
	// Run function
	err := chc.fn()

	return nil, err
}
