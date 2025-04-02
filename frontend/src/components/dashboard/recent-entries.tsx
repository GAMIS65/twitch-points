import { AlertCircle, Clock, RefreshCcw } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { useRecentEntries } from "@/hooks/use-api";
import { Skeleton } from "../ui/skeleton";
import { Button } from "../ui/button";

export function RecentEntries() {
  const { data: recentEntries, error, isLoading, mutate } = useRecentEntries();
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
        {isLoading && (
          <div className="space-y-2">
            {Array(5)
              .fill(0)
              .map((_, i) => (
                <div
                  key={i}
                  className="flex items-center justify-between py-2 border-b last:border-0"
                >
                  <Skeleton className="h-5 w-20" />
                  <div className="flex items-center gap-4">
                    <Skeleton className="h-5 w-16" />
                    <Skeleton className="h-5 w-12" />
                  </div>
                </div>
              ))}
          </div>
        )}

        {error && (
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <AlertCircle className="h-12 w-12 text-red-500 mb-3" />
            <h3 className="text-lg font-semibold mb-1">
              Failed to load recent entries
            </h3>
            <p className="text-muted-foreground mb-4">
              There was an error loading the recent entries list.
            </p>
            <Button
              variant="outline"
              onClick={() => mutate()}
              className="gap-2"
            >
              <RefreshCcw className="h-4 w-4" />
              Try Again
            </Button>
          </div>
        )}

        {recentEntries && !isLoading && !error && (
          <div className="space-y-2">
            {recentEntries.length > 0 ? (
              recentEntries.map((entry, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between py-2 border-b last:border-0"
                >
                  <div className="font-medium">{entry.viewer_username}</div>
                  <div className="flex items-center gap-4">
                    <span className="text-sm text-muted-foreground">
                      {entry.streamer_username}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {formatRelativeTime(entry.redeemed_at)}
                    </span>
                  </div>
                </div>
              ))
            ) : (
              <div className="text-center py-6 text-muted-foreground">
                No recent entries found.
              </div>
            )}
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
