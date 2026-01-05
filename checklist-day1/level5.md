# 📁 Cẩm Nang Xử Lý Concurrency trong Golang

## 1. Race Condition (Điều kiện đua)
**Định nghĩa:** Xảy ra khi nhiều goroutines truy cập chung một tài nguyên và có ít nhất một thao tác ghi cùng lúc.

* **Kiểm tra nhanh:** `go run -race main.go`
* **Ví dụ:** Hai Goroutine cùng `count++` dẫn đến giá trị cuối cùng bị thiếu hụt.

---
Example code 
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var balance = 0
	var wg sync.WaitGroup

	// Chạy 1000 goroutines, mỗi cái nạp 1 đồng
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			balance = balance + 1 // Race condition xảy ra ở đây
		}()
	}

	wg.Wait()
	fmt.Printf("Số dư cuối cùng: %d\n", balance)
}
```
Tại sao kết quả thường ít hơn 1000? Phép tính balance = balance + 1 thực chất gồm 3 bước ở cấp độ CPU:

 - Đọc giá trị hiện tại của balance (ví dụ là 10).

 - Cộng thêm 1 vào giá trị đó (10 + 1 = 11).

 - Ghi giá trị mới (11) quay lại biến balance.

Nếu hai goroutines cùng thực hiện bước 1 khi balance đang là 10, cả hai đều tính ra 11 và cùng ghi đè số 11 vào biến. Kết quả là ta mất đi một lần tăng giá trị.

## 2. Giải pháp: Mutex vs Atomic

| Đặc điểm | Mutex (`sync.Mutex`) | Atomic (`sync/atomic`) |
| :--- | :--- | :--- |
| **Cơ chế** | Khóa phần mềm (Blocking) | Thao tác CPU (Non-blocking) |
| **Tốc độ** | Chậm hơn (do quản lý sleep/wakeup) | Cực nhanh (Hardware level) |
| **Sử dụng** | Bảo vệ logic phức tạp, struct, map | Biến số đơn giản (int, bool, pointer) |

### Mutex (The Guard)
Mutex giống như một chiếc vé duy nhất để vào phòng. Nếu Goroutine A đang giữ vé, Goroutine B phải đứng đợi ở cửa (trạng thái Blocked). Khi A trả vé, Runtime của Go mới đánh thức B dậy để vào. Quá trình "ngủ" và "thức dậy" này tốn chi phí quản lý của hệ điều hành.

### Atomic (The Specialist)
Atomic không bắt ai phải đợi. Nó sử dụng các lệnh đặc biệt của CPU (như LOCK XADD). Nếu có 2 Goroutine cùng tác động, CPU sẽ xếp hàng chúng ở mức vi mạch. Không có Goroutine nào bị đưa vào trạng thái "ngủ", vì vậy nó cực kỳ nhanh.

### Ví dụ Atomic (Go 1.19+):
```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic" // Import gói atomic
)

func main() {
	var balance int64 = 0 // Phải dùng kiểu dữ liệu cố định như int32 hoặc int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Thay vì balance = balance + 1
			atomic.AddInt64(&balance, 1) 
		}()
	}

	wg.Wait()
	// Sử dụng Load để đọc dữ liệu an toàn
	fmt.Printf("Số dư cuối cùng: %d\n", atomic.LoadInt64(&balance))
}
```
### Các thao tác phổ biến với Atomic

| Thao tác | Hàm tương ứng (ví dụ với int64) | Ý nghĩa |
| :--- | :--- | :--- |
| Cộng | `atomic.AddInt64(&addr, delta)` | Cộng một giá trị vào biến |
| Ghi | `atomic.StoreInt64(&addr, val)` | Gán giá trị mới cho biến |
| Đọc | `atomic.LoadInt64(&addr)` | Đọc giá trị hiện tại một cách an toàn |
| Tráo đổi | `atomic.SwapInt64(&addr, new)` | Gán giá trị mới và trả về giá trị cũ |
| So sánh & Đổi | `atomic.CompareAndSwapInt64(...)` | Chỉ đổi nếu giá trị hiện tại bằng giá trị mong đợi (CAS) |

### Cơ chế Compare-and-Swap (CAS)
Đây là "linh hồn" của lập trình không khóa (lock-free). Nó chỉ cập nhật giá trị nếu giá trị hiện tại đúng như ta nghĩ.

Ví dụ: Bạn muốn cập nhật số dư nhưng chỉ khi số dư hiện tại đang là 100.
```go
var balance int64 = 100
atomic.CompareAndSwapInt64(&balance, 100, 200) // Chỉ đổi nếu balance == 100
```

### Khi nào dùng Atomic, khi nào dùng Mutex?
Dùng Atomic khi:

- Bạn chỉ cần thao tác trên một biến đơn lẻ (biến đếm, cờ hiệu, trạng thái).
- Cần hiệu năng cực cao và biến đó là các kiểu số cơ bản (int32, int64, uint32,...).

Dùng Mutex khi:

- Bạn cần bảo vệ một đoạn code phức tạp (ví dụ: vừa đọc map, vừa ghi struct, vừa check điều kiện).
- Thao tác trên các kiểu dữ liệu không được atomic hỗ trợ trực tiếp như map, slice, hoặc string.

## 3. Tối ưu hóa Lock Contention (Hệ thống lớn)
Khi nhiều Goroutine cùng đợi một Mutex, hệ thống bị nghẽn (giống như kẹt xe tại trạm thu phí). Để xử lý, ta có các chiến lược:

- Giảm phạm vi khóa (Reduce Lock Granularity): Chỉ Lock đoạn code thật sự cần thiết, không Lock cả hàm dài.

- Phân mảnh khóa (Lock Sharding): Thay vì 1 Map dùng 1 Mutex, ta chia Map thành 32 phần nhỏ, mỗi phần 1 Mutex. Điều này giảm tỷ lệ các Goroutine đụng độ nhau.

- Sử dụng sync.RWMutex: Nếu hệ thống đọc nhiều, ghi ít. RWMutex cho phép nhiều Goroutine cùng đọc một lúc, chỉ khóa hoàn toàn khi có thao tác ghi.

- Copy-on-Write: Thay vì sửa dữ liệu cũ, ta tạo một bản copy mới, cập nhật dữ liệu trên đó rồi tráo đổi (swap) con trỏ bằng atomic.

Khi pprof báo hiệu Mutex đang là nút thắt cổ chai, hãy áp dụng các kỹ thuật sau:

### A. Sử dụng Read-Write Mutex (sync.RWMutex)
Nếu bạn có 90% thao tác là đọc, dùng RWMutex để các luồng đọc không chặn nhau.

```go
var mu sync.RWMutex
mu.RLock()   // Nhiều luồng vào đây cùng lúc được
// ... read data ...
mu.RUnlock()
```

### B. Lock Sharding (Phân mảnh khóa)
Thay vì dùng 1 Mutex cho toàn bộ Database/Map, hãy chia nhỏ ra:

Bucket[0] - Mutex[0]

Bucket[1] - Mutex[1] Giúp giảm xác suất 2 Goroutine cần cùng 1 khóa xuống nhiều lần.

### C. Giảm Critical Section
Di chuyển các thao tác nặng (I/O, tính toán toán học) ra ngoài khối Lock() / Unlock().

```go

// TỐT
data := calculateComplexStuff() // Tính trước
mu.Lock()
result = data                   // Chỉ Lock khi gán
mu.Unlock()
```

### D. Ưu tiên Atomic cho cờ hiệu (Flags)
Dùng atomic.Value hoặc atomic.Bool cho các biến trạng thái (is_closed, is_running) để tránh dùng Mutex không cần thiết.

Ghi chú: Luôn ưu tiên sự đơn giản của Mutex trước, chỉ tối ưu bằng Atomic hoặc Sharding khi pprof chỉ ra vấn đề thực sự về hiệu năng.

## 4. Tối ưu hóa Lock Contention (Hệ thống lớn)
Khi pprof báo hiệu Mutex đang là nút thắt cổ chai, hãy áp dụng các kỹ thuật sau:

### A. Sử dụng Read-Write Mutex (sync.RWMutex)
Nếu bạn có 90% thao tác là đọc, dùng RWMutex để các luồng đọc không chặn nhau.

```go
var mu sync.RWMutex
mu.RLock()   // Nhiều luồng vào đây cùng lúc được
// ... read data ...
mu.RUnlock()
```

### B. Lock Sharding (Phân mảnh khóa)
Thay vì dùng 1 Mutex cho toàn bộ Database/Map, hãy chia nhỏ ra:

Bucket[0] - Mutex[0]

Bucket[1] - Mutex[1] Giúp giảm xác suất 2 Goroutine cần cùng 1 khóa xuống nhiều lần.

### C. Giảm Critical Section
Di chuyển các thao tác nặng (I/O, tính toán toán học) ra ngoài khối Lock() / Unlock().

```go
// TỐT
data := calculateComplexStuff() // Tính trước
mu.Lock()
result = data                   // Chỉ Lock khi gán
mu.Unlock()
```

### D. Ưu tiên Atomic cho cờ hiệu (Flags)
Dùng ***atomic.Value*** hoặc ***atomic.Bool*** cho các biến trạng thái ***is_closed***, ***is_running*** để tránh dùng Mutex không cần thiết.

Ghi chú: Luôn ưu tiên sự đơn giản của Mutex trước, chỉ tối ưu bằng Atomic hoặc Sharding khi pprof chỉ ra vấn đề thực sự về hiệu năng.