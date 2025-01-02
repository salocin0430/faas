package nats

import (
	"context"
	"encoding/json"
	"faas/internal/features/function_objects/domain/entity"
	"faas/internal/shared/infrastructure/nats"
	"fmt"
	"strings"
)

type NatsObjectRepository struct {
	kv nats.KeyValue
}

func NewNatsObjectRepository(js nats.JetStreamContext) (*NatsObjectRepository, error) {
	kv, err := js.KeyValue(nats.OBJECTS_BUCKET)
	if err != nil {
		return nil, err
	}
	return &NatsObjectRepository{kv: nats.NewKeyValueAdapter(kv)}, nil
}

func (r *NatsObjectRepository) Save(ctx context.Context, obj *entity.FunctionObject, data []byte) error {
	key := fmt.Sprintf("%s/%s", obj.FunctionID, obj.Name)

	// Save metadata
	metaKey := fmt.Sprintf("%s.meta", key)
	metaData, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err := r.kv.Put(metaKey, metaData); err != nil {
		return err
	}

	// Save data
	_, err = r.kv.Put(key, data)
	return err
}

func (r *NatsObjectRepository) Get(ctx context.Context, functionID, objectName string) (*entity.FunctionObject, []byte, error) {
	key := fmt.Sprintf("%s/%s", functionID, objectName)

	// Get metadata
	metaKey := fmt.Sprintf("%s.meta", key)
	metaEntry, err := r.kv.Get(metaKey)
	if err != nil {
		return nil, nil, err
	}

	var obj entity.FunctionObject
	if err := json.Unmarshal(metaEntry.Value(), &obj); err != nil {
		return nil, nil, err
	}

	// Get data
	dataEntry, err := r.kv.Get(key)
	if err != nil {
		return nil, nil, err
	}

	return &obj, dataEntry.Value(), nil
}

func (r *NatsObjectRepository) List(ctx context.Context, functionID string) ([]*entity.FunctionObject, error) {
	prefix := fmt.Sprintf("%s/", functionID)
	keys, err := r.kv.Keys()
	if err != nil {
		return nil, err
	}

	var objects []*entity.FunctionObject
	for _, key := range keys {
		if strings.HasPrefix(key, prefix) && strings.HasSuffix(key, ".meta") {
			entry, err := r.kv.Get(key)
			if err != nil {
				continue
			}

			var obj entity.FunctionObject
			if err := json.Unmarshal(entry.Value(), &obj); err != nil {
				continue
			}
			objects = append(objects, &obj)
		}
	}

	return objects, nil
}

func (r *NatsObjectRepository) Delete(ctx context.Context, functionID, objectName string) error {
	key := fmt.Sprintf("%s/%s", functionID, objectName)

	// Delete metadata
	metaKey := fmt.Sprintf("%s.meta", key)
	if err := r.kv.Delete(metaKey); err != nil {
		return err
	}

	// Delete data
	return r.kv.Delete(key)
}
