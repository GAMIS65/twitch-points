import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import Dashboard from "./Dashboard.tsx";
import SignInPage from "./SignIn.tsx";
import { BrowserRouter, Route, Routes } from "react-router";
import { SWRProvider } from "./providers/SWRProvider";
import AddRewardPage from "./AddReward.tsx";
import WheelPage from "./Wheel.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <SWRProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/SignIn" element={<SignInPage />} />
          <Route path="/AddReward" element={<AddRewardPage />} />
          <Route path="/Wheel" element={<WheelPage />} />
        </Routes>
      </BrowserRouter>
    </SWRProvider>
  </StrictMode>,
);
