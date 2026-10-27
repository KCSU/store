// Send critical alerts to kcsu-admin.
package alerts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/kcsu/store/logging"
	alertsv1 "github.com/kcsu/store/proto/kcsu/alerts/v1"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/protobuf/proto"
)

const source = "store"

// Log and publish an alert. err and extra may be nil.
func Send(dedupKey, summary string, err error, extra map[string]string) {
	alert := &alertsv1.Alert{
		Source:         source,
		DedupKey:       dedupKey,
		OccurredAtUnix: time.Now().Unix(),
		Summary:        summary,
		Extra:          extra,
	}
	if err != nil {
		detail := err.Error()
		alert.Detail = &detail
		logging.Errorf("alert %s: %s: %s", dedupKey, summary, detail)
	} else {
		logging.Errorf("alert %s: %s", dedupKey, summary)
	}
	topic := os.Getenv("ALERTS_TOPIC")
	if topic == "" {
		logging.Warnf("ALERTS_TOPIC unset, not publishing alert %s", dedupKey)
		return
	}
	if err := publish(topic, alert); err != nil {
		logging.Errorf("publishing alert %s failed: %v", dedupKey, err)
	}
}

// Send a startup-failed alert and exit
func Fatal(summary string, err error) {
	Send("startup-failed", summary, err, nil)
	logging.Fatalf("%s: %v", summary, err)
}

func publish(topic string, alert *alertsv1.Alert) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, err := proto.Marshal(alert)
	if err != nil {
		return err
	}
	tokens, err := google.DefaultTokenSource(ctx, "https://www.googleapis.com/auth/pubsub")
	if err != nil {
		return err
	}
	body, err := json.Marshal(map[string]interface{}{
		"messages": []map[string]interface{}{{
			"data":       data,
			"attributes": map[string]string{"source": source},
		}},
	})
	if err != nil {
		return err
	}
	res, err := oauth2.NewClient(ctx, tokens).Post(
		"https://pubsub.googleapis.com/v1/"+topic+":publish", "application/json", bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("pubsub answered %s", res.Status)
	}
	return nil
}
