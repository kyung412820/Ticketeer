import { test, expect } from "@playwright/test";

test.describe("Ticketeer booking flow", () => {
  test("events page renders and allows navigation", async ({ page }) => {
    await page.goto("/events");

    await expect(page.getByRole("heading", { name: "공연 목록" })).toBeVisible();
    await expect(page.getByText("TICKETEER")).toBeVisible();

    const detailLink = page.getByRole("link", { name: "상세 보기" }).first();
    await expect(detailLink).toBeVisible();
    await detailLink.click();

    await expect(page).toHaveURL(/\/events\/\d+/);
  });

  test("queue page renders status area", async ({ page }) => {
    await page.goto("/queue/1");

    await expect(page.getByRole("heading", { name: "대기열" })).toBeVisible();
    await expect(page.getByText(/대기 상태를 주기적으로 갱신합니다/)).toBeVisible();

    await expect(
      page.getByTestId("loading-state")
        .or(page.getByTestId("error-banner"))
        .or(page.getByText("Queue Token"))
    ).toBeVisible();
  });

  test("booking page blocks invalid direct entry", async ({ page }) => {
    await page.goto("/booking/1");

    await expect(page.getByRole("heading", { name: "좌석 선택" })).toBeVisible();
    await expect(page.getByTestId("error-banner")).toContainText("유효한 접근이 아닙니다");
  });
});
