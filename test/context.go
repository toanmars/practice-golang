package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done(): // Lắng nghe tín hiệu từ bộ đàm
			fmt.Printf("Nhân viên %d: Đã nhận lệnh dừng, đang dọn dẹp đồ đạc...\n", id)
			return // Kết thúc công việc ngay lập tức
		default:
			fmt.Printf("Nhân viên %d: Đang làm việc hăng say...\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	// 1. Tạo một Context có thời gian hết hạn là 2 giây
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Giải phóng tài nguyên khi xong việc

	// 2. Cho nhân viên đi làm
	go worker(ctx, 1)
	go worker(ctx, 2)

	// 3. Leader đợi xem chuyện gì xảy ra
	time.Sleep(3 * time.Second)
	fmt.Println("Leader: Kết thúc buổi làm việc.")
}
