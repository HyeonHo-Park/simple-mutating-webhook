package model

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type JSONPatchEntry struct {
	OP    string          `json:"op"`
	Path  string          `json:"path"`
	Value json.RawMessage `json:"value,omitempty"`
}

func getSuccessResponse(uid types.UID, patch []*JSONPatchEntry) *admissionv1.AdmissionResponse {
	patchBytes, err := json.Marshal(&patch)
	if err != nil {
		log.Errorf("Failed Marshal: %v", err)
		return nil
	}

	patchType := admissionv1.PatchTypeJSONPatch

	return &admissionv1.AdmissionResponse{
		UID:       uid,
		Allowed:   true,
		Patch:     patchBytes,
		PatchType: &patchType,
	}
}

func getFailedResponse(uid types.UID, err error) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{
		UID:     uid,
		Allowed: false,
		Result:  &metav1.Status{Message: err.Error()},
	}
}

func getEmptyResponse(uid types.UID) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{ UID: uid }
}

func makeResponseObj(uid types.UID, data interface{}) *admissionv1.AdmissionReview {
	var response *admissionv1.AdmissionResponse

	switch data.(type) {
	case []*JSONPatchEntry:
		response = getSuccessResponse(uid, data.([]*JSONPatchEntry))
	case error:
		response = getFailedResponse(uid, data.(error))
	default:
		response = getEmptyResponse(uid)
	}

	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: response,
	}
}

func WriteResponse(ctx *gin.Context, ar *admissionv1.AdmissionReview, data interface{}) {
	obj := makeResponseObj(ar.Request.UID, data)

	resp, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("Failed Marshal: %v", err)
	}

	ctx.Writer.Write(resp)
}
