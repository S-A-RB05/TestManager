package kubernetes

import (
	"context"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var mt5_image string = "stockbrood/mt5_nogui"

func CreateJob(namespace, jobName, command string) error {
	// Load kubeconfig file and create clientset
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}

	// Define the job
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mt5-container",
							Image: mt5_image,
							Command: []string{
								"/bin/sh",
								"-c",
								command,
							},
						},
					},
					RestartPolicy: "Never",
				},
			},
			BackoffLimit: func() *int32 {
				i := int32(3)
				return &i
			}(),
		},
	}

	// Create the job
	jobClient := clientset.BatchV1().Jobs(namespace)
	result, err := jobClient.Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// Print the job name and UID
	fmt.Printf("Job %s created with UID %s\n", result.Name, result.UID)

	return nil
}
