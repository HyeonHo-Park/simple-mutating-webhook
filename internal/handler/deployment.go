package handler

import (
	"encoding/json"
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/controller"
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/model"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (d DeploymentHandler) Mutate(ctx *gin.Context)  {
	admissionReview := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
	}

	ctx.Bind(admissionReview)
	//ctx.Bind(admissionReview.Request)

	var deployment appsv1.Deployment
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, deployment); err != nil {
		model.ErrResponse(ctx, admissionReview, err)
		return
	}

	if deployment.Namespace != NamespaceNeedToBeMutated {
		result := model.EmptyAdmissionReviewResponse(admissionReview)
		model.APIResponse(ctx, "OK", http.StatusOK, ctx.Request.Method, result)
		return
	}

	patch, err := d.controller.Mutate(deployment)
	if err != nil {
		model.ErrResponse(ctx, admissionReview, err)
		return
	}

	result, err := model.SuccessAdmissionReviewResponse(admissionReview, patch)
	if err != nil {
		model.ErrResponse(ctx, admissionReview, err)
		return
	}

	model.APIResponse(ctx, "OK", http.StatusOK, ctx.Request.Method, result)
}