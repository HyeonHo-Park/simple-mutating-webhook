package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/model"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	eachCPUMin = 200
	eachCPUMax = 500
	totalCPULimit = 1000
)

type DeploymentController struct {}

func (d DeploymentController) Mutate(deployment appsv1.Deployment) ([]model.JSONPatchEntry, error){
	replicaBytes, err := checkReplicas(deployment)
	if err != nil {
		log.Errorf("marshall replicas: %v", err)
		return nil, err
	}

	containersBytes, err := checkResource(deployment)
	if err != nil {
		log.Errorf("marshall containers: %v", err)
		return nil, err
	}

	return []model.JSONPatchEntry{
		model.JSONPatchEntry{
			OP:    "replace",
			Path:  "/spec/replicas",
			Value: replicaBytes,
		},
		model.JSONPatchEntry{
			OP:    "replace",
			Path:  "/spec/template/spec/containers",
			Value: containersBytes,
		},
	}, nil
}

func checkReplicas(deployment appsv1.Deployment) ([]byte, error) {
	maxCount := int32(3)
	if *deployment.Spec.Replicas > maxCount {
		deployment.Spec.Replicas = &maxCount
	}

	return json.Marshal(&deployment.Spec.Replicas)
}

func checkResource(deployment appsv1.Deployment) ([]byte, error) {
	var rTotal, lTotal int64 = 0, 0
	for i, c := range deployment.Spec.Template.Spec.Containers{
		if req, ok := c.Resources.Requests.Cpu().AsInt64(); ok {
			r, err := checkCPU(req)
			if err != nil {
				log.Error(err)
				return nil, err
			}

			rTotal, err = checkCPU(rTotal + r)
			if err != nil {
				log.Error(err)
				return nil, err
			}

			deployment.Spec.Template.Spec.Containers[i].Resources.Requests.Cpu().Set(r)
		}

		if limit, ok := c.Resources.Limits.Cpu().AsInt64(); ok {
			l, err := checkCPU(limit)
			if err != nil {
				log.Error(err)
				return nil, err
			}

			lTotal, err = checkCPU(lTotal + l)
			if err != nil {
				log.Error(err)
				return nil, err
			}

			deployment.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu().Set(l)
		}
	}

	return json.Marshal(&deployment.Spec.Template.Spec.Containers)
}

func checkCPU(cpu int64) (int64, error){
	switch true {
	case cpu > totalCPULimit:
		return 0, errors.New(fmt.Sprintf("total usage of CPU > %dm", totalCPULimit))
	case cpu > eachCPUMax:
		return 0, errors.New(fmt.Sprintf("usage of CPU > %dm", eachCPUMax))
	case cpu < eachCPUMin:
		return eachCPUMin, nil
	default:
		return cpu, nil
	}
}