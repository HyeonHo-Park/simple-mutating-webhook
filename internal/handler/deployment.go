package handler

import (
	"encoding/json"
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/controller"
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/model"
	"github.com/gin-gonic/gin"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
)

const NamespaceNeedToBeMutated = "test"

type DeploymentHandler struct {
	controller *controller.DeploymentController
}

func NewDeploymentHandler() *DeploymentHandler {
	return &DeploymentHandler{
		controller: &controller.DeploymentController{},
	}
}

func (d DeploymentHandler) Mutate(ctx *gin.Context) {
	admissionReview := &admissionv1.AdmissionReview{}
	ctx.Bind(admissionReview)

	var deployment appsv1.Deployment
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, &deployment); err != nil {
		model.WriteResponse(ctx, admissionReview, err)
		return
	}

	if deployment.Namespace != NamespaceNeedToBeMutated {
		model.WriteResponse(ctx, admissionReview, nil)
		return
	}

	patch, err := d.controller.Mutate(&deployment)
	if err != nil {
		model.WriteResponse(ctx, admissionReview, patch)
		return
	}
}