package config

func init() {
	registerValidator("awslambda", func(s Settings) error {
		c := defaultAWSLambdaHandlerConfig()
		if v, ok := s["function_name"].(string); ok {
			c.FunctionName = v
		}
		if v, ok := s["region"].(string); ok {
			c.Region = v
		}
		if v, ok := s["access_key"].(string); ok {
			c.AccessKey = v
		}
		if v, ok := s["secret_key"].(string); ok {
			c.SecretKey = v
		}
		if v, ok := s["invocation_type"].(string); ok {
			c.InvocationType = v
		}
		if v, ok := s["timeout_seconds"].(int64); ok {
			c.Timeout = int(v)
		}
		return validateAWSLambdaHandler(c)
	})
}
