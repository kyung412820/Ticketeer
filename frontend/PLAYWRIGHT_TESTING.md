# Playwright Testing Guide for Accessibility

## Key Tests
1. Focus movement for seat selection
2. Focus on booking complete result
3. Error and loading states visibility with aria-live regions

## Expected Outcomes
- Focus is correctly moved to selected seat and booking completion results
- Proper use of aria-pressed for interactive elements
- Alert regions are announced for screen reader users

### Running Tests
Ensure the application is running before executing the tests:
```bash
npm run dev
npx playwright test
```
