package jwt

import (
	"fmt"

	"github.com/google/uuid"
)

func ExampleGetUserID() {
	userID, err := uuid.Parse("3462f28c-1c3a-457f-8849-c5216fbf9e16")
	if err != nil {
		fmt.Println(err)
	}

	token, err := BuildJWTString(userID)
	if err != nil {
		fmt.Println(err)
	}

	out1, err := GetUserID(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out1)

	// Output:
	// 3462f28c-1c3a-457f-8849-c5216fbf9e16
}
