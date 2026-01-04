package main

import (
	"errors"
	"fmt"
)

func CheckAge(age int) error {
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
