import React, { useEffect } from "react";

type BookingCompleteProps = {
  bookingCode: string;
  onClose: () => void;
};

export const BookingComplete = ({ bookingCode, onClose }: BookingCompleteProps) => {
  useEffect(() => {
    // Focus on the booking completion heading for accessibility
    const bookingHeading = document.getElementById("booking-heading");
    if (bookingHeading) {
      bookingHeading.focus();
    }
  }, []);

  return (
    <div role="dialog" aria-labelledby="booking-heading">
      <h2 id="booking-heading" tabIndex={-1}>
        예매 완료
      </h2>
      <p>예약 코드: {bookingCode}</p>
      <button onClick={onClose}>닫기</button>
    </div>
  );
};
