package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func processOrder(ctx context.Context, orderID int, logChan chan<- string) {
	// Giả lập thời gian xử lý đơn hàng ngẫu nhiên từ 1 đến 5 giây
	// (Có đơn nhanh, có đơn chậm hơn cả thời gian timeout)
	workTime := time.Duration(orderID%5+1) * time.Second

	select {
	case <-time.After(workTime):
		// Nếu xử lý xong trước khi bị hủy
		msg := fmt.Sprintf("✅ Đơn hàng %d: Xử lý xong sau %v", orderID, workTime)
		logChan <- fmt.Sprintf("DONE: Order %d", orderID)
		fmt.Println(msg)
	case <-ctx.Done():
		// Nếu bộ đàm báo HỦY (do quá 3 giây)
		fmt.Sprintf("❌ Đơn hàng %d: Bị hủy do quá thời gian!", orderID)
	}
}

func main() {
	// 1. Khởi tạo Context - Hết hạn sau 3 giây
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 2. Khởi tạo Channel (Buffered) và WaitGroup
	logChan := make(chan string, 10)
	var wg sync.WaitGroup

	fmt.Println("=== HỆ THỐNG BẮT ĐẦU NHẬN 5 ĐƠN HÀNG ===")

	for i := 1; i <= 5; i++ {
		go processOrder(ctx, i, logChan)
	}

	// 3. Goroutine riêng biệt để ghi Log (Dùng WaitGroup ở đây)
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("--- Bộ phận ghi Log đang trực... ---")
		for {
			select {
			case msg := <-logChan:
				fmt.Println("   📝 Ghi vào DB:", msg)
			case <-ctx.Done():
				// Khi hệ thống dừng, bộ phận log cũng nghỉ
				fmt.Println("--- Bộ phận ghi Log đóng cửa ---")
				return
			}
		}
	}()

	// Đợi lâu hơn timeout một chút để quan sát kết quả
	time.Sleep(4 * time.Second)

	// Đảm bảo bộ phận ghi log đã xong việc
	wg.Wait()
	fmt.Println("=== HỆ THỐNG NGHỈ NGƠI ===")
}
