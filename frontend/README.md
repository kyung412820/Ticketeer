# Ticketeer - Ticket Booking System

## Overview
Ticketeer is a ticket booking platform that handles event ticketing, seat selection, queue management, and booking confirmation. The project aims to demonstrate the usage of frontend technologies like React, Next.js, and React Query, alongside the implementation of real-world performance and accessibility best practices.

## Frontend Design Decisions

### 1. Why /events is server-side rendered and /queue, /booking is client-side rendered?
- **/events page** is server-side rendered (SSR) because it is a public-facing page that needs to be indexed by search engines. Additionally, it provides a list of events that do not depend on user-specific data and can be fetched efficiently on the server.
- **/queue and /booking pages** are client-side rendered (CSR) to prioritize fast, dynamic user interactions. These pages require real-time data like the queue status and seat availability, which is better suited for client-side updates.

### 2. Why we minimized optimistic updates?
We minimized optimistic updates because the application’s core functionality revolves around ensuring data consistency, especially when it comes to booking and seat holding. Optimistic updates in these areas could lead to discrepancies between the client and the server, which could harm the user experience and cause errors such as booking confirmation mismatches.

### 3. Why we store the queue token in sessionStorage?
Queue tokens are stored in `sessionStorage` because they are session-specific and are not meant to persist beyond the user’s session. This ensures that the queue token is cleared when the session ends, maintaining the security of the token and ensuring it is not accessible after the user has logged out or closed the browser.

### 4. Why we did not implement retry for booking?
Retrying booking requests could lead to duplicate bookings and other issues related to transaction integrity. Instead, we ensure that the user is informed when the booking cannot be processed, and they are given the option to try again with a clear indication of what went wrong.

### 5. Why we implemented React Query for queue and seats management?
We implemented React Query to manage the server state for the queue and seat information to improve the efficiency of data fetching, caching, and synchronization between the server and the client. React Query also simplifies the logic for managing the loading and error states associated with data fetching, and it allows us to invalidate data when it is no longer valid.

### 6. Why we integrated Playwright for accessibility testing?
Playwright was integrated to ensure that our application is fully accessible to screen reader users and other assistive technology users. By automating accessibility tests, we can ensure that dynamic changes in the UI (such as the selection of seats or booking completion) are properly announced and do not cause confusion.

## Troubleshooting (Frontend-Specific)

### 1. Issue: Stale UI states due to real-time data updates
- **Cause**: When polling data from the server (like queue status and seat availability), the UI may display outdated information due to stale client-side state.
- **Solution**: We implemented React Query’s cache invalidation and refetching mechanisms to ensure that the UI remains in sync with the server at all times. In addition, the UI is re-rendered only when necessary to avoid unnecessary performance hits.

### 2. Issue: Multiple form submissions due to repeated clicks
- **Cause**: Users could submit forms (like seat booking) multiple times, leading to unintended consequences.
- **Solution**: We implemented button disabling strategies, preventing users from submitting the form multiple times. This is accompanied by feedback indicating that the form is being processed.

### 3. Issue: Invalid booking attempts
- **Cause**: Some users attempted to directly access the booking page without going through the queue process.
- **Solution**: We implemented guards and error messages that prevent direct access to the booking page unless the user has completed the queue process.

## Future Improvements
- Add deeper integration of React Query across all data-fetching operations.
- Introduce server-side error boundaries and retry strategies for network requests.
- Expand on accessibility testing with more detailed screen reader checks.

