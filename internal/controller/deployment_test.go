package controller

import (
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"reflect"
	"testing"
)

func TestDeploymentController_Mutate(t *testing.T) {
	lessReplica := int32(1)
	moreReplica := int32(5)

	cpuLimit100 := resource.NewQuantity(100, resource.DecimalSI)
	cpuLimit300 := resource.NewQuantity(300, resource.DecimalSI)
	cpuLimit700 := resource.NewQuantity(700, resource.DecimalSI)

	cpuReq100 := resource.NewQuantity(100, resource.DecimalSI)
	cpuReq300 := resource.NewQuantity(300, resource.DecimalSI)
	cpuReq700 := resource.NewQuantity(700, resource.DecimalSI)

	type args struct {
		deployment *appsv1.Deployment
	}
	tests := []struct {
		name    string
		args    args
		want    []*model.JSONPatchEntry
		wantErr bool
	}{
		{
			"mutated replicas and resources",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: &moreReplica,
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit100,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq100,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			[]*model.JSONPatchEntry{
				&model.JSONPatchEntry{
					OP:    "replace",
					Path:  "/spec/replicas",
					Value: []byte("3"),
				},
				&model.JSONPatchEntry{
					OP:    "replace",
					Path:  "/spec/template/spec/containers",
					Value: []byte("[{\"name\":\"\",\"resources\":{\"limits\":{\"cpu\":\"200\"},\"requests\":{\"cpu\":\"200\"}}}]"),
				},
			},
			false,
		},
		{
			"don't need to mutate",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: &lessReplica,
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit300,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq300,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			[]*model.JSONPatchEntry{
				&model.JSONPatchEntry{
					OP:    "replace",
					Path:  "/spec/replicas",
					Value: []byte("1"),
				},
				&model.JSONPatchEntry{
					OP:    "replace",
					Path:  "/spec/template/spec/containers",
					Value: []byte("[{\"name\":\"\",\"resources\":{\"limits\":{\"cpu\":\"300\"},\"requests\":{\"cpu\":\"300\"}}}]"),
				},
			},
			false,
		},
		{
			"inject resources",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: &lessReplica,
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit700,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq700,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DeploymentController{}
			got, err := d.Mutate(tt.args.deployment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mutate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mutate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkReplicas(t *testing.T) {
	lessReplica := int32(1)
	moreReplica := int32(5)

	type args struct {
		deployment *appsv1.Deployment
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"less than replicas 3",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: &lessReplica,
					},
				},
			},
			[]byte("1"),
			false,
		},
		{
			"more than replicas 3",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Replicas: &moreReplica,
					},
				},
			},
			[]byte("3"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkReplicas(tt.args.deployment)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkReplicas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkReplicas() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_checkResource(t *testing.T) {
	cpuLimit100 := resource.NewQuantity(100, resource.DecimalSI)
	cpuLimit300 := resource.NewQuantity(300, resource.DecimalSI)
	cpuLimit700 := resource.NewQuantity(700, resource.DecimalSI)

	cpuReq100 := resource.NewQuantity(100, resource.DecimalSI)
	cpuReq300 := resource.NewQuantity(300, resource.DecimalSI)
	cpuReq700 := resource.NewQuantity(700, resource.DecimalSI)

	type args struct {
		deployment *appsv1.Deployment
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"less than min all",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit100,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq100,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			[]byte("[{\"name\":\"\",\"resources\":{\"limits\":{\"cpu\":\"200\"},\"requests\":{\"cpu\":\"200\"}}}]"),
			false,
		},
		{
			"less than min req",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit300,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq100,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			[]byte("[{\"name\":\"\",\"resources\":{\"limits\":{\"cpu\":\"300\"},\"requests\":{\"cpu\":\"200\"}}}]"),
			false,
		},
		{
			"less than min limit",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit100,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq300,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			[]byte("[{\"name\":\"\",\"resources\":{\"limits\":{\"cpu\":\"200\"},\"requests\":{\"cpu\":\"300\"}}}]"),
			false,
		},
		{
			"fit value",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit300,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq300,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			[]byte("[{\"name\":\"\",\"resources\":{\"limits\":{\"cpu\":\"300\"},\"requests\":{\"cpu\":\"300\"}}}]"),
			false,
		},
		{
			"more than max all",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit700,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq700,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nil,
			true,
		},
		{
			"more than max limit",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit700,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq300,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nil,
			true,
		},
		{
			"more than max req",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit300,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq700,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nil,
			true,
		},
		{
			"more than max total",
			args{
				deployment: &appsv1.Deployment{
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit100,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq100,
											},
										},
									},
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit300,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq300,
											},
										},
									},
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit300,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq300,
											},
										},
									},
									corev1.Container{
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												corev1.ResourceCPU: *cpuLimit300,
											},
											Requests: corev1.ResourceList{
												corev1.ResourceCPU: *cpuReq300,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkResource(tt.args.deployment)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkResource() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_checkCPU(t *testing.T) {
	type args struct {
		cpu int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"less than min",
			args{
				cpu: 100,
			},
			int64(200),
			false,
		},
		{
			"fit value",
			args{
				cpu: 300,
			},
			int64(300),
			false,
		},
		{
			"more than max",
			args{
				cpu: 700,
			},
			int64(0),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkEachCPU(tt.args.cpu)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkEachCPU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkEachCPU() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkTotalCPU(t *testing.T) {
	type args struct {
		cpu int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"less than total",
			args{
				cpu: 700,
			},
			int64(700),
			false,
		},
		{
			"more than total",
			args{
				cpu: 1500,
			},
			int64(0),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkTotalCPU(tt.args.cpu)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkTotalCPU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkTotalCPU() got = %v, want %v", got, tt.want)
			}
		})
	}
}
