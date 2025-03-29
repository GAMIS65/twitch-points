import { Trophy } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/table";

// Sample data
const leaderboardUsers = [
  { rank: 1, name: "TwitchUser123", entries: 87, winPercentage: 2.3 },
  { rank: 2, name: "StreamFan42", entries: 76, winPercentage: 1.3 },
  { rank: 3, name: "GamerPro99", entries: 65, winPercentage: 0 },
  { rank: 4, name: "ViewerElite", entries: 54, winPercentage: 1.9 },
  { rank: 5, name: "LuckyCharm", entries: 43, winPercentage: 7.0 },
  { rank: 6, name: "TwitchFan88", entries: 38, winPercentage: 0 },
  { rank: 7, name: "StreamerSub", entries: 32, winPercentage: 3.1 },
  { rank: 8, name: "GamingWizard", entries: 28, winPercentage: 0 },
  { rank: 9, name: "PurplePanda", entries: 25, winPercentage: 0 },
  { rank: 10, name: "NightOwl", entries: 22, winPercentage: 4.5 },
  { rank: 11, name: "StreamerFan1", entries: 20, winPercentage: 0 },
  { rank: 12, name: "GamerGirl42", entries: 18, winPercentage: 5.5 },
  { rank: 13, name: "TwitchViewer", entries: 15, winPercentage: 0 },
  { rank: 14, name: "StreamerSub2", entries: 12, winPercentage: 0 },
  { rank: 15, name: "GamingWizard2", entries: 10, winPercentage: 10.0 },
];

export function LeaderboardTable() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Trophy className="h-5 w-5 text-purple-500" />
          Leaderboard
        </CardTitle>
        <CardDescription>
          Users with the most entries across all channels
        </CardDescription>
      </CardHeader>
      <CardContent className="p-0">
        <div className="max-h-[400px] overflow-y-auto">
          <Table>
            <TableHeader className="sticky top-0 bg-card z-10">
              <TableRow>
                <TableHead className="w-12">Rank</TableHead>
                <TableHead>User</TableHead>
                <TableHead>Total Entries</TableHead>
                <TableHead className="text-right">Win %</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {leaderboardUsers.map((user) => (
                <TableRow key={user.name}>
                  <TableCell className="font-medium">#{user.rank}</TableCell>
                  <TableCell>{user.name}</TableCell>
                  <TableCell>{user.entries}</TableCell>
                  <TableCell className="text-right">
                    {user.winPercentage.toFixed(1)}%
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}
