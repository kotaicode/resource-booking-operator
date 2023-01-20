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

var _ = Describe("Resource controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		ResourceName      = "test-resource-monitor"
		ResourceTag       = "analytics"
		ResourceType      = "ec2"
		ResourceNamespace = "default"

		timeout  = time.Second * 12
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	ctx := context.Background()

	Context("TODO", func() {
		BeforeEach(func() {
			resourceMonitor := &managerv1.ResourceMonitor{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "manager.kotaico.de/v1",
					Kind:       "ResourceMonitor",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ResourceName,
					Namespace: ResourceNamespace,
				},
			}
			Expect(k8sClient.Create(ctx, resourceMonitor)).Should(Succeed())
		})

		It("TODO", func() {
			By("By get Resource monitor")
			resourceList := &managerv1.ResourceList{}
			Eventually(func() bool {
				err := k8sClient.List(ctx, resourceList)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(len(resourceList.Items)).Should(BeNumerically(">", 0)) // TODO: fix errors
		})
	})
})
