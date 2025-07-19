import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useTotalParticipantsStatic } from "@/hooks/use-api";
import { Users } from "lucide-react";

export default function TotalParticipants() {
  const { data: participantsCount } = useTotalParticipantsStatic();

  return (
    <Card className="bg-gradient-to-r from-purple-100 to-purple-200 border-purple-200">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-purple-800">
          Total Participants
        </CardTitle>
        <Users className="h-4 w-4 text-purple-600" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold text-purple-900">
          {participantsCount ? participantsCount.toLocaleString() : "0"}
        </div>
        <p className="text-xs text-purple-700">viewers participating</p>
      </CardContent>
    </Card>
  );
}
