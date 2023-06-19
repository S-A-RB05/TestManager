package kubernetes

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Docker hub image to start inside job
var mt5Image string = "stockbrood/mt5_nogui"

// Function to create a job inside the Kubernetes cluster
func CreateJob(namespace string) (jobId string, err error) {
	// Create a Kubernetes client using the local configuration
	fmt.Println("trying to create client from local config")
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("config failed")
		return "", err
	}

	// Create the Kubernetes clientset
	fmt.Println("trying to create clientset")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("clientset failed")
		return "", nil
	}

	// Generate a unique job ID
	jobID := GenerateJobID()

	// Define the environment variable for the container to retrieve
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
							Image: mt5Image,
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

	fmt.Println("trying to create job ")
	// Create the job
	jobClient := clientset.BatchV1().Jobs(namespace)
	result, err := jobClient.Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		fmt.Println("fialed to create job")
		return "", err
	}

	// Print the job name and UID
	fmt.Printf("Job %s created with UID %s\n", result.Name, result.UID)

	return jobID, err
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
