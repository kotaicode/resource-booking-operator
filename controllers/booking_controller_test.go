package controllers

import (
	"context"
	"fmt"
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
		timeout             = time.Second * 10
		duration            = time.Second * 10
		interval            = time.Millisecond * 250
	)

	// Keep the format visible for easier debugging and just increment with an year
	var (
		BookingStart = fmt.Sprintf("%d-01-01T00:00:00Z", time.Now().AddDate(1, 0, 0).Year())
		BookingEnd   = fmt.Sprintf("%d-01-02T00:00:00Z", time.Now().AddDate(1, 0, 0).Year())
	)

	Context("When creating a booking", func() {
		It("Should update booking status", func() {
			By("By creating a new Booking")
			ctx := context.Background()

			BookingSpec := managerv1.BookingSpec{
				StartAt:      BookingStart,
				EndAt:        BookingEnd,
				ResourceName: BookingResourceName,
			}

			booking := &managerv1.Booking{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "manager.kotaico.de/v1",
					Kind:       "Booking",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      BookingName,
					Namespace: BookingNamespace,
				},
				Spec: BookingSpec,
			}
			Expect(k8sClient.Create(ctx, booking)).Should(Succeed())

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

			Expect(createdBooking.Spec).Should(Equal(BookingSpec))

			By("By checking if the booking status is update")
			Consistently(func() (string, error) {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				if err != nil {
					return "", err
				}
				return createdBooking.Status.Status, nil
			}, duration, interval).Should(Equal(managerv1.BookingScheduled), "should show that the booking status is scheduled")
		})
	})
})
