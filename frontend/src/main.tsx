import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "@/index.css";
import Dashboard from "@/pages/Dashboard.tsx";
import SignInPage from "@/pages//SignIn.tsx";
import { BrowserRouter, Route, Routes } from "react-router";
import { SWRProvider } from "@/providers/SWRProvider";
import AddRewardPage from "@/pages/AddReward.tsx";
import { Layout } from "@/components/dashboard/layout.tsx";
import WheelPage from "@/pages/Wheel.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <SWRProvider>
      <BrowserRouter>
        <Routes>
          <Route element={<Layout />}>
            <Route path="/" element={<Dashboard />} />
            <Route path="/sign-in" element={<SignInPage />} />
            <Route path="/AddReward" element={<AddRewardPage />} />
            <Route path="/wheel" element={<WheelPage />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </SWRProvider>
  </StrictMode>,
);
