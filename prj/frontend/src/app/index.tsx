import * as React from "react";
import * as ReactDOM from "react-dom/client";
import { Provider } from "app/provider";
import "./style/index.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider />
  </React.StrictMode>,
);
