import { TwitchIcon } from "lucide-react";
import { Button } from "./components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "./components/ui/card";
import { Header } from "./components/dashboard/header";

export default function AddRewardPage() {
  const baseUrl = import.meta.env.VITE_BACKEND_URL || "";

  async function addReward() {
    try {
      const response = await fetch(baseUrl + "/add-reward", {
        method: "POST",
        credentials: "include",
        mode: "cors",
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Failed to add reward: ${errorText}`);
      }

      return await response.text();
    } catch (error) {
      console.error("Error adding reward:", error);
      throw error;
    }
  }

  const handleClick = async () => {
    try {
      const result = await addReward();
      alert("Reward added successfully!");
      console.log(result);
    } catch (error) {
      alert("Failed to add reward. See console for details.");
    }
  };

  return (
    <div className="flex min-h-screen flex-col bg-background">
      <Header />
      <main className="flex-1">
        <div className="container mx-auto max-w-3xl py-8">
          <Card className="text-center">
            <CardHeader>
              <CardTitle className="text-2xl">
                Add Giveaway Reward to Your Channel
              </CardTitle>
            </CardHeader>

            <CardContent className="space-y-6">
              <div className="flex justify-center">
                <Button
                  className="bg-purple-600 hover:bg-purple-700 w-full max-w-md"
                  size="lg"
                  onClick={handleClick}
                >
                  <TwitchIcon className="mr-2 h-5 w-5" />
                  Add Reward to Channel
                </Button>
              </div>

              <div className="text-sm text-muted-foreground">
                <p>
                  If you encounter any issues, please contact GAMIS65 on Discord
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
}
