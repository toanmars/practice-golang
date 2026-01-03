# Hệ thống xử lý đơn hàng.

Kịch bản:
1. Leader (Main): Nhận đơn hàng và giới hạn thời gian xử lý là 3 giây (Context).

2. Công nhân (Goroutine): Thực hiện xử lý đơn hàng.

3. Băng chuyền (Channel): Gửi kết quả xử lý về để lưu vào lịch sử.

4. Chốt chặn (WaitGroup): Đảm bảo tất cả báo cáo đã được ghi lại trước khi tắt hệ thống.

```go 
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
		msg := fmt.Printf("✅ Đơn hàng %d: Xử lý xong sau %v", orderID, workTime)
		logChan <- fmt.Sprintf("DONE: Order %d", orderID)
		fmt.Println(msg)
	case <-ctx.Done():
		// Nếu bộ đàm báo HỦY (do quá 3 giây)
		fmt.Printf("❌ Đơn hàng %d: Bị hủy do quá thời gian!\n", orderID)
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
```

Giải thích sự phối hợp "nhịp nhàng":

- Context (Bộ đàm): Nó là "nhịp tim" của cả chương trình. Khi 3 giây trôi qua, ctx.Done() phát tín hiệu cho tất cả các đơn hàng đang chạy dở và cả bộ phận ghi Log cùng biết để dừng lại.

- Goroutine (Nhân viên): 5 đơn hàng được xử lý cùng lúc (Parallelism nếu CPU đa nhân). Đơn nào làm kịp 3 giây thì xong, đơn nào chậm thì tự nghỉ.

- Channel (Băng chuyền): Giúp chuyển dữ liệu an toàn từ các đơn hàng về bộ phận Log mà không bị xung đột bộ nhớ.

- WaitGroup (Chốt chặn): Dùng để chắc chắn rằng bộ phận Log đã nhận đủ tin nhắn và in ra màn hình xong xuôi thì mới cho phép hàm main thoát hoàn toàn.

## Viết unit test cho bài toán trên

Viết Unit Test cho code đa nhiệm (Concurrency) trong Go là một thử thách thú vị, vì bạn phải kiểm tra xem các thành phần có phối hợp đúng như ý muốn hay không mà không bị "treo" test.

Để test được đoạn code trên, chúng ta cần cấu trúc lại code một chút để có thể truyền vào các giá trị giả lập (mock).

## 1. Cấu trúc lại Code để Test được (Testable Code)
Thay vì chạy thẳng trong main, ta tách logic xử lý đơn hàng ra một hàm riêng.

```go 
package main

import (
	"context"
	"fmt"
	"time"
)

≈```
### 2. Viết Unit Test (main_test.go)
Chúng ta sẽ dùng thư viện testing chuẩn của Go. Ta cần test 2 trường hợp:

- Trường hợp thành công: Xử lý xong trước khi hết thời gian.

- Trường hợp quá hạn (Timeout): Context báo hủy trước khi xử lý xong.

```go 
package main

import (
	"testing"
	"time"
)

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
```

### 3. Cách chạy Test
Bạn mở terminal và gõ:

```Bash
 go test -v
``` 
Đặc biệt: Kiểm tra xung đột dữ liệu (Race Condition) Khi làm việc với Goroutine, có một lệnh cực kỳ quan trọng mà mọi Dev Go phải dùng để xem các Goroutine có đang "đánh nhau" để giành giật biến nào không:

```Bash
go test -race -v
``` 
Nếu kết quả hiện PASS và không có cảnh báo màu đỏ nào, nghĩa là code của bạn an toàn.

### 4. Giải thích kỹ thuật Test Concurrency
- Dùng t.Run: Chia nhỏ các trường hợp (Success/Timeout) để dễ quản lý.

- Dùng select trong test: Khi đọc dữ liệu từ channel trong test, hãy luôn dùng select kèm default hoặc một cái time.After. Điều này giúp test của bạn không bị "treo" mãi mãi nếu code logic bị lỗi không gửi dữ liệu vào channel.

- Kiểm soát thời gian: Trong Unit Test, chúng ta thường dùng thời gian rất ngắn (milisecond) để test chạy nhanh, không nên để time.Sleep quá lâu.