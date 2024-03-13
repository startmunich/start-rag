package vector_queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type VectorQueue struct {
	basePath string
}

func New(basePath string) *VectorQueue {
	return &VectorQueue{
		basePath: basePath,
	}
}

func (v *VectorQueue) Enqueue(payload *EnqueuePayload) error {
	url := fmt.Sprintf("%s/enqueue", v.basePath)

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if _, err := http.NewRequest("POST", url, bytes.NewBuffer(body)); err != nil {
		return err
	}
	return nil
}
