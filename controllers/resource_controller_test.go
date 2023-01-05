package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	//+kubebuilder:scaffold:imports
)

var _ = Describe("Resource controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		ResourceName      = "test-resource"
		ResourceTag       = "analytics"
		ResourceType      = "ec2"
		ResourceNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	ctx := context.Background()

	Context("TODO", func() {
		BeforeEach(func() {
			resource := &managerv1.Resource{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "manager.kotaico.de/v1",
					Kind:       "Resource",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      ResourceName,
					Namespace: ResourceNamespace,
				},
				Spec: managerv1.ResourceSpec{
					Booked: false,
					Tag:    ResourceTag,
				},
			}
			Expect(k8sClient.Create(ctx, resource)).Should(Succeed())
		})

		It("TODO", func() {
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
		// TODO: THe case where booked is true
	})
})
