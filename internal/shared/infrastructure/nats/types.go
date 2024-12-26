package nats

import natspkg "github.com/nats-io/nats.go"

// Interfaces para abstraer los tipos de NATS
type KeyValue interface {
	Get(key string) (KeyValueEntry, error)
	Put(key string, value []byte) (uint64, error)
	Delete(key string, opts ...natspkg.DeleteOpt) error
	Keys() ([]string, error)
}

type KeyValueEntry interface {
	Value() []byte
}

type JetStreamContext interface {
	CreateKeyValue(cfg *natspkg.KeyValueConfig) (natspkg.KeyValue, error)
	KeyValue(bucket string) (natspkg.KeyValue, error)
	Publish(subj string, data []byte, opts ...natspkg.PubOpt) (*natspkg.PubAck, error)
	AddStream(cfg *natspkg.StreamConfig, opts ...natspkg.JSOpt) (*natspkg.StreamInfo, error)
}

// Adaptador para convertir nats.KeyValue a nuestra interfaz
type keyValueAdapter struct {
	natsKV natspkg.KeyValue
}

func NewKeyValueAdapter(kv natspkg.KeyValue) KeyValue {
	return &keyValueAdapter{natsKV: kv}
}

func (a *keyValueAdapter) Get(key string) (KeyValueEntry, error) {
	return a.natsKV.Get(key)
}

func (a *keyValueAdapter) Put(key string, value []byte) (uint64, error) {
	return a.natsKV.Put(key, value)
}

func (a *keyValueAdapter) Delete(key string, opts ...natspkg.DeleteOpt) error {
	return a.natsKV.Delete(key, opts...)
}

func (a *keyValueAdapter) Keys() ([]string, error) {
	return a.natsKV.Keys()
}
