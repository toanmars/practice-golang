# Mức 2: Channel - Truyền dữ liệu giữa các Goroutine

Ở **Mức 1**, chúng ta chỉ biết "đợi nhau". Ở **Mức 2**, chúng ta sẽ học cách "gửi đồ" cho nhau.

---

## 1. Channel là gì?

Hãy tưởng tượng **Channel** là một cái đường ống:

- **Gửi dữ liệu:** Một nhân viên bỏ kết quả vào đầu ống (`ch <- data`)
- **Nhận dữ liệu:** Một nhân viên khác (thường là Leader) đứng ở đầu ống kia để lấy ra (`data := <-ch`)

---

## 2. Tại sao cần Channel?

Trong Go, các Goroutine chạy riêng biệt. Nếu bạn muốn lấy dữ liệu từ một Goroutine con về hàm `main`, bạn **không nên dùng biến dùng chung** (vì dễ gây xung đột). 

**Channel là cách an toàn nhất để giao tiếp** giữa các Goroutine.

---

## 3. Code thực hành: Trình tải file (Nâng cấp)

Lần này, chúng ta không chỉ tải file, mà mỗi nhân viên sau khi tải xong sẽ **gửi dung lượng file** về cho Leader để tính tổng.

```go
package main

import (
	"fmt"
	"time"
)

func downloadFile(fileName string, resultChan chan int) {
	fmt.Printf("Đang tải %s...\n", fileName)
	time.Sleep(2 * time.Second) // Giả lập tải file

	// Giả lập dung lượng file tải được
	fileSize := len(fileName) * 100 
	
	// GỬI dữ liệu vào channel
	resultChan <- fileSize 
}

func main() {
	// 1. Tạo một đường ống (channel) chứa số nguyên (int)
	resultChan := make(chan int)

	files := []string{"video.mp4", "music.mp3", "image.jpg"}

	// 2. Chạy các Goroutine
	for _, file := range files {
		go downloadFile(file, resultChan)
	}

	totalSize := 0

	// 3. NHẬN dữ liệu từ đường ống
	// Vì có 3 file, chúng ta cần lấy dữ liệu ra 3 lần
	for i := 0; i < len(files); i++ {
		size := <-resultChan // Leader đứng đợi ở đây cho đến khi có nhân viên gửi đồ về
		totalSize += size
		fmt.Printf("Nhận được kết quả: %d KB\n", size)
	}

	fmt.Printf("--- TỔNG DUNG LƯỢNG ĐÃ TẢI: %d KB ---\n", totalSize)
}
```

---

## 4. Tại sao ở đây không cần `sync.WaitGroup`?

Đây là một điểm cực kỳ thú vị của Channel: **Nó có tính chất chặn (Blocking)**.

- Khi Leader chạy đến dòng `size := <-resultChan`, nếu chưa có nhân viên nào gửi dữ liệu về, Leader sẽ **đứng đợi** tại đó.
- Cơ chế này vô tình giúp chúng ta giữ hàm `main` không bị thoát sớm, tương tự như `wg.Wait()`.

---

## 💡 Bài tập nhỏ cho Mức 2

Trong đoạn code trên, nếu tôi chỉ chạy vòng lặp nhận dữ liệu **2 lần** (`i < 2`) thay vì 3 lần, thì chuyện gì sẽ xảy ra?

- [ ] Chương trình vẫn chạy bình thường nhưng tính thiếu 1 file.
- [ ] Chương trình bị lỗi.
- [ ] Chương trình sẽ kết thúc và bỏ qua file thứ 3.

> **Gợi ý:** Hãy nghĩ về việc nhân viên thứ 3 cố gắng gửi đồ vào ống nhưng không còn ai đứng đợi để lấy ra nữa.

## ✅ Đáp án: Lựa chọn số 3

Nếu bạn chỉ nhận **2 lần** trong khi có **3 nhân viên** gửi:

- Leader lấy xong 2 món đồ rồi... bỏ về luôn (hàm `main` kết thúc).
- Nhân viên thứ 3 đang định bỏ hàng vào ống thì "rầm", xưởng đóng cửa.
- Nếu đây là một chương trình chạy liên tục (như web server), nhân viên thứ 3 sẽ bị **treo mãi mãi** ở dòng gửi dữ liệu. 

> **⚠️ Lưu ý:** Trong lập trình, lỗi này gọi là **Goroutine Leak** (Rò rỉ bộ nhớ) – một lỗi cực kỳ nguy hiểm vì nó sẽ làm app tốn RAM dần theo thời gian.

---

# Mức 2.5: Buffered Channel (Đường ống có ngăn chứa)

Ở ví dụ trên, đường ống của chúng ta là **Unbuffered** (không có ngăn chứa). Nghĩa là: Người gửi và Người nhận phải "chạm mặt" nhau thì hàng mới đi qua được.

Bây giờ, hãy tưởng tượng xưởng của bạn có thêm một cái **Hòm thư (Buffer)**.

---

## 1. Sự khác biệt

- **Unbuffered** (`make(chan int)`): Không có hòm thư. Người gửi phải đứng đợi cho đến khi người nhận lấy hàng ra.
- **Buffered** (`make(chan int, 10)`): Có hòm thư chứa được 10 món đồ. Người gửi cứ ném đồ vào hòm rồi đi làm việc tiếp, không cần đợi người nhận có mặt ngay lúc đó. Người gửi chỉ phải đợi khi nào hòm thư bị đầy.

---

## 2. Ví dụ thực tế

Hãy tưởng tượng bạn có một "Người gửi tin nhắn" và một "Người in tin nhắn". Nếu máy in chậm hơn người gửi, chúng ta nên có một cái hòm thư để chứa tin nhắn tạm thời.

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Tạo channel có sức chứa là 3 (hòm thư có 3 ngăn)
	messenger := make(chan string, 3)

	// Người gửi ném liên tục 3 tin nhắn vào hòm
	messenger <- "Tin nhắn 1"
	messenger <- "Tin nhắn 2"
	messenger <- "Tin nhắn 3"

	fmt.Println("Đã gửi xong 3 tin nhắn vào hòm thư, không cần đợi ai lấy!")

	// Nếu gửi thêm tin thứ 4 ở đây, chương trình sẽ bị treo vì hòm đã đầy
	// messenger <- "Tin nhắn 4" 

	// Người nhận bắt đầu lấy ra
	fmt.Println("Người nhận lấy ra:", <-messenger)
	fmt.Println("Người nhận lấy ra:", <-messenger)
	fmt.Println("Người nhận lấy ra:", <-messenger)
}
```

---

## 3. Khi nào dùng cái nào?

- **Unbuffered:** Dùng khi bạn muốn sự chắc chắn. Người gửi biết chắc chắn là người nhận đã cầm được hàng thì mới đi làm việc tiếp.
- **Buffered:** Dùng khi bạn muốn tăng tốc độ. Người gửi không muốn bị lãng phí thời gian đứng đợi người nhận (vốn dĩ có thể đang bận việc khác).

---

## 💡 Bài tập tổng kết Mức 2

Giả sử bạn có **100 nhân viên** cùng làm việc, nhưng bạn chỉ có một cái hòm thư chứa được **10 món đồ**. Chuyện gì xảy ra nếu:

- Cả 100 người cùng làm xong và ném đồ vào hòm?
- Có ai bị "đứng hình" không? Và làm sao để giải quyết?

> **Gợi ý:** Hãy nghĩ về việc kết hợp `WaitGroup` (để biết khi nào 100 người xong) và `Channel` (để lấy 100 kết quả).

---