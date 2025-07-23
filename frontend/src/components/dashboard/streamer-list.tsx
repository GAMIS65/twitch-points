import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useStreamers } from "@/hooks/use-api";
import { Users } from "lucide-react";
import { Badge } from "@/components/ui/badge";

export function StreamersList() {
  const { data: streamers, error, isLoading } = useStreamers();

  const sortedStreamers = streamers?.slice().sort((a, b) => {
    return Number(b.is_live) - Number(a.is_live);
  });

  return (
    <Card className="bg-white/70 backdrop-blur-sm border-purple-200">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-purple-800">
          <Users className="h-5 w-5 text-purple-600" />
          Host Streamers
        </CardTitle>
        <CardDescription className="text-purple-600">
          Streamers hosting this giveaway
        </CardDescription>
      </CardHeader>

      <CardContent className="p-0">
        {streamers && sortedStreamers && !isLoading && !error && (
          <div className="max-h-96 overflow-y-auto px-6 pb-6">
            <div className="space-y-4">
              {sortedStreamers.map((streamer) => (
                <div
                  key={streamer.username}
                  className="flex items-center gap-3 p-3 rounded-lg bg-gradient-to-r from-indigo-50 to-purple-50 hover:from-indigo-100 hover:to-purple-100 transition-colors duration-100 "
                >
                  <div className="relative">
                    <Avatar className="h-12 w-12">
                      <AvatarImage
                        src={streamer.profile_image_url || "/placeholder.svg"}
                      />
                      <AvatarFallback className="bg-purple-200 text-purple-700">
                        {streamer.username.slice(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                  </div>
                  <div className="flex-1">
                    <p className="font-medium text-purple-900">
                      {streamer.username}
                    </p>
                  </div>
                  {streamer.is_live && (
                    <Badge className="bg-red-500 hover:bg-red-600">LIVE</Badge>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
