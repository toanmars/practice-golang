# 🧩 Interface Embedding & Composition

## 🎯 Khái niệm cơ bản

Trong Go, người ta **không dùng "Kế thừa" (Inheritance)** như Java hay C++. Thay vào đó, Go dùng **Composition** (Sự kết hợp).

> 💡 Hãy tưởng tượng bạn có các mảnh **Lego**, bạn lắp ghép chúng lại để tạo ra một thứ to lớn hơn.

---

## 1️⃣ Interface Embedding (Nhúng Interface)

Bạn có thể tạo ra một **Interface lớn** từ những **Interface nhỏ hơn**. Đây là cách thư viện chuẩn của Go được xây dựng (rất tinh tế!).

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// ReadWriter là sự kết hợp của cả hai
type ReadWriter interface {
    Reader
    Writer
}
```

### 🔍 Giải thích

- `ReadWriter` **kế thừa** tất cả các phương thức từ `Reader` và `Writer`
- Bất kỳ struct nào implement `ReadWriter` phải có cả `Read()` và `Write()`
- Đây là cách Go xây dựng các interface phức tạp từ các interface đơn giản

---

## 2️⃣ Struct Embedding (Nhúng Struct - Giả kế thừa)

Bạn có thể **nhúng một Struct này vào Struct khác** để "mượn" các thuộc tính và phương thức của nó.

```go
type User struct {
    Name string
}

func (u User) SayHi() {
    fmt.Println("Hi, I am", u.Name)
}

type Admin struct {
    User  // Nhúng trực tiếp User vào Admin
    Level int
}

func main() {
    ad := Admin{
        User: User{Name: "Tèo"},
        Level: 1,
    }
    // Admin có thể gọi luôn phương thức của User
    ad.SayHi() // Output: Hi, I am Tèo
}
```

### 🔍 Giải thích

- `Admin` **nhúng** `User` bên trong (không cần tên field)
- `Admin` tự động có tất cả thuộc tính và phương thức của `User`
- Có thể gọi `ad.SayHi()` thay vì `ad.User.SayHi()`

---

## 🤔 Tại sao Go lại làm vậy mà không dùng Kế thừa (Inheritance)?

Go muốn tránh **"Cây phả hệ"** quá phức tạp. Trong Java, một lớp con kế thừa lớp cha, lớp cha kế thừa lớp ông nội... Nếu ông nội thay đổi, cả dòng họ bị ảnh hưởng.

### ✅ Với Go:

- ✨ **Bạn chỉ cần những gì bạn cần** - Không bị ép phải kế thừa những thứ không dùng đến
- 🔧 **Bạn lắp ghép các tính năng như lắp máy tính** - Modular và linh hoạt
- 🔗 **Giúp code ít bị phụ thuộc chặt chẽ (Decoupled)** - Dễ bảo trì hơn rất nhiều

---

## 💡 Bài tập thử thách

Bây giờ hãy thử kết hợp cả **Interface** và **Embedding**:

### 📝 Yêu cầu

1. Tạo một interface `Logger` có hàm `Log(msg string)`
2. Tạo một struct `ConsoleLogger` thực hiện hàm đó bằng cách in ra màn hình
3. Tạo một struct `Server` có nhúng `ConsoleLogger` bên trong

### ❓ Câu hỏi

**Làm sao để struct `Server` có thể sử dụng hàm `Log` mà không cần phải viết lại hàm đó?**

---

## 🎓 Ứng dụng thực tế

Đây là cách chúng ta xây dựng các **Service** trong Go (như Database Service, Mail Service...) bằng cách **lắp ghép các thành phần nhỏ lại với nhau**.

```go
type DatabaseService struct {
    ConsoleLogger  // Nhúng logger
    ConnectionPool
}

// DatabaseService tự động có phương thức Log()
func (db *DatabaseService) Connect() {
    db.Log("Connecting to database...")
}
```

---

**Bạn thử viết code giải bài này xem nhé!** 🚀

---

## 📝 Giải bài tập Interface Embedding

### Cách tiếp cận cơ bản

```go
package main

import "fmt"

type Logger interface {
    Log(msg string)
}

type ConsoleLogger struct {
    message string
}

func (c *ConsoleLogger) Log(msg string) {
    fmt.Println(c.message + msg)
}

type Server struct {
    ConsoleLogger
}

func main() {
    server := Server{
        ConsoleLogger: ConsoleLogger{
            message: "Server Log: ",
        },
    }
    server.Log("Server started")
}
```

---

## 🎯 Một chút tinh chỉnh để "đẳng cấp" hơn

Trong thực tế, người ta thường **nhúng Interface** thay vì nhúng một Struct cụ thể. 

### Tại sao?

Vì nếu bạn nhúng Interface, bạn có thể thay đổi "linh hồn" của Server bất cứ lúc nào.

### Sự khác biệt

```go
type Server struct {
    Logger // Nhúng Interface thay vì Struct
}

func main() {
    // Server dùng ConsoleLogger
    s1 := Server{Logger: &ConsoleLogger{message: "Console: "}}
    s1.Log("Running")

    // Sau này bạn có FileLogger, bạn chỉ cần thay vào mà không sửa Struct Server
    // s2 := Server{Logger: &FileLogger{}} 
}
```

---