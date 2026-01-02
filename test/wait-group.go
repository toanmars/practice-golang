package main

import (
	"fmt" // Thư viện chứa WaitGroup
	"time"
)

// Hàm giả lập việc tải một file
func downloadFile(fileName string, resultChan chan int) {

	fmt.Printf("Bắt đầu tải: %s...\n", fileName)

	// Giả lập thời gian tải mất 2 giây
	time.Sleep(2 * time.Second)

	fmt.Printf("Thành công: %s đã tải xong!\n", fileName)

	// Gửi kết quả vào channel
	resultChan <- len(fileName) * 100
}

func main() {
	resultChan := make(chan int)
	// Khai báo "bộ đếm" WaitGroup
	files := []string{"Phim_Hanh_Dong.mp4", "Nhac_Hay.mp3", "Tai_Lieu_Hoc_Tap.pdf"}

	fmt.Println("--- TRÌNH TẢI FILE BẮT ĐẦU ---")

	for _, file := range files {
		// 1. Trước khi gọi nhân viên, Leader ghi thêm 1 người vào danh sách đợi

		// Chạy hàm tải file bằng Goroutine
		go downloadFile(file, resultChan)
	}

	// 3. Chốt chặn: Main sẽ đứng yên ở đây cho đến khi bộ đếm về 0
	totalSize := 0

	for i := 0; i < len(files); i++ {
		size := <-resultChan // Leader đứng đợi ở đây cho đến khi có nhân viên gửi đồ về
		totalSize += size
		fmt.Printf("Nhận được kết quả: %d KB\n", size)
	}

	fmt.Println("--- TẤT CẢ FILE ĐÃ TẢI XONG. CHƯƠNG TRÌNH KẾT THÚC ---")

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
