import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Trophy, Crown, Star } from "lucide-react";
import { useLeaderboard, useTotalEntries } from "@/hooks/use-api";

export function Leaderboard() {
  const { data: leaderboardUsers, error, isLoading } = useLeaderboard();
  const { data: entriesCount } = useTotalEntries();

  const calculateChanceToWin = (
    entries: number,
    totalEntries: number,
  ): number => {
    if (totalEntries === 0) {
      return 0;
    }
    return (entries / totalEntries) * 100;
  };

  return (
    <Card className="bg-white/70 backdrop-blur-sm border-purple-200">
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-purple-800">
          <Trophy className="h-5 w-5 text-yellow-500" />
          Top Participants
        </CardTitle>
        <CardDescription className="text-purple-600">
          Viewers with the most entries
        </CardDescription>
      </CardHeader>
      <CardContent className="p-0">
        {leaderboardUsers && !isLoading && !error && entriesCount && (
          <div className="max-h-96 overflow-y-auto px-6 pb-6">
            <div className="space-y-3">
              {leaderboardUsers.map((user, index) => (
                <div
                  key={user.username}
                  className="flex items-center gap-3 p-3 rounded-lg bg-gradient-to-r from-purple-50 to-pink-50 hover:from-purple-100 hover:to-pink-100 duration-100 transition-colors "
                >
                  <div className="flex items-center gap-2 min-w-0">
                    {index === 0 && (
                      <Crown className="h-4 w-4 text-yellow-500 flex-shrink-0" />
                    )}
                    {index === 1 && (
                      <Star className="h-4 w-4 text-gray-400 flex-shrink-0" />
                    )}
                    {index === 2 && (
                      <Star className="h-4 w-4 text-amber-600 flex-shrink-0" />
                    )}
                    <span className="text-sm font-medium text-purple-700 flex-shrink-0">
                      #{index + 1}
                    </span>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-purple-900 truncate">
                      {user.username}
                    </p>
                    <p className="text-xs text-purple-600">
                      {calculateChanceToWin(
                        user.total_redemptions,
                        entriesCount.total_entries,
                      )}
                      % chance to win
                    </p>
                  </div>
                  <div className="flex flex-col items-end gap-1 flex-shrink-0">
                    <Badge
                      variant="secondary"
                      className="bg-purple-100 text-purple-800"
                    >
                      {user.total_redemptions} entries
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
