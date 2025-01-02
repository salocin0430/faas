package nats

import (
	"time"

	natspkg "github.com/nats-io/nats.go"
)

const (
	FUNCTIONS_BUCKET  = "functions"
	EXECUTIONS_BUCKET = "executions"
	USERS_BUCKET      = "users"
	OBJECTS_BUCKET    = "function_objects"
	SECRETS_BUCKET    = "secrets"
)

func Connect(url string) (*natspkg.Conn, error) {
	return natspkg.Connect(url)
}

func CreateBuckets(js JetStreamContext) error {
	// Bucket for functions
	_, err := js.CreateKeyValue(&natspkg.KeyValueConfig{
		Bucket:      FUNCTIONS_BUCKET,
		Description: "Functions storage",
	})
	if err != nil {
		return err
	}

	// Bucket for users
	_, err = js.CreateKeyValue(&natspkg.KeyValueConfig{
		Bucket:      USERS_BUCKET,
		Description: "Users storage",
	})
	if err != nil {
		return err
	}

	// Bucket for executions
	_, err = js.CreateKeyValue(&natspkg.KeyValueConfig{
		Bucket:      EXECUTIONS_BUCKET,
		Description: "Executions storage",
	})
	if err != nil {
		return err
	}

	// Bucket for function objects
	_, err = js.CreateKeyValue(&natspkg.KeyValueConfig{
		Bucket:      OBJECTS_BUCKET,
		Description: "Function objects storage",
	})
	if err != nil {
		return err
	}

	// Bucket for secrets
	_, err = js.CreateKeyValue(&natspkg.KeyValueConfig{
		Bucket:      SECRETS_BUCKET,
		Description: "Secrets storage",
	})
	if err != nil {
		return err
	}

	return nil
}

func CreateStreams(js JetStreamContext) error {
	// Create persistent stream for executions
	_, err := js.AddStream(&natspkg.StreamConfig{
		Name:        "EXECUTIONS",
		Subjects:    []string{"executions.pending"},
		Storage:     natspkg.FileStorage,
		Retention:   natspkg.WorkQueuePolicy,
		MaxAge:      24 * time.Hour,
		Discard:     natspkg.DiscardOld,
		AllowDirect: true,
		AllowRollup: true,
	})
	return err
}
