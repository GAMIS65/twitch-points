import { Clock } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";

// Sample data for recent entries
const recentEntries = [
  {
    username: "TwitchUser123",
    timestamp: "2 minutes ago",
    channel: "GamingPro",
  },
  {
    username: "StreamFan42",
    timestamp: "5 minutes ago",
    channel: "StreamQueen",
  },
  { username: "GamerPro99", timestamp: "7 minutes ago", channel: "GamingPro" },
  {
    username: "ViewerElite",
    timestamp: "12 minutes ago",
    channel: "RetroGamer",
  },
  {
    username: "LuckyCharm",
    timestamp: "15 minutes ago",
    channel: "StreamQueen",
  },
  {
    username: "TwitchFan88",
    timestamp: "20 minutes ago",
    channel: "MorningStreamer",
  },
  {
    username: "StreamerSub",
    timestamp: "25 minutes ago",
    channel: "GamingPro",
  },
  {
    username: "GamingWizard",
    timestamp: "30 minutes ago",
    channel: "RetroGamer",
  },
];

export function RecentEntries() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Clock className="h-5 w-5 text-purple-500" />
          Recent Entries
        </CardTitle>
        <CardDescription>Latest viewers who entered giveaways</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          {recentEntries.map((entry, index) => (
            <div
              key={index}
              className="flex items-center justify-between py-2 border-b last:border-0"
            >
              <div className="font-medium">{entry.username}</div>
              <div className="flex items-center gap-4">
                <span className="text-sm text-muted-foreground">
                  {entry.channel}
                </span>
                <span className="text-xs text-muted-foreground">
                  {entry.timestamp}
                </span>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
