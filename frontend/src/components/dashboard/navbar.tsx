import { Home, Settings, TwitchIcon, Crown } from "lucide-react";
import { ReactNode } from "react";
import { useLocation, useNavigate } from "react-router";

type NavItemProps = {
  children: ReactNode;
  isActive: boolean;
  path: string;
};

type NavItem = {
  name: string;
  path: string;
  icon: ReactNode;
};

export function Navbar() {
  const location = useLocation();

  const navItems: NavItem[] = [
    { name: "Home", path: "/", icon: <Home size={20} /> },
    { name: "Pick a winner", path: "/wheel", icon: <Crown size={20} /> },
    { name: "Settings", path: "/settings", icon: <Settings size={20} /> },
    {
      name: "Streamer Sign-In",
      path: "/sign-in",
      icon: <TwitchIcon size={20} />,
    },
  ];

  return (
    <nav className="fixed top-8 flex items-center justify-center p-2 bg-white/20 backdrop-blur-lg rounded-full shadow-lg border border-white/30 z-10">
      <ul className="flex items-center gap-2">
        {navItems.map((item) => (
          <NavItem
            key={item.name}
            isActive={location.pathname === item.path}
            path={item.path}
          >
            {item.icon}
            <span className="hidden sm:inline-block">{item.name}</span>
          </NavItem>
        ))}
      </ul>
    </nav>
  );
}

function NavItem({ children, isActive, path }: NavItemProps) {
  const navigate = useNavigate();

  return (
    <li
      onClick={() => navigate(path)}
      className={`
        flex items-center justify-center gap-2 px-4 py-2 rounded-full cursor-pointer transition-all duration-100 ease-in-out
        ${
          isActive
            ? "bg-white text-purple-600 font-semibold shadow-md"
            : "text-gray-700 hover:shadow-md hover:bg-gray-200/30"
        }
      `}
    >
      {children}
    </li>
  );
}
