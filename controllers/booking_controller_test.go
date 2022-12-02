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

var _ = Describe("Booking controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		BookingName         = "test-resource"
		BookingNamespace    = "default"
		BookingResourceName = "analytics"
		BookingStart        = "2022-11-21T17:35:05Z"
		BookingEnd          = "2022-11-22T07:40:05Z"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When updating CronJob Status", func() { // TODO
		It("Should increase CronJob Status.Active count when new Jobs are created", func() { // TODO
			By("By creating a new Booking")
			ctx := context.Background()
			resource := &managerv1.Booking{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "manager.kotaico.de/v1",
					Kind:       "Booking",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      BookingName,
					Namespace: BookingNamespace,
				},
				Spec: managerv1.BookingSpec{
					StartAt:      BookingStart,
					EndAt:        BookingEnd,
					ResourceName: BookingResourceName,
				},
			}
			Expect(k8sClient.Create(ctx, resource)).Should(Succeed())

			// Check that the spec we passed are matching
			resourceLookupKey := types.NamespacedName{Name: BookingName, Namespace: BookingNamespace}
			createdBooking := &managerv1.Booking{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(createdBooking.Spec.ResourceName).Should(Equal(BookingResourceName))
		})
	})
})
