package sentryClient

import (
	"os"
	"testing"

	"github.com/getsentry/sentry-go"

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

	err := InitSentry(cfg)
	if err != nil {
		assert.False(t, false, "创建sentry对象失败")
	}

	_, err = os.Open("a.txt")
	if err != nil {
		sentry.CaptureException(err)
	}
	assert.True(t, true)
}
