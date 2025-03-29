import { Link } from "react-router";
import { TwitchIcon, AlertCircle } from "lucide-react";
import { Button } from "./components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./components/ui/card";
import { Header } from "./components/dashboard/header";

export default function SignInPage() {
  const BACKEND_URL = import.meta.env.VITE_BACKEND_URL;
  return (
    <div className="flex min-h-screen flex-col bg-background">
      <Header />
      <main className="flex-1 flex items-center justify-center">
        <div className="container max-w-3xl py-8 mx-auto">
          <Card className="border-2 border-purple-200 dark:border-purple-900">
            <CardHeader className="pb-4">
              <CardTitle className="text-2xl flex items-center gap-2">
                <TwitchIcon className="h-6 w-6 text-purple-500" />
                Connect Your Twitch Channel
              </CardTitle>
            </CardHeader>

            <CardContent className="space-y-6">
              <div className="rounded-lg bg-amber-50 dark:bg-amber-950/50 p-4 border border-amber-200 dark:border-amber-900">
                <div className="flex gap-2 text-amber-800 dark:text-amber-300 font-medium mb-2">
                  <AlertCircle className="h-5 w-5 flex-shrink-0" />
                  <h3>Important Information</h3>
                </div>
                <ul className="space-y-2 text-amber-700 dark:text-amber-400 pl-7 list-disc">
                  <li>
                    This integration is{" "}
                    <strong>
                      only for Twitch streamers hosting the giveaway
                    </strong>
                    . If you're not a streamer, connecting won't provide any
                    functionality.
                  </li>
                  <li>
                    After connecting, your channel will have to be manually
                    verified before your viewers appear on the leaderboard. You
                    don't need to do anything, just wait a few hours at most.
                    Don't worry - all entries will still be counted during the
                    verification period, and nothing will be lost.
                  </li>
                  <li>
                    If you encounter any issues, please contact GAMIS65 on
                    Discord. The bugs won't fix themselves.
                  </li>
                </ul>
              </div>
            </CardContent>

            <CardFooter className="flex flex-col gap-4">
              <Link to={BACKEND_URL + "/auth/twitch"}>
                <Button
                  className="w-full bg-purple-600 hover:bg-purple-700"
                  size="lg"
                >
                  Connect Twitch channel
                </Button>
              </Link>
            </CardFooter>
          </Card>
        </div>
      </main>
    </div>
  );
}
