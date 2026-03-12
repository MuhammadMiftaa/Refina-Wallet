package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"refina-wallet/internal/service/mocks"
	"refina-wallet/internal/types/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- helpers ----------

func newOutboxPublisher(
	outboxRepo *mocks.MockOutboxRepository,
	rabbitMQ *mocks.MockRabbitMQClient,
) *OutboxPublisher {
	return NewOutboxPublisher(outboxRepo, rabbitMQ)
}

func sampleOutboxMessages() []model.OutboxMessage {
	return []model.OutboxMessage{
		{
			ID:          1,
			AggregateID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			EventType:   "wallet.created",
			Payload:     []byte(`{"id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","name":"BCA"}`),
			Published:   false,
			Retries:     0,
			MaxRetries:  5,
			CreatedAt:   time.Now(),
		},
	}
}

// =====================================================================
// NewOutboxPublisher
// =====================================================================

func TestNewOutboxPublisher(t *testing.T) {
	repo := new(mocks.MockOutboxRepository)
	rabbitMQ := new(mocks.MockRabbitMQClient)

	publisher := newOutboxPublisher(repo, rabbitMQ)

	assert.NotNil(t, publisher)
	assert.Equal(t, repo, publisher.outboxRepo)
	assert.Equal(t, rabbitMQ, publisher.queue)
	assert.Equal(t, 100, publisher.batchSize)
	assert.Equal(t, 5*time.Second, publisher.interval)
}

// =====================================================================
// publishPendingMessages
// =====================================================================

func TestPublishPendingMessages_NoPendingMessages(t *testing.T) {
	repo := new(mocks.MockOutboxRepository)
	rabbitMQ := new(mocks.MockRabbitMQClient)

	publisher := newOutboxPublisher(repo, rabbitMQ)

	repo.On("GetPendingMessages", mock.Anything, publisher.batchSize).
		Return([]model.OutboxMessage{}, nil)

	err := publisher.publishPendingMessages(context.Background())

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	rabbitMQ.AssertNotCalled(t, "GetChannel")
}

func TestPublishPendingMessages_GetPendingError(t *testing.T) {
	repo := new(mocks.MockOutboxRepository)
	rabbitMQ := new(mocks.MockRabbitMQClient)

	publisher := newOutboxPublisher(repo, rabbitMQ)

	repo.On("GetPendingMessages", mock.Anything, publisher.batchSize).
		Return([]model.OutboxMessage{}, errors.New("db error"))

	err := publisher.publishPendingMessages(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get pending messages")
	repo.AssertExpectations(t)
}

func TestPublishPendingMessages_GetChannelError(t *testing.T) {
	repo := new(mocks.MockOutboxRepository)
	rabbitMQ := new(mocks.MockRabbitMQClient)

	publisher := newOutboxPublisher(repo, rabbitMQ)

	messages := sampleOutboxMessages()
	repo.On("GetPendingMessages", mock.Anything, publisher.batchSize).Return(messages, nil)
	rabbitMQ.On("GetChannel").Return(nil, errors.New("channel error"))
	repo.On("IncrementRetries", mock.Anything, messages[0].ID).Return(nil)

	err := publisher.publishPendingMessages(context.Background())

	assert.NoError(t, err) // errors are logged per-message, not returned
	repo.AssertExpectations(t)
	rabbitMQ.AssertExpectations(t)
}

func TestPublishPendingMessages_IncrementRetriesError(t *testing.T) {
	repo := new(mocks.MockOutboxRepository)
	rabbitMQ := new(mocks.MockRabbitMQClient)

	publisher := newOutboxPublisher(repo, rabbitMQ)

	messages := sampleOutboxMessages()
	repo.On("GetPendingMessages", mock.Anything, publisher.batchSize).Return(messages, nil)
	rabbitMQ.On("GetChannel").Return(nil, errors.New("channel error"))
	repo.On("IncrementRetries", mock.Anything, messages[0].ID).Return(errors.New("db error"))

	err := publisher.publishPendingMessages(context.Background())

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPublishPendingMessages_MaxRetriesExceeded(t *testing.T) {
	repo := new(mocks.MockOutboxRepository)
	rabbitMQ := new(mocks.MockRabbitMQClient)

	publisher := newOutboxPublisher(repo, rabbitMQ)

	messages := sampleOutboxMessages()
	messages[0].Retries = 4 // at MaxRetries-1 = 4, should log max retries
	messages[0].MaxRetries = 5

	repo.On("GetPendingMessages", mock.Anything, publisher.batchSize).Return(messages, nil)
	rabbitMQ.On("GetChannel").Return(nil, errors.New("channel error"))
	repo.On("IncrementRetries", mock.Anything, messages[0].ID).Return(nil)

	err := publisher.publishPendingMessages(context.Background())

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPublishPendingMessages_MarkPublishedError(t *testing.T) {
	// NOTE: This test cannot be fully exercised because publishMessage uses
	// amqp091.Channel directly (not an interface), so we cannot mock a successful
	// publish followed by a MarkAsPublished error without a real RabbitMQ connection.
	// The GetChannel-error path is already covered by TestPublishPendingMessages_GetChannelError.
	// If the amqp channel were wrapped behind an interface, this could be fully tested.
	t.Skip("publishMessage requires a real amqp091.Channel; covered by integration tests")
}

// =====================================================================
// Start — context cancellation
// =====================================================================

func TestStart_ContextCancellation(t *testing.T) {
	repo := new(mocks.MockOutboxRepository)
	rabbitMQ := new(mocks.MockRabbitMQClient)

	publisher := newOutboxPublisher(repo, rabbitMQ)
	publisher.interval = 10 * time.Millisecond // speed up for test

	ctx, cancel := context.WithCancel(context.Background())

	// Mock for any GetPendingMessages call that might fire during the tick
	repo.On("GetPendingMessages", mock.Anything, publisher.batchSize).
		Return([]model.OutboxMessage{}, nil).Maybe()

	done := make(chan struct{})
	go func() {
		publisher.Start(ctx)
		close(done)
	}()

	// Let at least one tick fire
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// goroutine exited properly
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not exit after context cancellation")
	}
}

// =====================================================================
// StartCleanupJob — context cancellation
// =====================================================================

func TestStartCleanupJob_ContextCancellation(t *testing.T) {
	repo := new(mocks.MockOutboxRepository)
	rabbitMQ := new(mocks.MockRabbitMQClient)

	publisher := newOutboxPublisher(repo, rabbitMQ)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		publisher.StartCleanupJob(ctx)
		close(done)
	}()

	// Cancel immediately since cleanup interval is 1 hour
	cancel()

	select {
	case <-done:
		// goroutine exited properly
	case <-time.After(2 * time.Second):
		t.Fatal("StartCleanupJob did not exit after context cancellation")
	}
}

// =====================================================================
// cleanupOldMessages
// =====================================================================

func TestCleanupOldMessages(t *testing.T) {
	repo := new(mocks.MockOutboxRepository)
	rabbitMQ := new(mocks.MockRabbitMQClient)

	publisher := newOutboxPublisher(repo, rabbitMQ)

	err := publisher.cleanupOldMessages(context.Background())

	assert.NoError(t, err) // currently returns nil
}
