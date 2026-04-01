import { test, expect } from "@playwright/test";

test.describe("Interaction states", () => {
  test("booking page shows disabled action buttons on invalid direct entry", async ({ page }) => {
    await page.goto("/booking/1");

    await expect(page.getByRole("heading", { name: "좌석 선택" })).toBeVisible();
    await expect(page.getByTestId("error-banner")).toBeVisible();

    const holdButton = page.getByRole("button", { name: /좌석 홀드|좌석 홀드 중/ });
    const bookingButton = page.getByRole("button", { name: /예매 확정|예매 확정 중|예매 완료/ });

    await expect(holdButton).toBeDisabled();
    await expect(bookingButton).toBeDisabled();
  });

  test("queue page exposes loading, error, or queue content", async ({ page }) => {
    await page.goto("/queue/1");

    await expect(
      page.getByTestId("loading-state")
        .or(page.getByTestId("error-banner"))
        .or(page.getByText("Queue Token"))
    ).toBeVisible();
  });
});
