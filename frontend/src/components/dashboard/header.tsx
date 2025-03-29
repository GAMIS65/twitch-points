import { Link } from "react-router";
import { TwitchIcon } from "lucide-react";
import { Button } from "../ui/button";

export function Header() {
  return (
    <header className="sticky top-0 z-10 flex h-14 items-center justify-between border-b bg-background px-4 sm:px-6">
      <div className="flex items-center gap-2">
        <div className="md:hidden">
          <TwitchIcon className="h-6 w-6 text-purple-500" />
        </div>
        <Link to="/" className="text-lg font-medium">
          GAMIS65's giveaway tool
        </Link>
      </div>
      <Link to="/signin">
        <Button
          className="gap-1 bg-purple-600 hover:bg-purple-700"
          variant="default"
        >
          <TwitchIcon className="h-4 w-4" />
          Streamer Sign In
        </Button>
      </Link>
    </header>
  );
}
