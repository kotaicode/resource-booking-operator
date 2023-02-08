package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//+kubebuilder:scaffold:imports
)

var _ = Describe("Resource Monitor controller", func() {

	const (
		ResourceName      = "test-resource-monitor"
		ResourceType      = "ec2"
		ResourceNamespace = "default"

		timeout  = time.Second * 3
		interval = time.Millisecond * 250
	)

	ctx := context.Background()

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
