package executor

import (
	"fmt"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"

	tasks "github.com/owainlewis/frequency/pkg/tasks"
)

// PodTaskExecutor ...
type PodTaskExecutor struct {
	Client kubernetes.Interface
}

// NewPodTaskExecutor creates a properly configured PodTaskExecutor
func NewPodTaskExecutor(clientset kubernetes.Interface) PodTaskExecutor {
	return PodTaskExecutor{Client: clientset}
}

// Execute will execute a single job
func (e PodTaskExecutor) Execute(task tasks.Task) error {
	glog.Infof("Executing Pod task: %+v", task)

	if task.GetKind() != "PodTask" {
		return fmt.Errorf("Invalid task kind for PodTaskExecutor")
	}

	podTask := task.(tasks.PodTask)
	taskPod := e.newPod(podTask)

	// TODO which namespace to run in (must be configurable)
	_, err := e.Client.CoreV1().Pods(v1.NamespaceDefault).Create(taskPod)
	if err != nil {
		glog.Infof("Failed to create Pod: %s", err)
		return err
	}

	return nil
}

func (e PodTaskExecutor) newPod(task tasks.PodTask) *v1.Pod {
	primary := v1.Container{
		Name:       "primary",
		Image:      task.Image,
		WorkingDir: task.Workspace,
		Env:        task.Env,
		Command:    task.Command.Cmd,
		Args:       task.Command.Args,
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{},
		},
		Spec: v1.PodSpec{
			Containers:    []v1.Container{primary},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	pod.SetGenerateName("task-")

	return pod
}
