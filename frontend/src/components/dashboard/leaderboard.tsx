import { Trophy, RefreshCcw, AlertCircle } from "lucide-react";
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
import { useLeaderboard, useTotalEntries } from "../../hooks/use-api";
import { Button } from "../ui/button";
import { Skeleton } from "../ui/skeleton";

export function LeaderboardTable() {
  const {
    data: leaderboardUsers,
    error,
    isLoading,
    mutate: mutateLeaderboard,
  } = useLeaderboard();

  const { data: entriesCount, mutate: mutateEntries } = useTotalEntries();

  // Function to refresh all data sources
  const refreshData = () => {
    mutateLeaderboard();
    mutateEntries();
  };

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
    <Card>
      <CardHeader className="flex flex-row items-start justify-between pb-2">
        <div>
          <CardTitle className="flex items-center gap-2">
            <Trophy className="h-5 w-5 text-purple-500" />
            Leaderboard
          </CardTitle>
          <CardDescription>
            Users with the most entries across all channels
          </CardDescription>
        </div>
        <Button
          variant="outline"
          size="icon"
          onClick={refreshData}
          className="h-8 w-8"
        >
          <RefreshCcw className={`h-4 w-4`} />
          <span className="sr-only">Refresh</span>
        </Button>
      </CardHeader>
      <CardContent className="p-0">
        {isLoading && (
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
                {Array(10)
                  .fill(0)
                  .map((_, i) => (
                    <TableRow key={i}>
                      <TableCell className="font-medium">
                        <Skeleton className="h-4 w-8" />
                      </TableCell>
                      <TableCell>
                        <Skeleton className="h-4 w-24" />
                      </TableCell>
                      <TableCell>
                        <Skeleton className="h-4 w-16" />
                      </TableCell>
                      <TableCell className="text-right">
                        <Skeleton className="h-4 w-12" />
                      </TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>
          </div>
        )}

        {error && (
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <AlertCircle className="h-12 w-12 text-red-500 mb-3" />
            <h3 className="text-lg font-semibold mb-1">
              Failed to load leaderboard
            </h3>
            <p className="text-muted-foreground mb-4">
              There was an error loading the leaderboard.
            </p>
            <Button variant="outline" onClick={refreshData} className="gap-2">
              <RefreshCcw className="h-4 w-4" />
              Try Again
            </Button>
          </div>
        )}

        {leaderboardUsers && !isLoading && !error && entriesCount && (
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
                {leaderboardUsers.map((user, index) => (
                  <TableRow key={user.username}>
                    <TableCell className="font-medium">#{index + 1}</TableCell>
                    <TableCell>{user.username}</TableCell>
                    <TableCell>{user.total_redemptions}</TableCell>
                    <TableCell className="text-right">
                      {calculateChanceToWin(
                        user.total_redemptions,
                        entriesCount.total_entries,
                      ).toFixed(2)}
                      %
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
