import { render, screen } from "@testing-library/react";
import App from "./App";

test("renders header", () => {
  render(<App />);
  const headerElement = screen.getByRole("banner"); // Assumes your <Header /> uses <header> tag
  expect(headerElement).toBeInTheDocument();
});
