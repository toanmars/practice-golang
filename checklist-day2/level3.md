# Level 3: Array, Slice & Map

## 🎯 Chủ đề: Array và Slice (Tại sao Slice là "vua" trong Go?)

Nhiều người mới học Go thường nhầm lẫn giữa **Array** và **Slice**. Nhưng trong thực tế, **99% thời gian chúng ta dùng Slice**.

---

### 1️⃣ Sự khác biệt "chí mạng"

#### Array

Cố định kích thước. `[3]int` và `[4]int` là hai kiểu dữ liệu hoàn toàn khác nhau. Bạn không thể thay đổi độ dài của nó sau khi khai báo.

#### Slice

Linh hoạt. Nó thực chất là một cái **"cửa sổ"** soi vào một mảng bên dưới.

---

### 2️⃣ Cấu trúc của một Slice (Bạn cần nhớ cái này để đi phỏng vấn)

Một Slice gồm **3 phần**:

- **Pointer**: Trỏ đến mảng dữ liệu thực sự.
- **Length (`len`)**: Số lượng phần tử hiện có trong slice.
- **Capacity (`cap`)**: Số lượng phần tử tối đa mà mảng bên dưới có thể chứa (trước khi phải cấp phát mảng mới to hơn).

---

## 💡 Bài tập thực hành về Slice

Hãy chạy đoạn code sau trong đầu (hoặc chạy thử) và dự đoán kết quả:

```go
package main

import "fmt"

func main() {
    a := []int{1, 2, 3}
    b := a
    b[0] = 100

    fmt.Println("a:", a)
    fmt.Println("b:", b)
}
```

### Câu hỏi

1. `a[0]` sẽ là bao nhiêu? Tại sao?
2. Chuyện gì xảy ra nếu mình dùng lệnh `b = append(b, 4)`? Liệu `a` có thấy số `4` đó không?

> Phần này sẽ giúp bạn hiểu về **Memory Management** trong Go cực kỳ tốt! Bạn thử giải đáp xem sao.

---

## 📋 Câu trả lời của bạn

- `a[0]` sẽ là `100`
- Nếu dùng `b = append(b, 4)` thì sẽ đưa số `4` vào cuối mảng `b`
- `a` sẽ thấy số `4` đó luôn, vì `b` đang copy cả địa chỉ và value của `a` rồi (`b:=a`)

---

## ✅ Phân tích câu trả lời

Câu trả lời của bạn cho thấy bạn đã nắm được một phần quan trọng, nhưng có một **"cú lừa"** kinh điển về Slice trong Go mà bạn cần lưu ý. Hãy cùng phân tích từng ý nhé:

### ✅ `a[0]` sẽ là 100: CHÍNH XÁC

Vì `a` và `b` cùng trỏ chung vào một mảng dữ liệu (**Backing Array**) nằm ở dưới. Sửa nhà này thì nhà kia cũng thấy.

### ❌ `a` sẽ thấy số 4 đó luôn: SAI

**(Đây chính là cái bẫy!)**

---

## 🔍 Tại sao `a` không thấy số 4?

Dù `a` và `b` trỏ chung vào một mảng, nhưng chúng là **hai cái cửa sổ riêng biệt** (`len` và `cap` riêng).

1. Khi bạn tạo `a := []int{1, 2, 3}`, cả `len` và `cap` của `a` đều là `3`.

2. Khi bạn `append(b, 4)`, Go nhận thấy cái mảng hiện tại (`capacity = 3`) đã hết chỗ để chứa thêm phần tử thứ 4.

3. **Hành động của Go**: Nó sẽ âm thầm tạo ra một mảng mới to hơn (thường là gấp đôi, `capacity = 6`), copy 3 phần tử cũ sang, rồi mới thêm số `4` vào.

4. **Kết quả**: Lúc này `b` đã trỏ sang một địa chỉ mới hoàn toàn, còn `a` vẫn trỏ về mảng cũ có 3 phần tử.

---

## 🎓 Thử thách thực tế: "Cắt lát" (Slicing)

Hãy xem ví dụ này, nó rất hay xuất hiện trong các bài kiểm tra logic:

```go
package main

import "fmt"

func main() {
    parent := []int{0, 1, 2, 3, 4, 5}
    child := parent[1:4] // Lấy từ index 1 đến 3 (không lấy 4)

    fmt.Println("child:", child) // [1 2 3]
    fmt.Println("len:", len(child), "cap:", cap(child))

    child[0] = 99
    fmt.Println("parent sau khi sửa child:", parent)
}
```

### Câu hỏi dành cho bạn

1. Theo bạn `cap` của `child` lúc này là bao nhiêu? 
   > **Gợi ý**: Tính từ vị trí nó bắt đầu cho đến hết mảng gốc.

2. Khi sửa `child[0] = 99`, thì `parent` sẽ bị thay đổi ở index nào? Giá trị đó là bao nhiêu?

---

## 🗺️ Chủ đề bổ trợ: Map (Bảng băm)

Sau khi xong Slice, chúng ta không thể bỏ qua **Map** – cách Go lưu trữ dữ liệu dạng **Key-Value**.

### Khai báo

```go
m := make(map[string]int)
```

### ⚠️ Lưu ý quan trọng

Map trong Go **không đảm bảo thứ tự**. Nếu bạn duyệt map 2 lần, kết quả có thể hiện ra theo thứ tự khác nhau.

### Kiểm tra tồn tại

Trong Go, chúng ta dùng cú pháp **"comma ok"**:

```go
val, ok := m["key"]
if ok {
    fmt.Println("Tìm thấy:", val)
}
```

---

## 🚀 Thử sức

Bạn thử giải đáp 2 câu hỏi về Slice ở trên nhé, phần `cap` của Slice là kiến thức cực kỳ quan trọng để tối ưu bộ nhớ đấy!

### Câu trả lời của bạn
1. Theo bạn cap của child lúc này là bao nhiêu? (Gợi ý: Tính từ vị trí nó bắt đầu cho đến hết mảng gốc). -> 6

2. Khi sửa child[0] = 99, thì parent sẽ bị thay đổi ở index nào? Giá trị đó là bao nhiêu? -> index 0 giá trị 99

### Phân tích câu trả lời

### 1. Phân tích Capacity (Sức chứa)
Mảng gốc parent có 6 phần tử: [0, 1, 2, 3, 4, 5].

child := parent[1:4] bắt đầu từ index 1 của parent.

Công thức: cap của một slice mới sẽ tính từ vị trí bắt đầu của nó cho đến hết mảng gốc bên dưới.

Kết quả: Từ index 1 đến index 5 là có 5 phần tử. Vậy cap(child) = 5. (Số 6 là tổng của parent, nhưng child đã bỏ qua phần tử đầu tiên rồi).

### 2. Sửa child[0] = 99 thì parent thay đổi ở đâu?
Đây là chỗ dễ nhầm nhất:

child[0] trỏ vào phần tử đầu tiên của chính nó.

Vì child bắt đầu từ parent[1], nên child[0] chính là parent[1].

Kết quả: parent sẽ thay đổi ở index 1, giá trị là 99.

Mảng parent lúc này sẽ là: [0, 99, 2, 3, 4, 5].
### 
