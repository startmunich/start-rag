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

	if res, err := http.Post(url, "application/json", bytes.NewBuffer(body)); err != nil {
		return err
	} else if res == nil {
		return errors.New(fmt.Sprintf("Could not reach %s", url))
	} else if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Failed with status code %d", res.StatusCode))
	}
	return nil
}

func (v *VectorQueue) PurgeQueue() error {
	url := fmt.Sprintf("%s/empty_redis", v.basePath)

	if res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{})); err != nil {
		return err
	} else if res == nil {
		return errors.New(fmt.Sprintf("Could not reach %s", url))
	} else if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Failed with status code %d", res.StatusCode))
	}
	return nil
}

func (v *VectorQueue) PurgeVectorDb() error {
	url := fmt.Sprintf("%s/empty_qdrant", v.basePath)

	if res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte{})); err != nil {
		return err
	} else if res == nil {
		return errors.New(fmt.Sprintf("Could not reach %s", url))
	} else if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Failed with status code %d", res.StatusCode))
	}
	return nil
}
