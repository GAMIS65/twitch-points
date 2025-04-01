import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import Dashboard from "./Dashboard.tsx";
import SignInPage from "./SignIn.tsx";
import { BrowserRouter, Route, Routes } from "react-router";
import { SWRProvider } from "./providers/SWRProvider";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <SWRProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/SignIn" element={<SignInPage />} />
        </Routes>
      </BrowserRouter>
    </SWRProvider>
  </StrictMode>,
);
