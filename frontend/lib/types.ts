export type EventStatus = "OPEN_PENDING" | "OPEN" | "CLOSED";
export type SeatStatus = "AVAILABLE" | "HELD" | "BOOKED";
export type QueueStatus = "WAITING" | "READY" | "EXPIRED";

export type Event = {
  id: number;
  title: string;
  description: string;
  venue: string;
  event_at: string;
  booking_open_at: string;
  booking_close_at: string;
  status: EventStatus;
};

export type Seat = {
  id: number;
  seat_no: string;
  section: string;
  price: number;
  status: SeatStatus;
};

export type QueueEnterResponse = {
  queue_token: string;
  status: QueueStatus;
  position: number;
};

export type QueueStatusResponse = {
  queue_token: string;
  status: QueueStatus;
  position?: number;
  ready_at?: string;
  expired_at?: string;
};

export type HoldSeatResponse = {
  seat_id: number;
  status: "HELD";
  hold_expires_at: string;
};

export type BookingResponse = {
  booking_id: number;
  booking_code: string;
  event_id: number;
  seat_id: number;
  status: "CONFIRMED";
  booked_at: string;
};

export type SeatsResponse = {
  event_id: number;
  seats: Seat[];
};
