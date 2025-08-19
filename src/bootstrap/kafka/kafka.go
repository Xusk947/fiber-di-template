package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gitlab.stat4market.com/reelsmarket/fiber-di-server-template/src/bootstrap/config"
	"go.uber.org/zap"
	"time"
)

// Kafka represents a Kafka client with reader and writer functionality
type Kafka struct {
	Config     *config.KafkaConfig
	Logger     *zap.Logger
	Writers    map[string]*kafka.Writer
	Readers    map[string]*kafka.Reader
	Dialer     *kafka.Dialer
	Conn       *kafka.Conn
	Partitions []kafka.Partition
}

// NewKafka creates a new Kafka client with the provided config
func NewKafka(cfg *config.Config, logger *zap.Logger) (*Kafka, error) {
	logger.Info("🔌 Initializing Kafka connection")

	// Create a custom dialer for authentication if needed
	dialer := &kafka.Dialer{
		Timeout:   time.Duration(cfg.Kafka.TimeoutSeconds) * time.Second,
		DualStack: true,
	}

	// Set up SASL authentication if credentials are provided
	if cfg.Kafka.Username != "" && cfg.Kafka.Password != "" {
		dialer.SASLMechanism = kafka.PLAIN
		dialer.SASLUsername = cfg.Kafka.Username
		dialer.SASLPassword = cfg.Kafka.Password
	}

	// Establish a connection to test connectivity and discover partitions
	kafkaBrokers := cfg.Kafka.Brokers
	if len(kafkaBrokers) == 0 {
		return nil, fmt.Errorf("no Kafka brokers configured")
	}

	conn, err := dialer.Dial("tcp", kafkaBrokers[0])
	if err != nil {
		logger.Error("❌ Failed to connect to Kafka", zap.Error(err))
		return nil, err
	}

	client := &Kafka{
		Config:  &cfg.Kafka,
		Logger:  logger.Named("kafka"),
		Writers: make(map[string]*kafka.Writer),
		Readers: make(map[string]*kafka.Reader),
		Dialer:  dialer,
		Conn:    conn,
	}

	return client, nil
}

// CreateWriter creates a new Kafka writer for the given topic
func (k *Kafka) CreateWriter(topic string) *kafka.Writer {
	if writer, exists := k.Writers[topic]; exists {
		return writer
	}

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(k.Config.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: k.Config.AllowAutoTopicCreation,
		BatchTimeout:           time.Duration(k.Config.BatchTimeoutMs) * time.Millisecond,
		BatchSize:              k.Config.BatchSize,
		Async:                  k.Config.Async,
		Logger:                 kafka.LoggerFunc(k.logKafkaMessage),
		ErrorLogger:            kafka.LoggerFunc(k.logKafkaError),
	}

	k.Writers[topic] = writer
	k.Logger.Info("Created Kafka writer", zap.String("topic", topic))
	return writer
}

// CreateReader creates a new Kafka reader for the given topic and group
func (k *Kafka) CreateReader(topic, groupID string) *kafka.Reader {
	key := fmt.Sprintf("%s-%s", topic, groupID)
	if reader, exists := k.Readers[key]; exists {
		return reader
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        k.Config.Brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       k.Config.MinBytes,
		MaxBytes:       k.Config.MaxBytes,
		CommitInterval: time.Duration(k.Config.CommitIntervalMs) * time.Millisecond,
		StartOffset:    kafka.FirstOffset,
		Logger:         kafka.LoggerFunc(k.logKafkaMessage),
		ErrorLogger:    kafka.LoggerFunc(k.logKafkaError),
		Dialer:         k.Dialer,
	})

	k.Readers[key] = reader
	k.Logger.Info("Created Kafka reader", zap.String("topic", topic), zap.String("groupID", groupID))
	return reader
}

// WriteMessage writes a message to the specified topic
func (k *Kafka) WriteMessage(ctx context.Context, topic string, key, value []byte) error {
	writer := k.CreateWriter(topic)
	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	})

	if err != nil {
		k.Logger.Error("Failed to write message to Kafka",
			zap.String("topic", topic),
			zap.Error(err))
		return err
	}

	return nil
}

// ConsumeMessages starts consuming messages from the specified topic and group
// The handler function is called for each message received
func (k *Kafka) ConsumeMessages(ctx context.Context, topic, groupID string, handler func(context.Context, kafka.Message) error) {
	reader := k.CreateReader(topic, groupID)

	go func() {
		k.Logger.Info("Starting Kafka consumer",
			zap.String("topic", topic),
			zap.String("groupID", groupID))

		for {
			select {
			case <-ctx.Done():
				k.Logger.Info("Stopping Kafka consumer",
					zap.String("topic", topic),
					zap.String("groupID", groupID))
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					// Check if context was canceled
					if ctx.Err() != nil {
						return
					}
					
					k.Logger.Error("Error reading message from Kafka",
						zap.String("topic", topic),
						zap.Error(err))
					continue
				}

				// Process the message
				if err := handler(ctx, msg); err != nil {
					k.Logger.Error("Error processing Kafka message",
						zap.String("topic", topic),
						zap.Error(err))
				}
			}
		}
	}()
}

// Close closes all Kafka connections
func (k *Kafka) Close() error {
	var lastErr error

	// Close all readers
	for key, reader := range k.Readers {
		if err := reader.Close(); err != nil {
			k.Logger.Error("Failed to close Kafka reader",
				zap.String("reader", key),
				zap.Error(err))
			lastErr = err
		}
	}

	// Close all writers
	for topic, writer := range k.Writers {
		if err := writer.Close(); err != nil {
			k.Logger.Error("Failed to close Kafka writer",
				zap.String("topic", topic),
				zap.Error(err))
			lastErr = err
		}
	}

	// Close main connection
	if k.Conn != nil {
		if err := k.Conn.Close(); err != nil {
			k.Logger.Error("Failed to close Kafka connection", zap.Error(err))
			lastErr = err
		}
	}

	k.Logger.Info("Kafka connections closed")
	return lastErr
}

// onStart is called during application startup
func (k *Kafka) onStart(ctx context.Context) error {
	k.Logger.Info("✅ Kafka client initialized")
	return nil
}

// onStop is called during application shutdown
func (k *Kafka) onStop(ctx context.Context) error {
	k.Logger.Info("🛑 Shutting down Kafka client")
	return k.Close()
}

// logKafkaMessage is a helper function to log Kafka messages to zap
func (k *Kafka) logKafkaMessage(msg string, args ...interface{}) {
	k.Logger.Debug(fmt.Sprintf(msg, args...))
}

// logKafkaError is a helper function to log Kafka errors to zap
func (k *Kafka) logKafkaError(msg string, args ...interface{}) {
	k.Logger.Error(fmt.Sprintf(msg, args...))
}
