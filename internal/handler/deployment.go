package handler

import (
	"encoding/json"
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/controller"
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"

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
		model.ErrResponse(ctx, admissionReview, err)
		return
	}

	if deployment.Namespace != NamespaceNeedToBeMutated {
		result := model.EmptyAdmissionReviewResponse(admissionReview)
		model.APIResponse(ctx, "OK", http.StatusOK, ctx.Request.Method, result)
		return
	}

	patch, err := d.controller.Mutate(&deployment)
	if err != nil {
		model.ErrResponse(ctx, admissionReview, err)
		return
	}
	// deleted log
	log.Infof("patch replicas value : %v", string(patch[0].Value))
	log.Infof("patch resources value : %v", string(patch[1].Value))

	result, err := model.SuccessAdmissionReviewResponse(admissionReview, patch)
	if err != nil {
		model.ErrResponse(ctx, admissionReview, err)
		return
	}

	// deleted log
	log.Infof("result : %v", string(result.Response.Patch))
	model.APIResponse(ctx, "OK", http.StatusOK, ctx.Request.Method, result)
}