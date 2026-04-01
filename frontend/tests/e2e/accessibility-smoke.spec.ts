import { test, expect } from "@playwright/test";

test.describe("Accessibility smoke", () => {
  test("queue page exposes alert or live status area", async ({ page }) => {
    await page.goto("/queue/1");
    await expect(page.getByRole("heading", { name: "대기열" })).toBeVisible();

    const errorBannerCount = await page.getByTestId("error-banner").count();
    if (errorBannerCount > 0) {
      await expect(page.getByTestId("error-banner").first()).toBeVisible();
    } else {
      await expect(page.locator('[aria-live="polite"]').first()).toBeVisible();
    }
  });

  test("booking page seat controls expose buttons", async ({ page }) => {
    await page.goto("/booking/1");
    await expect(page.getByRole("heading", { name: "좌석 선택" })).toBeVisible();

    await expect(page.getByRole("button", { name: /좌석 홀드|좌석 홀드 중/ })).toBeVisible();
    await expect(page.getByRole("button", { name: /예매 확정|예매 확정 중|예매 완료/ })).toBeVisible();
  });
});
