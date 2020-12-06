package log

type cronLogger struct {
	logger Logger
}

// Following implementation here: https://github.com/robfig/cron/blob/master/logger.go

func (c *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	c.logger.WithFields(transformKeysAndValues(keysAndValues...)).Debug(msg)
}

func (c *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	c.logger.WithFields(transformKeysAndValues(keysAndValues...)).WithError(err).Error(msg)
}

func transformKeysAndValues(keysAndValues ...interface{}) map[string]interface{} {
	// Prepare result
	res := map[string]interface{}{}

	// loop over inputs
	for i := 0; i < len(keysAndValues); i += 2 {
		res[keysAndValues[i].(string)] = keysAndValues[i+1]
	}

	return res
}
