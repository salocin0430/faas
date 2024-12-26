package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"faas/internal/features/executions/domain/worker"

	"github.com/docker/docker/client"
	"github.com/nats-io/nats.go"
)

const (
	EXECUTIONS_SUBJECT = "executions.pending"
	WORKERS_QUEUE      = "execution-workers"
)

type NatsWorker struct {
	id           string
	natsClient   *nats.Conn
	dockerClient *client.Client
	containerMgr *worker.ContainerManager
	js           nats.JetStreamContext
}

func NewNatsWorker(natsURL string) (worker.Worker, error) {
	// Inicializar NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to NATS: %v", err)
	}

	// Inicializar Docker client
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating Docker client: %v", err)
	}

	w := &NatsWorker{
		id:           generateID(),
		natsClient:   nc,
		dockerClient: docker,
	}

	w.containerMgr = worker.NewContainerManager(docker)
	w.js = nc.JetStream()
	return w, nil
}

func (w *NatsWorker) Start(ctx context.Context) error {
	log.Printf("Worker %s starting...", w.id)

	// Suscribirse a tareas
	err := w.Subscribe()
	return err
}

func (w *NatsWorker) Subscribe() error {
	// Queue subscribe: cada mensaje va a un solo worker
	_, err := w.js.QueueSubscribe(
		EXECUTIONS_SUBJECT,
		WORKERS_QUEUE,
		w.handleExecution,
		nats.MaxConcurrent(1), // Un mensaje a la vez
		nats.ManualAck(),      // Ack manual después de procesar
		nats.DeliverNew(),     // Solo nuevos mensajes
	)
	return err
}

func (w *NatsWorker) handleExecution(msg *nats.Msg) {
	// Procesar mensaje
	// ...

	// Confirmar procesamiento
	msg.Ack()
}

func (w *NatsWorker) ProcessTask(ctx context.Context, task worker.Task) error {
	log.Printf("Worker %s processing task %s", w.id, task.ID)

	// Ejecutar función
	result, err := w.containerMgr.RunFunction(ctx, task.FunctionConfig)
	if err != nil {
		return err
	}

	// Publicar resultado
	w.publishResult(task.ID, result)
	return nil
}

func (w *NatsWorker) Stop() error {
	if w.natsClient != nil {
		w.natsClient.Close()
	}
	if w.dockerClient != nil {
		w.dockerClient.Close()
	}
	return nil
}

func (w *NatsWorker) publishResult(taskID string, result string) {
	taskResult := worker.TaskResult{
		TaskID: taskID,
		Result: result,
		Status: "completed",
	}

	data, _ := json.Marshal(taskResult)
	w.natsClient.Publish("results", data)
}

func (w *NatsWorker) publishError(taskID string, err error) {
	taskResult := worker.TaskResult{
		TaskID: taskID,
		Error:  err.Error(),
		Status: "error",
	}

	data, _ := json.Marshal(taskResult)
	w.natsClient.Publish("errors", data)
}

func generateID() string {
	return fmt.Sprintf("worker-%d", time.Now().UnixNano())
}
