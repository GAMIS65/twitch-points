import { Tv } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";

// Sample data
const streamers = [
  { name: "GamingPro", status: "Live" },
  { name: "StreamQueen", status: "Live" },
  { name: "RPGMaster", status: "Offline" },
  { name: "SpeedRunner", status: "Live" },
  { name: "CosplayGirl", status: "Offline" },
  { name: "RetroGamer", status: "Live" },
  { name: "MorningStreamer", status: "Live" },
  { name: "NightOwlGaming", status: "Offline" },
];

export function StreamersList() {
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
        <div className="grid gap-4 md:grid-cols-2">
          {streamers.map((streamer) => (
            <div
              key={streamer.name}
              className="flex items-center gap-3 rounded-lg border p-3"
            >
              <Avatar className="h-8 w-8 border-2 border-purple-500">
                <AvatarImage
                  src={`/placeholder.svg?height=32&width=32`}
                  alt={streamer.name}
                />
                <AvatarFallback>{streamer.name.substring(0, 2)}</AvatarFallback>
              </Avatar>
              <div className="flex items-center gap-2">
                <span className="font-medium">{streamer.name}</span>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
