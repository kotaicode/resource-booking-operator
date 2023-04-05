package controllers

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//+kubebuilder:scaffold:imports
)

var _ = Describe("Resource Monitor controller", func() {
	ctx := context.Background()

	const (
		ResourceName = "test-resource-monitor"
		ResourceType = "ec2"

		timeout  = time.Second * 3
		interval = time.Millisecond * 250
	)

	ResourceNamespace := os.Getenv("NAMESPACE")
	if ResourceNamespace == "" {
		ResourceNamespace = "default"
	}

	Context("Resource monitor management", func() {
		BeforeEach(func() {
			resourceMonitor := managerv1.ResourceMonitor{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ResourceName,
					Namespace: ResourceNamespace,
				},
				Spec: managerv1.ResourceMonitorSpec{
					Type: ResourceType,
				},
			}
			Expect(k8sClient.Create(ctx, &resourceMonitor)).Should(Succeed())
		})

		It("Asserts that the resource monitor created resources", func() {
			By("By getting the list of resources")
			resourceList := managerv1.ResourceList{}
			Eventually(func() bool {
				err := k8sClient.List(ctx, &resourceList)
				return err == nil && len(resourceList.Items) > 1
			}, timeout, interval).Should(BeTrue())
		})
	})
})
