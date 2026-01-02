# Mức 3: Context - Cách để Leader "Ra lệnh hủy quân"

Trong thực tế, không phải lúc nào công việc cũng suôn sẻ. Ví dụ: Bạn yêu cầu 10 nhân viên đi tìm thông tin khách hàng. Nhưng khách hàng vừa gọi điện báo: "Tôi không cần nữa". Lúc này, nếu 10 nhân viên vẫn tiếp tục hì hục làm việc thì sẽ cực kỳ lãng phí tài nguyên (CPU, RAM).

Context sinh ra để giải quyết việc này. Nó giống như một cái **"Bộ đàm"** mà Leader dùng để hét lên: **"TẤT CẢ DỪNG LẠI, HỦY KẾ HOẠCH!"**.

## 1. Cách tạo một lệnh Hủy (Timeout)

Thông thường nhất là chúng ta đặt một giới hạn thời gian. Nếu sau 2 giây mà nhân viên chưa làm xong thì... nghỉ, không làm nữa.

```go
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
```

## 2. Tại sao Senior/Dev giỏi luôn dùng Context?

- **Hủy hàng loạt (Propagation)**: Nếu nhân viên A gọi nhân viên B, nhân viên B gọi nhân viên C. Khi Leader báo hủy A, thì B và C cũng sẽ tự động nhận được lệnh hủy qua "bộ đàm" Context này.

- **Tránh "Goroutine mồ côi"**: Nếu hàm main bị timeout mà các Goroutine vẫn chạy ngầm bên dưới, chúng sẽ ngốn sạch tài nguyên server. Context giúp dọn dẹp sạch sẽ.

## 3. Tổng kết 3 Mức độ (Cheat Sheet cho bạn)

| Công cụ | Nhiệm vụ chính | Ví dụ đời thực |
|---------|----------------|----------------|
| Goroutine (go) | Chạy đa nhiệm | Thuê thêm nhân viên làm việc. |
| WaitGroup | Đợi hoàn thành | Chốt cửa khi mọi người đã về đủ. |
| Channel | Truyền dữ liệu | Băng chuyền chuyển hàng giữa các khâu. |
| Context | Quản lý vòng đời | Bộ đàm để ra lệnh dừng/hủy công việc. |

## 💡 Bài tập cuối cùng cho bạn (Tổng hợp)

Bạn hãy thử tưởng tượng: Bạn viết một con Web Server.

- Mỗi khi có khách truy cập, bạn tạo 1 Goroutine để xử lý.

- Bạn dùng Context để nếu khách hàng ngắt kết nối giữa chừng, Goroutine đó phải dừng lại ngay.

- Bạn dùng Channel để gửi log về cho một Goroutine khác chuyên ghi file.

Bạn thấy sự phối hợp này thế nào? Nó chính là cách mà những hệ thống cực lớn như **Google, Uber, Netflix** đang vận hành bằng Go đấy!

