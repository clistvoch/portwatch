package alert

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/patrickward/portwatch/internal/monitor"
)

type dingTalkPayload struct {
	MsgType  string          `json:"msgtype"`
	Text     *dingTalkText   `json:"text,omitempty"`
	Markdown *dingTalkMD     `json:"markdown,omitempty"`
}

type dingTalkText struct {
	Content string `json:"content"`
}

type dingTalkMD struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkHandler struct {
	webhookURL string
	secret     string
	msgType    string
	client     *http.Client
}

func NewDingTalkHandler(webhookURL, secret, msgType string) *DingTalkHandler {
	if msgType == "" {
		msgType = "text"
	}
	return &DingTalkHandler{
		webhookURL: webhookURL,
		secret:     secret,
		msgType:    msgType,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (h *DingTalkHandler) Handle(changes []monitor.Change) error {
	if len(changes) == 0 {
		return nil
	}

	body := formatDingTalkBody(changes)
	var payload dingTalkPayload
	payload.MsgType = h.msgType
	if h.msgType == "markdown" {
		payload.Markdown = &dingTalkMD{Title: "portwatch alert", Text: body}
	} else {
		payload.Text = &dingTalkText{Content: body}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("dingtalk: marshal payload: %w", err)
	}

	endpoint := h.webhookURL
	if h.secret != "" {
		timestamp := time.Now().UnixMilli()
		sign := dingTalkSign(h.secret, timestamp)
		endpoint = fmt.Sprintf("%s&timestamp=%d&sign=%s", endpoint, timestamp, url.QueryEscape(sign))
	}

	resp, err := h.client.Post(endpoint, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("dingtalk: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("dingtalk: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func dingTalkSign(secret string, timestamp int64) string {
	msg := fmt.Sprintf("%d\n%s", timestamp, secret)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(msg))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func formatDingTalkBody(changes []monitor.Change) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("portwatch detected %d change(s):\n", len(changes)))
	for _, c := range changes {
		buf.WriteString(fmt.Sprintf("  %s\n", c.String()))
	}
	return buf.String()
}
