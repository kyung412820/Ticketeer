import { test, expect } from "@playwright/test";

test.describe("Enhanced Accessibility Tests", () => {
  test("Focus moves to selected seat after click", async ({ page }) => {
    await page.goto("/booking/1");
    const seat = page.getByRole("button", { name: "A-1" });
    await seat.click();

    const selectedSeat = page.getByRole("button", { name: "A-1 선택됨" });
    await expect(selectedSeat).toBeFocused();
  });

  test("Booking complete moves focus to the result card", async ({ page }) => {
    await page.goto("/booking/1");
    await page.getByRole("button", { name: "예매 확정" }).click();

    const bookingComplete = page.locator("#booking-heading");
    await expect(bookingComplete).toBeFocused();
  });

  test("aria-live regions for error and loading states", async ({ page }) => {
    await page.goto("/queue/1");
    await expect(page.getByRole("alert")).toBeVisible();
  });
});
