import { Heart } from "lucide-react";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { useRecentEntries } from "@/hooks/use-api";
import { Badge } from "@/components/ui/badge";

// TODO: Amount of entries gained
export function RecentEntries() {
  const { data: recentEntries, error, isLoading } = useRecentEntries();
  return (
    <Card className="bg-white/70 backdrop-blur-sm border-purple-200">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-purple-800">
          <Heart className="h-5 w-5 text-pink-500" />
          Recent Entries
        </CardTitle>
        <CardDescription className="text-purple-600">
          Viewers who recently redeemed an entry
        </CardDescription>
      </CardHeader>
      <CardContent>
        {recentEntries && !isLoading && !error && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {recentEntries.map((user) => (
              <div
                key={user.viewer_username}
                className="flex items-center gap-3 p-3 rounded-lg bg-gradient-to-r from-pink-50 to-purple-50 border border-pink-100"
              >
                <div className="flex-1 min-w-0">
                  <p className="font-medium text-purple-900 truncate">
                    {user.viewer_username}
                  </p>
                  <p className="text-xs text-purple-600">
                    {formatRelativeTime(user.redeemed_at)}
                  </p>
                </div>

                <Badge
                  variant="outline"
                  className="border-purple-200 text-purple-700"
                >
                  1
                </Badge>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function formatRelativeTime(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffInSeconds = Math.round((now.getTime() - date.getTime()) / 1000);

  if (diffInSeconds < 60) {
    return `${diffInSeconds} second${diffInSeconds !== 1 ? "s" : ""} ago`;
  } else if (diffInSeconds < 3600) {
    const diffInMinutes = Math.floor(diffInSeconds / 60);
    return `${diffInMinutes} minute${diffInMinutes !== 1 ? "s" : ""} ago`;
  } else if (diffInSeconds < 86400) {
    const diffInHours = Math.floor(diffInSeconds / 3600);
    return `${diffInHours} hour${diffInHours !== 1 ? "s" : ""} ago`;
  } else if (diffInSeconds < 604800) {
    // 7 days
    const diffInDays = Math.floor(diffInSeconds / 86400);
    return `${diffInDays} day${diffInDays !== 1 ? "s" : ""} ago`;
  } else if (diffInSeconds < 2419200) {
    // 4 weeks
    const diffInWeeks = Math.floor(diffInSeconds / 604800);
    return `${diffInWeeks} week${diffInWeeks !== 1 ? "s" : ""} ago`;
  } else if (diffInSeconds < 29030400) {
    // 12 months
    const diffInMonths = Math.floor(diffInSeconds / 2419200);
    return `${diffInMonths} month${diffInMonths !== 1 ? "s" : ""} ago`;
  } else {
    const diffInYears = Math.floor(diffInSeconds / 29030400);
    return `${diffInYears} year${diffInYears !== 1 ? "s" : ""} ago`;
  }
}
