package async

type KafkaConfig struct {
	Name  string
	Topic string
}

type Job struct {
	Backend BackendType
	Kafka   *KafkaConfig
}
