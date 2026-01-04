## 🚨 Level 2: Error Handling (Xử lý lỗi) - "Đặc sản" của Go

Nếu bạn đã quen với `try-catch` trong Java hay Python, bạn sẽ thấy Go rất... lạ. 

> **Go coi Error là một giá trị** (Error is a value).

---

### 1️⃣ Tại sao Go không dùng try-catch?

Go muốn mọi lỗi đều phải được **xử lý ngay lập tức và tường minh**. 

Bạn sẽ thấy cú pháp "huyền thoại" này xuất hiện ở khắp mọi nơi:

```go
if err != nil {
    return err
}
```

---

### 2️⃣ Cách tạo một lỗi

Go cung cấp gói `errors` để tạo nhanh một thông báo lỗi.

```go
import "errors"

func Chia(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("không thể chia cho số 0")
    }
    return a / b, nil
}
```

---

### 3️⃣ Custom Error (Dành cho những bài toán phức tạp)

Vì `error` thực chất cũng chỉ là một **Interface**, nên bất cứ Struct nào có hàm `Error() string` đều được coi là một lỗi.

```go
type MyError struct {
    Code    int
    Message string
}

func (e *MyError) Error() string {
    return fmt.Sprintf("Lỗi %d: %s", e.Code, e.Message)
}
```

---

## 💡 Bài tập thực hành Error Handling

Hãy viết một hàm `CheckAge(age int) error`:

### Yêu cầu

- **Nếu `age < 0`**: Trả về một Custom Error (Struct) chứa cả số tuổi sai và thông báo "Tuổi không được âm".

- **Nếu `age < 18`**: Trả về một lỗi thông thường (`errors.New`) báo "Chưa đủ tuổi trưởng thành".

- **Nếu hợp lệ**: Trả về `nil`.

### Câu hỏi phụ

Trong hàm `main`, làm sao bạn phân biệt được lỗi trả về là lỗi "Chưa đủ tuổi" hay là cái Custom Error để in ra cái Code lỗi?

> **Gợi ý**: Hãy tìm hiểu về `errors.As` hoặc Type Assertion.

---

## 🎓 Kết luận

Phần này cực kỳ quan trọng vì nó giúp bạn viết code "sạch" và debug cực nhanh! Bạn thử sức nhé.

--- 

## Bài giải tự giải
```go 
Bài giải của tôi 
package main

import (
    "fmt"
)

func CheckAge(age int) *CustomError {
    if age < 0 {
        return &CustomError{Code: 1, Message: "Tuổi không được âm"}
    }
    if age < 18 {
        return &CustomError{Code: 2, Message: "Chưa đủ tuổi trưởng thành"}
    }

    return nil
}

type CustomError struct {
    Code    int
    Message string
}

func (e *CustomError) Error() string {
    return fmt.Sprintf("Lỗi %d: %s", e.Code, e.Message)
}

func main() {
    age := 17
    customError := CheckAge(age)
    if customError != nil {
        fmt.Println(customError)
    }
}
```
--- 
### Sửa bài giải 
Bài giải của bạn chạy rất tốt về mặt logic, nhưng ở đây có một **"cái bẫy"** về thiết kế mà hầu hết mọi người mới học Go đều mắc phải. Hãy để mình chỉ ra điểm thú vị này, vì nó sẽ giúp bạn hiểu sâu sắc về cách Go xử lý lỗi.

#### 1️⃣ Phân tích bài giải của bạn

Trong bài của bạn, hàm `CheckAge` trả về kiểu cụ thể là `*CustomError`:

```go
func CheckAge(age int) *CustomError { ... }
```

**Vấn đề là**: Trong Go, một hàm trả về lỗi thì luôn luôn nên trả về kiểu interface `error`.

**Tại sao?** 

Vì nếu bạn trả về `*CustomError`, người dùng hàm của bạn sẽ khó kết hợp với các loại lỗi khác, và đặc biệt là gặp vấn đề với việc kiểm tra `nil` (một con trỏ struct `nil` khi gán vào interface `error` đôi khi không thực sự là `nil` - đây là một lỗi rất khó chịu trong Go).

#### 2️⃣ Cách giải "Chuẩn Go" (Sử dụng `errors.As`)

Để phân biệt các loại lỗi khác nhau, Go cung cấp công cụ:
- `errors.Is` (để so sánh giá trị)
- `errors.As` (để kiểm tra kiểu dữ liệu)

```go
package main

import (
	"errors"
	"fmt"
)

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Code %d: %s", e.Code, e.Message)
}

// Trả về interface 'error' thay vì struct cụ thể
func CheckAge(age int) error {
	if age < 0 {
		return &CustomError{Code: 400, Message: "Tuổi không được âm"}
	}
	if age < 18 {
		// Trả về một lỗi thông thường bằng thư viện chuẩn
		return errors.New("chưa đủ tuổi trưởng thành")
	}
	return nil
}

func main() {
	err := CheckAge(-5)

	if err != nil {
		// Dùng errors.As để "ép kiểu" xem nó có phải CustomError không
		var cErr *CustomError
		if errors.As(err, &cErr) {
			fmt.Printf("Đây là lỗi đặc biệt! Mã lỗi: %d, Nội dung: %s\n", cErr.Code, cErr.Message)
		} else {
			fmt.Println("Đây là lỗi thông thường:", err)
		}
	}
}
```

#### 3️⃣ Tại sao dùng `errors.As` lại chuyên nghiệp hơn?

- **Tính đóng gói**: Hàm `CheckAge` có thể trả về rất nhiều loại lỗi khác nhau (lỗi DB, lỗi logic, lỗi custom...) nhưng người gọi chỉ cần biết nó là một cái `error`.

- **Mở rộng**: Bạn có thể bọc lỗi này vào trong một lỗi khác (Error Wrapping) mà vẫn dùng `errors.As` để tìm ra lỗi gốc ở sâu bên trong được.
