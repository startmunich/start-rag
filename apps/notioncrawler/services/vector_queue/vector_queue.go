package vector_queue

import (
	"bytes"
	"encoding/json"
	"errors"
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

	if res, err := http.NewRequest("POST", url, bytes.NewBuffer(body)); err != nil {
		return err
	} else if res.Response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Failed with status code %d", res.Response.StatusCode))
	}
	return nil
}
