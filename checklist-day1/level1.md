# Mức 1: Goroutines và WaitGroup - Làm chủ đồng thời trong Go

Chào mừng bạn đến với **Mức 1**! Đây là bước đầu tiên và quan trọng nhất để làm chủ Go. Ở mức này, mục tiêu là bạn biết cách tạo ra các tác vụ chạy song song và làm sao để quản lý chúng không bị "mất tích".

**Ví dụ thực tế:** Xây dựng một chương trình tải 3 file cùng lúc.

---

## 1. Cách dùng `go func()` (Tạo tác vụ)

Trong Go, mỗi khi bạn thêm từ khóa `go` trước một hàm, bạn đang tạo ra một **Goroutine**.

- **Bình thường:** Hàm chạy xong mới đến hàm sau (Tuần tự).
- **Với `go`:** Hàm sẽ "nhảy" ra chạy riêng, và dòng code tiếp theo sẽ thực thi ngay lập tức mà không đợi hàm đó xong.

---

## 2. Cách dùng `sync.WaitGroup` (Đợi tác vụ xong)

Nếu bạn chỉ dùng `go func()` mà không có `WaitGroup`, chương trình `main` sẽ chạy vèo một cái đến cuối và kết thúc, trong khi các tác vụ con chưa kịp làm gì đã bị đóng lại.

### Quy trình 3 bước với WaitGroup:

1. **`wg.Add(n)`**: "Này Leader, có `n` nhân viên chuẩn bị đi làm đấy nhé."
2. **`wg.Done()`**: Nhân viên báo cáo: "Tôi làm xong việc của tôi rồi."
3. **`wg.Wait()`**: Leader đứng ở cửa: "Chưa đủ `n` báo cáo thì tôi chưa cho kết thúc chương trình."

---

## 3. Code thực hành: Trình tải file đa nhiệm

Hãy đọc kỹ các chú thích trong code dưới đây:

```go
package main

import (
	"fmt"
	"sync" // Thư viện chứa WaitGroup
	"time"
)

// Hàm giả lập việc tải một file
func downloadFile(fileName string, wg *sync.WaitGroup) {
	// 2. Khi hàm này chạy xong, tự động gọi Done() để báo cáo
	defer wg.Done()

	fmt.Printf("Bắt đầu tải: %s...\n", fileName)
	
	// Giả lập thời gian tải mất 2 giây
	time.Sleep(2 * time.Second)
	
	fmt.Printf("Thành công: %s đã tải xong!\n", fileName)
}

func main() {
	// Khai báo "bộ đếm" WaitGroup
	var wg sync.WaitGroup

	files := []string{"Phim_Hanh_Dong.mp4", "Nhac_Hay.mp3", "Tai_Lieu_Hoc_Tap.pdf"}

	fmt.Println("--- TRÌNH TẢI FILE BẮT ĐẦU ---")

	for _, file := range files {
		// 1. Trước khi gọi nhân viên, Leader ghi thêm 1 người vào danh sách đợi
		wg.Add(1)
		
		// Chạy hàm tải file bằng Goroutine
		go downloadFile(file, &wg)
	}

	// 3. Chốt chặn: Main sẽ đứng yên ở đây cho đến khi bộ đếm về 0
	wg.Wait()

	fmt.Println("--- TẤT CẢ FILE ĐÃ TẢI XONG. CHƯƠNG TRÌNH KẾT THÚC ---")
}
```

### Giải thích hiện tượng khi chạy:

- **Không có `go`:** Bạn sẽ mất **6 giây** (3 file × 2 giây).
- **Với đoạn code trên:** Bạn chỉ mất **2 giây** để xong cả 3 file. Vì cả 3 nhân viên cùng làm việc một lúc trên các "làn đường" khác nhau.

---

## 💡 Bài tập nhỏ cho bạn

Nếu bạn thử xóa dòng `wg.Wait()` ở cuối đi và chạy lại chương trình, bạn đoán xem kết quả sẽ hiện ra như thế nào?

- [ ] Chương trình vẫn chạy mất 2 giây rồi xong.
- [ ] Chương trình kết thúc ngay lập tức và không hiện chữ "Thành công" nào cả.

*(Bạn hãy thử đoán hoặc chạy thử code để xác nhận nhé!)*

---