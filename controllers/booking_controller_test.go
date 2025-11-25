package controllers

import (
	"context"
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	managerv1 "github.com/kotaicode/resource-booking-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	//+kubebuilder:scaffold:imports
)

func getBookingNamespace() string {
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		return "default"
	}
	return namespace
}

var _ = Describe("Booking controller", func() {
	ctx := context.Background()

	const (
		BookingName         = "test-booking"
		BookingResourceName = "ec2.analytics"
	)

	// Keep the format visible for easier debugging and just increment with an year
	var (
		ScheduledBookingStart = fmt.Sprintf("%d-01-01T00:00:00Z", time.Now().AddDate(1, 0, 0).Year())
		ScheduledBookingEnd   = fmt.Sprintf("%d-01-02T00:00:00Z", time.Now().AddDate(1, 0, 0).Year())

		InProgressBookingStart = fmt.Sprintf("%d-01-01T00:00:00Z", time.Now().Year())
		InProgressBookingEnd   = fmt.Sprintf("%d-01-01T00:00:00Z", time.Now().AddDate(1, 0, 0).Year())

		FinishedBookingStart = fmt.Sprintf("%d-01-01T00:00:00Z", time.Now().AddDate(-1, 0, 0).Year())
		FinishedBookingEnd   = fmt.Sprintf("%d-01-02T00:00:00Z", time.Now().AddDate(-1, 0, 0).Year())

		BookingNamespace = getBookingNamespace()
	)

	var bookingSpec managerv1.BookingSpec
	var booking *managerv1.Booking
	var resource *managerv1.Resource

	Context("Booking changes", func() {
		BeforeEach(func() {
			bookingSpec = managerv1.BookingSpec{
				StartAt:      ScheduledBookingStart,
				EndAt:        ScheduledBookingEnd,
				ResourceName: BookingResourceName,
			}

			booking = &managerv1.Booking{
				ObjectMeta: metav1.ObjectMeta{
					Name:      BookingName,
					Namespace: BookingNamespace,
				},
				Spec: bookingSpec,
			}
		})

		AfterEach(func() {
			// Clean up booking and resource
			Expect(k8sClient.Delete(ctx, booking)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, resource)).Should(Succeed())
		})

		It("Should update booking status to SCHEDULED for future bookings", func() {
			By("By creating a new Booking with future dates")
			booking.Spec.StartAt = ScheduledBookingStart
			booking.Spec.EndAt = ScheduledBookingEnd
			Expect(k8sClient.Create(ctx, booking)).Should(Succeed())

			// Check that the spec we passed is matching
			resourceLookupKey := types.NamespacedName{Name: BookingName, Namespace: BookingNamespace}
			createdBooking := &managerv1.Booking{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				return err == nil
			}).Should(BeTrue())

			Expect(createdBooking.Spec.StartAt).Should(Equal(ScheduledBookingStart))
			Expect(createdBooking.Spec.EndAt).Should(Equal(ScheduledBookingEnd))

			By("By checking if the booking status is updated to SCHEDULED")
			Eventually(func() (string, error) {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				if err != nil {
					return "", err
				}
				return createdBooking.Status.Status, nil
			}).Should(Equal(managerv1.BookingScheduled), "should show that the booking status is scheduled")
		})

		It("Should update booking status to IN PROGRESS for current bookings", func() {
			By("By creating a new Booking with current dates")

			booking.Spec.StartAt = InProgressBookingStart
			booking.Spec.EndAt = InProgressBookingEnd
			Expect(k8sClient.Create(ctx, booking)).Should(Succeed())

			// Check that the spec we passed is matching
			resourceLookupKey := types.NamespacedName{Name: BookingName, Namespace: BookingNamespace}
			createdBooking := &managerv1.Booking{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				return err == nil
			}).Should(BeTrue())

			Expect(createdBooking.Spec.StartAt).Should(Equal(InProgressBookingStart))
			Expect(createdBooking.Spec.EndAt).Should(Equal(InProgressBookingEnd))

			By("By checking if the booking status is updated to IN PROGRESS")
			Eventually(func() (string, error) {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				if err != nil {
					return "", err
				}
				return createdBooking.Status.Status, nil
			}).Should(Equal(managerv1.BookingInProgress), "should show that the booking status is in progress")

			// TODO Check if the resource spec.booked got updated
		})

		It("Should update booking status to FINISHED for past bookings", func() {
			By("By creating a new Booking with past dates")

			booking.Spec.StartAt = FinishedBookingStart
			booking.Spec.EndAt = FinishedBookingEnd
			Expect(k8sClient.Create(ctx, booking)).Should(Succeed())

			// Check that the spec we passed is matching
			resourceLookupKey := types.NamespacedName{Name: BookingName, Namespace: BookingNamespace}
			createdBooking := &managerv1.Booking{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				return err == nil
			}).Should(BeTrue())

			Expect(createdBooking.Spec.StartAt).Should(Equal(FinishedBookingStart))
			Expect(createdBooking.Spec.EndAt).Should(Equal(FinishedBookingEnd))

			By("By checking if the booking status is updated to FINISHED")
			Eventually(func() (string, error) {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				if err != nil {
					return "", err
				}
				return createdBooking.Status.Status, nil
			}).Should(Equal(managerv1.BookingFinished), "should show that the booking status is finished")

			// TODO Check if the resource spec.booked got updated
		})
	})
})
