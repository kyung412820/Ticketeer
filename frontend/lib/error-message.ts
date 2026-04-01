const ERROR_MESSAGE_MAP: Record<string, string> = {
  QUEUE_EXPIRED: "입장 시간이 만료되었습니다. 다시 대기열에 진입해주세요.",
  SEAT_ALREADY_HELD: "다른 사용자가 먼저 선택한 좌석입니다.",
  SEAT_ALREADY_BOOKED: "이미 예매가 완료된 좌석입니다.",
  HOLD_NOT_FOUND: "좌석 홀드 시간이 만료되었습니다. 다시 선택해주세요.",
  HOLD_EXPIRED: "좌석 홀드 시간이 만료되었습니다. 다시 선택해주세요.",
  INVALID_QUEUE_TOKEN: "유효하지 않은 접근입니다. 다시 대기열부터 진행해주세요.",
  QUEUE_TOKEN_NOT_FOUND: "대기열 정보가 없습니다. 다시 대기열에 진입해주세요.",
  BOOKING_NOT_OPEN: "아직 예매가 시작되지 않았습니다.",
  EVENT_CLOSED: "예매가 종료된 공연입니다.",
  ALREADY_IN_QUEUE: "이미 대기열에 진입한 사용자입니다.",
  EVENT_NOT_FOUND: "존재하지 않는 공연입니다.",
  SEAT_NOT_FOUND: "존재하지 않는 좌석입니다.",
  SEAT_EVENT_MISMATCH: "공연과 좌석 정보가 일치하지 않습니다.",
  INVALID_REQUEST: "요청 정보를 다시 확인해주세요.",
  INTERNAL_SERVER_ERROR: "잠시 후 다시 시도해주세요.",
};

export function getUserFriendlyErrorMessage(code?: string, fallback?: string): string {
  if (code && ERROR_MESSAGE_MAP[code]) {
    return ERROR_MESSAGE_MAP[code];
  }

  if (fallback && fallback.trim()) {
    return fallback;
  }

  return "요청 처리 중 오류가 발생했습니다.";
}
