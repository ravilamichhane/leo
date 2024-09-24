import { createRouter, RouterProvider } from "@tanstack/react-router";
import ReactDOM from "react-dom/client";
import { routeTree } from "./routes";
const router = createRouter({
  routeTree,
  defaultPreload: "intent",
  defaultStaleTime: 5000,
});
const root = ReactDOM.hydrateRoot(
  document.getElementById("root"),
  <RouterProvider router={router} />,
);
