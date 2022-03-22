package model

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type JSONPatchEntry struct {
	OP    string          `json:"op"`
	Path  string          `json:"path"`
	Value json.RawMessage `json:"value,omitempty"`
}

func SuccessAdmissionReviewResponse(admissionReview *admissionv1.AdmissionReview, patch []JSONPatchEntry) (*admissionv1.AdmissionReview, error) {
	patchBytes, err := json.Marshal(&patch)
	if err != nil {
		return nil, err
	}

	patchType := admissionv1.PatchTypeJSONPatch

	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:       admissionReview.Request.UID,
			Allowed:   true,
			Patch:     patchBytes,
			PatchType: &patchType,
		},
	}, nil
}

func FailedAdmissionReviewResponse(admissionReview *admissionv1.AdmissionReview, err error) *admissionv1.AdmissionReview {
	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     admissionReview.Request.UID,
			Allowed: false,
			Result:  &metav1.Status{Message: err.Error()},
		},
	}
}

func EmptyAdmissionReviewResponse(admissionReview *admissionv1.AdmissionReview) *admissionv1.AdmissionReview {
	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID: admissionReview.Request.UID,
		},
	}
}

type ResponseBody struct {
	StatusCode int         `json:"statusCode"`
	Method     string      `json:"method"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func APIResponse(ctx *gin.Context, Message string, StatusCode int, Method string, Data interface{}) {
	jsonResponse := ResponseBody{
		StatusCode: StatusCode,
		Method:     Method,
		Message:    Message,
		Data:       Data,
	}

	if StatusCode >= 400 {
		ctx.JSON(StatusCode, jsonResponse)
		defer ctx.AbortWithStatus(StatusCode)
	} else {
		ctx.JSON(StatusCode, jsonResponse)
	}
}

func ErrResponse(ctx *gin.Context, admissionReview *admissionv1.AdmissionReview, err error) {
	inject := FailedAdmissionReviewResponse(admissionReview, err)
	APIResponse(ctx, err.Error(), http.StatusInternalServerError, ctx.Request.Method, inject)
}
