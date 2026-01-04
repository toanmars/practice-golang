# 🚗 Interface - "Tấm bằng lái xe" của Go

## 🎯 Khái niệm cơ bản

Hãy tưởng tượng thế này:

- **Struct**: Là cái xe cụ thể (Xe Honda, Xe Tesla, Xe Đạp)
- **Interface**: Là cái "Bằng lái xe". Nó chỉ quy định: *"Ai muốn lái thì phải biết Bành trướng, Phanh và Rẽ"*

### ✨ Implicit Implementation (Triển khai ngầm định)

Go có một cơ chế rất hay gọi là **Implicit Implementation**. Bạn không cần viết `implements InterfaceName`. Chỉ cần Struct của bạn có đủ các hàm mà Interface yêu cầu, Go sẽ tự hiểu là *"Ồ, anh này có bằng lái rồi!"*.

---

## 1️⃣ Cách định nghĩa và sử dụng

```go
package main

import "fmt"

// 1. Định nghĩa Interface (Bộ quy chuẩn)
type Animal interface {
    Speak() string
}

// 2. Các Struct cụ thể
type Dog struct{}
func (d Dog) Speak() string { return "Gâu Gâu" }

type Cat struct{}
func (c Cat) Speak() string { return "Meo Meo" }

// 3. Hàm dùng Interface (Cực kỳ linh hoạt)
func MakeAnimalSpeak(a Animal) {
    fmt.Println(a.Speak())
}

func main() {
    cho := Dog{}
    meo := Cat{}

    MakeAnimalSpeak(cho) // Gâu Gâu
    MakeAnimalSpeak(meo) // Meo Meo
}
```

---

## 2️⃣ Tại sao cái này lại quan trọng?

Nếu không có Interface, mỗi lần bạn thêm một con vật mới, bạn lại phải viết một hàm `MakeDogSpeak`, `MakeCatSpeak`... **Rất mệt!** 

Với Interface, bạn chỉ viết hàm đó **đúng 1 lần**.

### 🌟 Ứng dụng thực tế

#### 🗄️ Database
Bạn tạo interface `Repository`. Hôm nay dùng MySQL, mai đổi sang MongoDB chỉ cần đổi Struct, code logic chính không đổi.

#### 🧪 Unit Test
Bạn có thể tạo một "Database giả" (Mock) thỏa mãn Interface để test mà không cần kết nối thật.

---

## 3️⃣ Interface rỗng `interface{}` (Hoặc `any` trong bản mới)

Đây là một **"chiếc túi thần kỳ"**. Nó không yêu cầu hàm nào cả, nên mọi thứ đều thỏa mãn nó.

Bạn dùng nó khi bạn chưa biết dữ liệu truyền vào là gì (Số, Chữ, hay một Struct phức tạp).

```go
func PrintAnything(value any) {
    fmt.Println(value)
}

// Có thể truyền bất kỳ kiểu dữ liệu nào
PrintAnything(42)
PrintAnything("Hello")
PrintAnything(Dog{})
```

---

## 💡 Bài tập nhỏ để bạn "thấm" Interface

### 📝 Đề bài

Giả sử bạn đang viết **hệ thống thanh toán** cho shop online:

- Bạn có 2 phương thức: 
  - **Thanh toán bằng Thẻ** (Visa)
  - **Thanh toán bằng Ví** (Momo)
- Cả hai đều phải có hàm `Pay(amount int)`

### ❓ Câu hỏi

1. Bạn sẽ thiết kế Interface `PaymentMethod` như thế nào?
2. Nếu sau này sếp muốn thêm **"Thanh toán bằng Bitcoin"**, bạn có phải sửa code của cái hàm `ProcessPayment(p PaymentMethod)` không?

---

## 🎓 Tiếp theo

Khi bạn giải xong bài này, tôi sẽ chỉ cho bạn một tuyệt chiêu của Go: 

**Interface Embedding** (Nhúng interface này vào interface khác) - cách Go xây dựng nên những thư viện chuẩn cực mạnh!

---