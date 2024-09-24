import { createMemoryHistory, createRouter } from "@tanstack/react-router";
import { StartServer } from "@tanstack/start/server";
import ReactDOMServer from "react-dom/server";
import { routeTree } from "./routes";
globalThis.renderApp = renderApp;

async function renderApp(url: string) {
  // return "Hello";

  const memoryHistory = createMemoryHistory({
    initialEntries: [url],
    initialIndex: 0,
  });

  // Create a request handler

  const router = createRouter({
    routeTree,
    defaultPreload: "viewport",
    isServer: true,
  });
  await router.load();

  router.update({
    history: memoryHistory,
  });

  // return "Hello World";
  // Let's use the default stream handler to create the response
  return ReactDOMServer.renderToString(
    <StartServer router={router}></StartServer>,
  );
}
