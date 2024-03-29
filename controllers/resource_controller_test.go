package controllers

import (
	"context"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	//+kubebuilder:scaffold:imports
)

var _ = Describe("Resource controller", func() {
	ctx := context.Background()

	const (
		ResourceName = "test-resource"
		ResourceTag  = "analytics"
		ResourceType = "ec2"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	ResourceNamespace := os.Getenv("NAMESPACE")
	if ResourceNamespace == "" {
		ResourceNamespace = "default"
	}

	Context("Resource basics", func() {
		BeforeEach(func() {
			resource := &managerv1.Resource{
				ObjectMeta: metav1.ObjectMeta{
					Name:      ResourceName,
					Namespace: ResourceNamespace,
				},
				Spec: managerv1.ResourceSpec{
					BookedBy: "test",
					Tag:      ResourceTag,
					Type:     ResourceType,
				},
			}
			Expect(k8sClient.Create(ctx, resource)).Should(Succeed())
		})

		It("Asserts resource creation", func() {
			By("By creating a new Resource")
			// Check that the spec we passed is matching
			resourceLookupKey := types.NamespacedName{Name: ResourceName, Namespace: ResourceNamespace}
			createdResource := &managerv1.Resource{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, resourceLookupKey, createdResource)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdResource.Spec.Tag).Should(Equal(ResourceTag))

			Expect(createdResource.Status).Should(Equal(managerv1.ResourceStatus{
				Instances: 0,
				Running:   0,
				Status:    "",
			}))
		})
		// TODO: The case where booked is true
	})
})
