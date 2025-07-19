import TotalParticipants from "@/components/dashboard/total-participants";
import TotalEntries from "@/components/dashboard/total-entries";
import { Leaderboard } from "@/components/dashboard/leaderboard";
import { StreamersList } from "@/components/dashboard/streamer-list";
import { RecentEntries } from "@/components/dashboard/recent-entries";

export default function Dashboard() {
  return (
    <div className="min-h-screen flex justify-center">
      <div className="container mx-auto p-6 space-y-6 pt-28">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <TotalParticipants />
          <TotalEntries />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Leaderboard />
          <StreamersList />
        </div>

        <RecentEntries />
      </div>
    </div>
  );
}
