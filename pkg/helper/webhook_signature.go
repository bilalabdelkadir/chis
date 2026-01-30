package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

type WebhookSignature struct {
	MessageID string
	Timestamp int64
	Signature string
}

func SignWebhookPayload(msgID, secret string, body []byte) *WebhookSignature {
	return SignWebhookPayloadWithTimestamp(msgID, secret, body, time.Now().Unix())
}

func SignWebhookPayloadWithTimestamp(msgID, secret string, body []byte, timestamp int64) *WebhookSignature {
	signedContent := fmt.Sprintf("%s.%d.%s", msgID, timestamp, string(body))

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedContent))
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return &WebhookSignature{
		MessageID: msgID,
		Timestamp: timestamp,
		Signature: "v1," + sig,
	}
}
