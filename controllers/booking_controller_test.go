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
		timeout             = time.Second * 50
		duration            = time.Second * 5
		interval            = time.Millisecond * 250
	)

	// Keep the format visible for easier debugging and just increment with an year
	var (
		ScheduledBookingStart = fmt.Sprintf("%d-01-01T00:00:00Z", time.Now().AddDate(1, 0, 0).Year())
		ScheduledBookingEnd   = fmt.Sprintf("%d-01-02T00:00:00Z", time.Now().AddDate(1, 0, 0).Year())

		InProgressBookingStart = fmt.Sprintf("%d-01-01T00:00:00Z", time.Now().Year())
		InProgressBookingEnd   = fmt.Sprintf("%d-01-01T00:00:00Z", time.Now().AddDate(1, 0, 0).Year())

		FinishedBookingStart = fmt.Sprintf("%d-01-01T00:00:00Z", time.Now().AddDate(-1, 0, 0).Year())
		FinishedBookingEnd   = fmt.Sprintf("%d-01-02T00:00:00Z", time.Now().AddDate(-1, 0, 0).Year())
	)

	ctx := context.Background()

	var bookingSpec managerv1.BookingSpec
	var booking *managerv1.Booking

	Context("TODO", func() {
		BeforeEach(func() {
			bookingSpec = managerv1.BookingSpec{
				StartAt:      ScheduledBookingStart,
				EndAt:        ScheduledBookingEnd,
				ResourceName: BookingResourceName,
			}

			booking = &managerv1.Booking{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "manager.kotaico.de/v1",
					Kind:       "Booking",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      BookingName,
					Namespace: BookingNamespace,
				},
				Spec: bookingSpec,
			}
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(ctx, booking))
		})

		It("Should update booking status", func() {
			By("By creating a new Booking")
			booking.Spec.StartAt = ScheduledBookingStart
			booking.Spec.EndAt = ScheduledBookingEnd
			Expect(k8sClient.Create(ctx, booking)).Should(Succeed())

			// Check that the spec we passed is matching
			resourceLookupKey := types.NamespacedName{Name: BookingName, Namespace: BookingNamespace}
			createdBooking := &managerv1.Booking{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				return err != nil
			}).Should(BeTrue())

			Expect(createdBooking.Spec.StartAt).Should(Equal(ScheduledBookingStart))
			Expect(createdBooking.Spec.EndAt).Should(Equal(ScheduledBookingEnd))

			By("By checking if the booking status is updated")
			Eventually(func() (string, error) {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				if err != nil {
					return "", err
				}
				return createdBooking.Status.Status, nil
			}).Should(Equal(managerv1.BookingScheduled), "should show that the booking status is scheduled")
		})

		It("Should update booking status", func() {
			By("By creating a new Booking")

			booking.Spec.StartAt = InProgressBookingStart
			booking.Spec.EndAt = InProgressBookingEnd
			Expect(k8sClient.Create(ctx, booking)).Should(Succeed())

			err := k8sClient.Update(ctx, booking)
			if err != nil {
				return
			}

			// Check that the spec we passed is matching
			resourceLookupKey := types.NamespacedName{Name: BookingName, Namespace: BookingNamespace}
			createdBooking := &managerv1.Booking{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				return err != nil
			}).Should(BeTrue())

			Expect(createdBooking.Spec.StartAt).Should(Equal(InProgressBookingStart))
			Expect(createdBooking.Spec.EndAt).Should(Equal(InProgressBookingEnd))

			By("By checking if the booking status is updated")
			Eventually(func() (string, error) {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				if err != nil {
					return "", err
				}
				return createdBooking.Status.Status, nil
			}).Should(Equal(managerv1.BookingInProgress), "should show that the booking status is in progres")

			// TODO Check if the resource spec.booked got updated
		})

		It("Should update booking status", func() {
			By("By creating a new Booking")

			booking.Spec.StartAt = FinishedBookingStart
			booking.Spec.EndAt = FinishedBookingEnd
			Expect(k8sClient.Create(ctx, booking)).Should(Succeed())

			err := k8sClient.Update(ctx, booking)
			if err != nil {
				return
			}

			// Check that the spec we passed is matching
			resourceLookupKey := types.NamespacedName{Name: BookingName, Namespace: BookingNamespace}
			createdBooking := &managerv1.Booking{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, resourceLookupKey, createdBooking)
				return err != nil
			}).Should(BeTrue())

			Expect(createdBooking.Spec.StartAt).Should(Equal(FinishedBookingStart))
			Expect(createdBooking.Spec.EndAt).Should(Equal(FinishedBookingEnd))

			By("By checking if the booking status is updated")
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
