# 🚀 Practice Golang

Chào mừng bạn đến với thế giới của Go! Tài liệu này giúp bạn hiểu cách Go vận hành đa nhiệm và cách dùng các công cụ "soi lỗi" để ứng dụng của bạn luôn chạy mượt mà.

---

## 📋 Mục lục

1. [Bộ Ba Quyền Lực: GOMAXPROCS - Goroutine - Channel](#1-bộ-ba-quyền-lực-gomaxprocs---goroutine---channel)
2. [Công cụ "Soi Lỗi" - pprof](#2-công-cụ-soi-lỗi---pprof)
3. [Công cụ "Quay Phim" - go tool trace](#3-công-cụ-quay-phim---go-tool-trace)
4. [Checklist Học Tập](#4-checklist-học-tập)

---

## 1. Bộ Ba Quyền Lực: GOMAXPROCS - Goroutine - Channel

> 💡 **Ví dụ:** Hãy tưởng tượng bạn đang quản lý một xưởng may

### 🔧 GOMAXPROCS: Số lượng máy may (CPU)

Đây là cấu hình quyết định xem bạn có bao nhiêu "máy may" có thể hoạt động cùng một lúc.

- **Mặc định:** Go sẽ tự đếm xem máy tính của bạn có bao nhiêu nhân (core) và cấp bấy nhiêu máy may.

> ⚠️ **Lưu ý:** Nếu bạn chạy Go trong Docker/Kubernetes, đôi khi Go "đếm nhầm" số nhân của máy chủ vật lý thay vì số nhân được cấp cho Container. Điều này làm app chạy giật lag.

**✅ Giải pháp:** Thêm dòng này để Go tự động chỉnh cho đúng:

```go
import _ "go.uber.org/automaxprocs"
```

### 👷 Goroutine: Công nhân (Task)

Mỗi khi bạn dùng từ khóa `go`, bạn đang gọi một "công nhân" ra làm việc.

- ✨ Công nhân Go rất nhẹ, bạn có thể gọi hàng ngàn người mà không sợ tốn nhiều RAM
- 🔄 Họ làm việc độc lập và không làm gián đoạn công việc của người khác

### 🔗 WaitGroup & Channel: Phối hợp công việc

| Công cụ | Ví dụ | Khi nào dùng |
|---------|-------|--------------|
| **WaitGroup** | Cái chốt cửa | Dùng khi bạn muốn đợi tất cả công nhân làm xong việc mới được đóng cửa xưởng |
| **Channel** | Băng chuyền | Dùng khi công nhân khâu A làm xong muốn chuyển sản phẩm sang cho công nhân khâu B |

---

## 2. Công cụ "Soi Lỗi" - pprof

> 🔍 **Chụp ảnh X-quang ứng dụng**

Khi app của bạn chạy chậm hoặc ngốn nhiều tài nguyên, đừng đoán mò. Hãy dùng `pprof`.

### 📝 Cách kích hoạt (Rất dễ)

Chỉ cần thêm dòng này vào đầu file `main.go`:

```go
import _ "net/http/pprof"
import "net/http"

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    // ... code của bạn ...
}
```

### 💻 Ví dụ soi lỗi CPU

Nếu máy tính bỗng dưng nóng lên khi chạy app, hãy gõ lệnh này ở terminal:

```bash
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

Sau 30 giây, bạn gõ lệnh:

- **`top`**: Xem 10 hàm nào đang ngốn CPU nhất
- **`web`**: Nó sẽ vẽ một sơ đồ, cái nào to và màu đỏ chính là nơi bị lỗi

---

## 3. Công cụ "Quay Phim" - go tool trace

> 🎬 **Theo dõi dòng thời gian**

Nếu `pprof` là chụp ảnh, thì `trace` là quay phim. Nó cho bạn thấy từng tích tắc các Goroutine đang làm gì.

### 📹 Cách dùng

**Bước 1:** Lấy dữ liệu

```bash
curl http://localhost:6060/debug/pprof/trace?seconds=5 > trace.out
```

**Bước 2:** Xem kết quả

```bash
go tool trace trace.out
```

### 👀 Bạn sẽ thấy gì?

- ⚔️ Thấy các Goroutine đang tranh giành nhau hay đang đợi nhau
- 🗑️ Thấy khi nào hệ thống tự động dọn rác (Garbage Collection) làm dừng app

---

## 4. Checklist Học Tập

> 📚 **Học theo lộ trình**

- [ ] **Mức 1:** Biết dùng `go func()` và `sync.WaitGroup` để chạy nhiều việc cùng lúc
- [ ] **Mức 2:** Biết dùng `channel` để truyền dữ liệu qua lại giữa các hàm
- [ ] **Mức 3:** Biết cài đặt `automaxprocs` khi đưa app lên Docker/Kubernetes để tránh bị chậm
- [ ] **Mức 4:** Khi app chậm, biết mở `pprof` lên để tìm xem dòng code nào là "thủ phạm"

---

## 💡 Lời khuyên

Đừng quá lo lắng về việc phải hiểu hết mọi thứ ngay lập tức. Hãy bắt đầu bằng việc dùng `pprof` để soi thử một ứng dụng nhỏ bạn đang viết, bạn sẽ thấy nó rất thú vị!

