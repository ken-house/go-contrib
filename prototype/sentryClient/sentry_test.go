package sentryClient

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitSentry(t *testing.T) {
	cfg := SentryConfig{
		Dsn:              "https://789e8b4d389e40c5994f6b09bd89d519@o435470.ingest.sentry.io/4504257001422848",
		ServerName:       "go_example",
		SampleRate:       1.0,
		AttachStacktrace: true,
		TracesSampleRate: 1.0,
		IgnoreErrors:     nil,
	}

	client, clean, err := NewSentryClient(cfg)
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	defer clean()

	_, err = os.Open("a.txt")
	if err != nil {
		client.CaptureException(err)
	}
	assert.Equal(t, err, nil)
}
