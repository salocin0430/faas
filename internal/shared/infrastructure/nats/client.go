package nats

import (
	"time"

	natspkg "github.com/nats-io/nats.go"
)

const (
	FUNCTIONS_BUCKET  = "functions"
	EXECUTIONS_BUCKET = "executions"
	USERS_BUCKET      = "users"
)

func Connect(url string) (*natspkg.Conn, error) {
	return natspkg.Connect(url)
}

func CreateBuckets(js JetStreamContext) error {
	// Bucket para funciones
	_, err := js.CreateKeyValue(&natspkg.KeyValueConfig{
		Bucket:      FUNCTIONS_BUCKET,
		Description: "Functions storage",
	})
	if err != nil {
		return err
	}

	// Bucket para usuarios
	_, err = js.CreateKeyValue(&natspkg.KeyValueConfig{
		Bucket:      USERS_BUCKET,
		Description: "Users storage",
	})
	if err != nil {
		return err
	}

	// Bucket para ejecuciones
	_, err = js.CreateKeyValue(&natspkg.KeyValueConfig{
		Bucket:      EXECUTIONS_BUCKET,
		Description: "Executions storage",
	})
	if err != nil {
		return err
	}

	return nil
}

func CreateStreams(js JetStreamContext) error {
	// Crear stream persistente para ejecuciones
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
