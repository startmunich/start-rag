package vector_queue

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type VectorQueue struct {
	basePath string
}

func New(basePath string) *VectorQueue {
	return &VectorQueue{
		basePath: basePath,
	}
}

func (v *VectorQueue) WaitForReady() {
	url := fmt.Sprintf("%s/ready", v.basePath)

	for {
		if res, err := http.Get(url); err == nil && res.StatusCode == 200 {
			return
		}
		time.Sleep(time.Duration(1) * time.Second)
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
	} else if res.Response == nil {
		return errors.New(fmt.Sprintf("Could not reach %s", url))
	} else if res.Response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Failed with status code %d", res.Response.StatusCode))
	}
	return nil
}
