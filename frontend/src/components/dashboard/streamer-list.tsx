import { AlertCircle, RefreshCcw, Tv } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { useStreamersStatic } from "../../hooks/use-api";

import { Button } from "../ui/button";
import { Skeleton } from "../ui/skeleton";

export function StreamersList() {
  const { data: streamers, error, isLoading, mutate } = useStreamersStatic();

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Tv className="h-5 w-5 text-purple-500" />
          Streamers hosting the giveaway
        </CardTitle>
        <CardDescription>
          Check out these streamers to gain more entries!
        </CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading && (
          <div className="grid gap-4 md:grid-cols-2">
            {Array(4)
              .fill(0)
              .map((_, i) => (
                <div
                  key={i}
                  className="flex items-center gap-3 rounded-lg border p-3"
                >
                  <Skeleton className="h-8 w-8 rounded-full" />
                  <Skeleton className="h-5 w-32" />
                </div>
              ))}
          </div>
        )}

        {error && (
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <AlertCircle className="h-12 w-12 text-red-500 mb-3" />
            <h3 className="text-lg font-semibold mb-1">
              Failed to load streamers
            </h3>
            <p className="text-muted-foreground mb-4">
              There was an error loading the streamer list.
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

        {streamers && !isLoading && !error && (
          <div className="grid gap-4 md:grid-cols-2">
            {streamers.length > 0 ? (
              streamers.map((streamer) => (
                <div
                  key={streamer.username}
                  className="flex items-center gap-3 rounded-lg border p-3"
                >
                  <Avatar className="h-8 w-8 border-2 border-purple-500">
                    <AvatarImage
                      src={streamer.profile_image_url}
                      alt={streamer.username}
                    />
                    <AvatarFallback>
                      {streamer.username.substring(0, 2)}
                    </AvatarFallback>
                  </Avatar>
                  <div className="flex items-center gap-2">
                    <span className="font-medium">{streamer.username}</span>
                  </div>
                </div>
              ))
            ) : (
              <div className="col-span-2 text-center py-6 text-muted-foreground">
                No streamers are currently hosting the giveaway.
              </div>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
