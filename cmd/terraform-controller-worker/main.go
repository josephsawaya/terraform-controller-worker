package main

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"

	"github.com/josephsawaya/terraform-controller-worker/cmd/terraform-controller-worker/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

type TerraformFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

var RootCmd = &cobra.Command{
	Use:   "work",
	Short: "Command to run the terraform-controller's worker",
	Run: func(cmd *cobra.Command, args []string) {
		rconf, err := rest.InClusterConfig()
		if err != nil {
			klog.Error(err)
			return
		}

		// TODO: Using dynamic client here feels incorrect
		dyn, err := dynamic.NewForConfig(rconf)
		if err != nil {
			klog.Error(err)
			return
		}

		res := schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "terraforms"}

		terraformList, err := dyn.Resource(res).Namespace("argocd").List(context.Background(), metav1.ListOptions{})
		if err != nil {
			klog.Error(err)
			return
		}

		for _, unstructuredTerraform := range terraformList.Items {
			listOfContents, exists, err := unstructured.NestedSlice(unstructuredTerraform.UnstructuredContent(), "spec", "list")
			if err != nil {
				klog.Error(err)
				return
			}
			if exists == false {
				klog.Errorf("List field does not exist for %v", unstructuredTerraform.GetName())
				return
			}

			for _, content := range listOfContents {
				klog.Infof("%v: %+v\n\n", unstructuredTerraform.GetName(), content)
				terraformFile, ok := content.(map[string]interface{})
				if !ok {
					klog.Errorf("Unable to convert %+v to TerraformFile", content)
					return
				}

				data, err := base64.StdEncoding.DecodeString(terraformFile["content"].(string))
				if err != nil {
					klog.Error(err)
					return
				}

				terraformPath := strings.TrimPrefix(terraformFile["name"].(string), "./")

				terraformDir := filepath.Dir(terraformPath)

				terraformFileName := filepath.Base(terraformPath)

				err = os.MkdirAll(terraformDir, os.ModePerm)
				if err != nil {
					klog.Error(err)
					return
				}

				err = os.WriteFile(terraformFileName, data, os.ModePerm)
				if err != nil {
					klog.Error(err)
					return
				}
			}
		}

		// TODO: Replace these with tfexec package
		err = util.RunTerraformCommand("init", nil)
		if err != nil {
			klog.Error(err)
			return
		}

		klog.Info("Init")
		err = util.RunTerraformCommand("plan", nil)
		if err != nil {
			klog.Error(err)
			return
		}

		klog.Info("Planned")

		yes := "yes\n"
		err = util.RunTerraformCommand("apply", &yes)
		if err != nil {
			klog.Error(err)
			return
		}

		klog.Info("Applied")

		// 	time.Sleep(time.Duration(3) * time.Minute)
	},
}

func main() {
	RootCmd.Execute()
}
