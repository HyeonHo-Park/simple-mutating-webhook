package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/model"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	eachCPUMin = 200
	eachCPUMax = 500
	totalCPULimit = 1000
)

type DeploymentController struct {}

//Mutate TODO: Check DecimalSI format 200 -> 200m
func (d DeploymentController) Mutate(deployment *appsv1.Deployment) ([]*model.JSONPatchEntry, error){
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

	return []*model.JSONPatchEntry{
		&model.JSONPatchEntry{
			OP:    "replace",
			Path:  "/spec/replicas",
			Value: replicaBytes,
		},
		&model.JSONPatchEntry{
			OP:    "replace",
			Path:  "/spec/template/spec/containers",
			Value: containersBytes,
		},
	}, nil
}

func checkReplicas(deployment *appsv1.Deployment) ([]byte, error) {
	maxCount := int32(3)
	if *deployment.Spec.Replicas > maxCount {
		deployment.Spec.Replicas = &maxCount
	}

	return json.Marshal(&deployment.Spec.Replicas)
}

func checkResource(deployment *appsv1.Deployment) ([]byte, error) {
	var rTotal, lTotal int64 = 0, 0
	for i, c := range deployment.Spec.Template.Spec.Containers{
		r, err := checkEachCPU(c.Resources.Requests.Cpu().Value())
		if err != nil {
			log.Error(err)
			return nil, err
		}

		rTotal, err = checkTotalCPU(rTotal + r)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		deployment.Spec.Template.Spec.Containers[i].Resources.Requests = corev1.ResourceList{
			corev1.ResourceCPU: *resource.NewQuantity(r, resource.DecimalSI),
		}

		l, err := checkEachCPU(c.Resources.Limits.Cpu().Value())
		if err != nil {
			log.Error(err)
			return nil, err
		}

		lTotal, err = checkTotalCPU(lTotal + l)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		deployment.Spec.Template.Spec.Containers[i].Resources.Limits = corev1.ResourceList{
			corev1.ResourceCPU: *resource.NewQuantity(l, resource.DecimalSI),
		}
	}

	return json.Marshal(&deployment.Spec.Template.Spec.Containers)
}

func checkEachCPU(cpu int64) (int64, error){
	switch true {
	case cpu > eachCPUMax:
		return 0, errors.New(fmt.Sprintf("usage of CPU > %dm", eachCPUMax))
	case cpu < eachCPUMin:
		return eachCPUMin, nil
	default:
		return cpu, nil
	}
}

func checkTotalCPU(cpu int64) (int64, error) {
	switch true {
	case cpu > totalCPULimit:
		return 0, errors.New(fmt.Sprintf("total usage of CPU > %dm", totalCPULimit))
	default:
		return cpu, nil
	}
}