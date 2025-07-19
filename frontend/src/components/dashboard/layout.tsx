import { Outlet } from "react-router";
import { Navbar } from "@/components/dashboard/navbar";

export function Layout() {
  return (
    <div>
      <div className="flex justify-center">
        <Navbar />
      </div>

      <main className="pt-12">
        <Outlet />
      </main>
    </div>
  );
}
