package auth

import (
	"fmt"
	"net/http"
)

func ExampleGetUserID() {
	authCookie := &http.Cookie{
		Name:  "Auth",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDY2OTM1MDQsIlVzZXJJRCI6IjM0NjJmMjhjLTFjM2EtNDU3Zi04ODQ5LWM1MjE2ZmJmOWUxNiJ9.AfZ55RiHueBTmvX0WyWmxA4MX4LrI7fN_D5-CbGb_xE",
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
