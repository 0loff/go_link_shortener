package jwt

import (
	"fmt"
	"net/http"

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

	authCookie := &http.Cookie{
		Name:  "Auth",
		Value: token,
		Path:  "/",
	}

	out1, err := GetUserID(authCookie)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out1)

	// Output:
	// 3462f28c-1c3a-457f-8849-c5216fbf9e16
}
