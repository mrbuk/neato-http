package neato

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// NeatoRobot
type NeatoRobot struct {
	SerialNumber string
	Secret       string
}

var (
	t = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    5 * time.Second,
		DisableCompression: true,
	}

	client = &http.Client{Transport: t}
)

// signCommand creates the signature used to authenticate against the
// nucleo API
func (n *NeatoRobot) signCommand(date, body string) string {
	// ensure the robot serial is lowercase as uppercase won't be found
	message := fmt.Sprintf("%s\n%s\n%s", strings.ToLower(n.SerialNumber), date, body)

	hash := hmac.New(sha256.New, []byte(n.Secret))
	hash.Write([]byte(message))

	signature := hex.EncodeToString(hash.Sum(nil))

	return signature
}

// formatRFC1123Date formats a given to the specific version of RFC1123
// required by the nucleo API authentication header.
func formatRFC1123Date(t time.Time) string {
	// build date in format "Mon, 02 Jan 2006 15:04:05 MST"
	// ensure to be in UTC
	raw := t.UTC().Format(time.RFC1123)

	// replace UTC with GMT as neato only understands "GMT"
	rfc1123UTCDate := strings.ReplaceAll(raw, "UTC", "GMT")

	return rfc1123UTCDate
}

// sendCommand communicates with the robot by talking to the
// nucleo API. It creates the proper authentication header using
// the robot secret
func (n *NeatoRobot) sendCommand(command string) ([]byte, error) {
	now := time.Now()
	date := formatRFC1123Date(now)

	signature := n.signCommand(date, command)

	url := fmt.Sprintf("https://nucleo.neatocloud.com:4443/vendors/neato/robots/%s/messages", n.SerialNumber)

	req, err := http.NewRequest("POST", url, strings.NewReader(command))
	if err != nil {
		return nil, err
	}
	// Add required Headers. Ensure that the date
	req.Header.Add("Accept", "application/vnd.neato.nucleo.v1")
	req.Header.Add("Authorization", fmt.Sprintf("NEATOAPP %s", signature))
	req.Header.Add("Date", date)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// assume that everything other than 200 is actually an error
	if resp.StatusCode != http.StatusOK {
		message, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(message))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return []byte(body), nil
}

// unmarshal
func unmarshal(b []byte, ref interface{}) error {
	err := json.Unmarshal([]byte(b), ref)
	if err != nil {
		return err
	}
	return nil
}

// GetState return capabilities and the actual state of the robot.
func (n *NeatoRobot) GetState() (*StateResponse, error) {
	command := `{"reqId":"77", "cmd":"getRobotState"}`
	b, err := n.sendCommand(command)

	if err != nil {
		return nil, err
	}

	var r StateResponse
	if err = unmarshal(b, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// IsCleaning reports true when the robot is clean, false if it is not
// an error is retured in case the cleaning state could be determined.
// Error is not returned in case the robot is in error state itself.
func (n *NeatoRobot) IsCleaning() (bool, error) {
	state, err := n.GetState()
	if err != nil {
		return false, err
	}

	if state.Action == StateBusy {
		return true, nil
	}

	return false, nil
}

// StartCleaning start map based house cleaning program
func (n *NeatoRobot) StartCleaning(withoutMap bool) (*StandardResponse, error) {
	var command string
	if withoutMap {
		command = `{"reqId":"77", "cmd":"startCleaning", "params": {"category": 2, "mode": 2, "navigationMode": 3}}`
	} else {
		command = `{"reqId":"77", "cmd":"startCleaning", "params": {"category": 4, "mode": 2, "navigationMode": 3}}`
	}

	b, err := n.sendCommand(command)

	if err != nil {
		return nil, err
	}

	var r StandardResponse
	if err = unmarshal(b, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// StartCleaning start map based house cleaning program
func (n *NeatoRobot) ResumeCleaning() (*StandardResponse, error) {
	command := `{"reqId": "77", "cmd": "resumeCleaning"}`
	b, err := n.sendCommand(command)

	if err != nil {
		return nil, err
	}

	var r StandardResponse
	if err = unmarshal(b, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// SendToBase moves robots back to base
func (n *NeatoRobot) SendToBase() (*StandardResponse, error) {
	command := `{"reqId":"77", "cmd":"sendToBase"}`

	b, err := n.sendCommand(command)

	if err != nil {
		return nil, err
	}

	var r StandardResponse
	if err = unmarshal(b, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// StopCleaning issues a stop command to the neato robot
func (n *NeatoRobot) StopCleaning() (*StandardResponse, error) {
	command := `{"reqId":"77", "cmd":"stopCleaning"}`

	b, err := n.sendCommand(command)

	if err != nil {
		return nil, err
	}

	var r StandardResponse
	if err = unmarshal(b, &r); err != nil {
		return nil, err
	}

	return &r, nil
}
