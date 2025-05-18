import { Header } from "./components/dashboard/header";
import { StatsCards } from "./components/dashboard/stats-card";
import { StreamersList } from "./components/dashboard/streamer-list";
import { LeaderboardTable } from "./components/dashboard/leaderboard";
import { RecentEntries } from "./components/dashboard/recent-entries";
import { useTotalEntries, useTotalParticipants } from "./hooks/use-api";

export default function Dashboard() {
  const { data: participantsCount } = useTotalParticipants();
  const { data: entriesCount } = useTotalEntries();

  return (
    <div className="flex min-h-screen bg-background">
      <div className="flex flex-col w-full">
        <Header />
        <main className="flex-1 p-4 md:p-6">
          <div className="flex flex-col gap-4 md:gap-8">
            <div className="flex justify-end"></div>
            <StatsCards
              totalParticipants={participantsCount?.total_participants}
              totalEntries={entriesCount?.total_entries}
            />
            <div className="grid gap-4 md:grid-cols-2">
              <StreamersList />
              <LeaderboardTable />
            </div>
            <RecentEntries />
          </div>
        </main>
      </div>
    </div>
  );
}
