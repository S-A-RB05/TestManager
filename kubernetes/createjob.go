package kubernetes

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var mt5_image string = "stockbrood/mt5_nogui"

func CreateJob(namespace string) error {
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

	// Generate a unique job ID
	jobID := GenerateJobID()

	// Define the environment variable
	jobIDEnvVar := corev1.EnvVar{
		Name:  "JOB_ID",
		Value: jobID,
	}

	// Define the job
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "job-", // Use a generate name instead of a fixed job name
			Namespace:    namespace,
			Labels: map[string]string{
				"job-id": jobID,
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mt5-container",
							Image: mt5_image,
							Env: []corev1.EnvVar{
								jobIDEnvVar,
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

// GenerateJobID generates a unique job ID based on timestamp and a unique identifier
func GenerateJobID() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	uniqueID := generateUniqueID(8)
	return fmt.Sprintf("job-%d-%s", timestamp, uniqueID)
}

// generateUniqueID generates a random alphanumeric string of the specified length
func generateUniqueID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
