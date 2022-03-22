package model

import (
	"encoding/json"
	"fmt"
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

func makeResponseObj(uid types.UID, data interface{}) *admissionv1.AdmissionReview {
	var response *admissionv1.AdmissionResponse

	switch data.(type) {
	case []*JSONPatchEntry:
		patchBytes, err := json.Marshal(&data)
		if err != nil {
		}

		patchType := admissionv1.PatchTypeJSONPatch

		response = &admissionv1.AdmissionResponse{
			UID:       uid,
			Allowed:   true,
			Patch:     patchBytes,
			PatchType: &patchType,
		}
	case error:
		response = &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: false,
			Result:  &metav1.Status{Message: fmt.Sprintf("%v", data)},
		}
	default:
		response = &admissionv1.AdmissionResponse{
			UID: uid,
		}
	}

	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: response,
	}
}

func APIResponse(ctx *gin.Context, ar *admissionv1.AdmissionReview, data interface{}) {
	obj := makeResponseObj(ar.Request.UID, data)
	resp, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("Failed Marshal: %v", err)
	}

	ctx.Writer.Write(resp)
}
