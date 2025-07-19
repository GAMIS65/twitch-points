import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useTotalEntriesStatic } from "@/hooks/use-api";
import { Ticket } from "lucide-react";

export default function TotalEntries() {
  const { data: totalEntries } = useTotalEntriesStatic();

  return (
    <Card className="bg-gradient-to-r from-pink-100 to-pink-200 border-pink-200">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-pink-800">
          Total Entries
        </CardTitle>
        <Ticket className="h-4 w-4 text-pink-600" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold text-pink-900">
          {totalEntries ? totalEntries.toLocaleString() : "0"}
        </div>
        <p className="text-xs text-pink-700">entries redeemed</p>
      </CardContent>
    </Card>
  );
}
