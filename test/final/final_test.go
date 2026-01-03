package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Hàm xử lý đơn hàng được tách ra để dễ test
func ProcessOrder(ctx context.Context, orderID int, logChan chan<- string, workTime time.Duration) string {
	select {
	case <-time.After(workTime):
		msg := fmt.Sprintf("DONE: Order %d", orderID)
		logChan <- msg
		return "success"
	case <-ctx.Done():
		return "timeout"
	}
}
func TestProcessOrder(t *testing.T) {
	// Test trường hợp thành công
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logChan := make(chan string)
	result := ProcessOrder(ctx, 1, logChan, 2*time.Second)

	if result != "success" {
		t.Errorf("Expected success, got %s", result)
	}

	// Test trường hợp quá hạn
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	logChan = make(chan string)
	result = ProcessOrder(ctx, 2, logChan, 2*time.Second)

	if result != "timeout" {
		t.Errorf("Expected timeout, got %s", result)
	}
}
