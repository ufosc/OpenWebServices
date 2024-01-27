package authapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ufosc/OpenWebServices/pkg/websmtp"
	"io/ioutil"
	"net/http"
)

// SendVerification sends the signup verification email, where id is the
// MongoDB object ID of the pending user and email is their email address.
func (cntrl *DefaultAPIController) SendVerification(id, email string) bool {
	reqBody, err := json.Marshal(websmtp.SendRequest{
		From:    cntrl.address,
		To:      []string{email},
		Subject: "UF Open Source Club: Verify Your Email Address",
		Body:    "go to api.testing.ufosc.org/auth/verify/" + id,
	})

	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(cntrl.websmtp)
	respBody := bytes.NewBuffer(reqBody)
	resp, err := http.Post(cntrl.websmtp, "application/json", respBody)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	sb := string(body)
	fmt.Println(sb)

	return true
}
