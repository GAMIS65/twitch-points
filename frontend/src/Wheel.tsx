import { useState } from "react";
import { Roulette, useRoulette } from "react-hook-roulette";
import Confetti from "react-confetti";
import { Button } from "./components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "./components/ui/card";
import { Header } from "./components/dashboard/header";
import { Play, StopCircle } from "lucide-react";
import { LeaderboardTable } from "./components/dashboard/leaderboard";
import { useLeaderboard } from "./hooks/use-api";

export default function WheelPage() {
  const [showConfetti, setShowConfetti] = useState(false);
  const [isSpinning, setIsSpinning] = useState(false);
  const { data: leaderboardUsers } = useLeaderboard();
  const [items] = useState(
    leaderboardUsers?.map((user) => ({
      name: user.username,
      weight: user.total_redemptions,
    })),
  );

  const options = {
    size: 500,
    deceleration: 0.02,
    maxSpeed: 15,
    determineAngle: 0,
    style: {
      label: {
        align: "left" as CanvasTextAlign,
      },
      canvas: {
        bg: "transparent",
      },
      arrow: {
        bg: "black",
        borderColor: "hsl(var(--primary))",
        borderWidth: 2,
      },
      pie: {
        border: false,
        borderColor: "transparent",
        borderWidth: 0,
        theme: [
          { bg: "oklch(85% 0.25 25)", color: "white" }, // Light pink
          { bg: "oklch(88% 0.18 45)", color: "black" }, // pastel orange
          { bg: "oklch(95% 0.1 100)", color: "black" }, // Pale yellow
          { bg: "oklch(90% 0.2 140)", color: "black" }, // Mint green
          { bg: "oklch(95% 0.15 210)", color: "black" }, // Baby blue
          { bg: "oklch(90% 0.18 250)", color: "white" }, // Periwinkle blue
          { bg: "oklch(85% 0.2 280)", color: "white" }, // Lavender
          { bg: "oklch(90% 0.25 320)", color: "white" }, // Light purple
        ],
      },
    },
  };

  // @ts-ignore
  const { roulette, onStart, onStop, result } = useRoulette({ items, options });

  const handleWheelAction = () => {
    if (isSpinning) {
      onStop();
      setShowConfetti(true);
      setIsSpinning(false);
    } else {
      setShowConfetti(false);
      onStart();
      setIsSpinning(true);
    }
  };

  return (
    <div className="flex min-h-screen bg-background">
      <div className="flex flex-col w-full">
        <Header />
        <main className="flex-1 p-4 md:p-6">
          <div className="flex flex-col gap-4 md:gap-8">
            <div className="grid gap-4 md:grid-cols-2">
              {/* Wheel Section */}
              <div className="flex flex-col items-center">
                <div className="relative mb-6">
                  {items && <Roulette roulette={roulette} />}
                  {showConfetti && result && (
                    <Confetti recycle={false} numberOfPieces={200} />
                  )}
                </div>

                <Button
                  onClick={handleWheelAction}
                  className="flex items-center gap-2 mb-6"
                  size="lg"
                  variant={isSpinning ? "destructive" : "default"}
                >
                  {isSpinning ? (
                    <>
                      <StopCircle className="h-4 w-4" />
                      <span>Stop Wheel</span>
                    </>
                  ) : (
                    <>
                      <Play className="h-4 w-4" />
                      <span>Start Wheel</span>
                    </>
                  )}
                </Button>

                {result && (
                  <Card className="w-full max-w-md">
                    <CardHeader>
                      <CardTitle className="text-center">Winner</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <p className="text-3xl font-bold text-center">{result}</p>
                    </CardContent>
                  </Card>
                )}
              </div>

              {/* Leaderboard Section */}
              <div>
                <LeaderboardTable />
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
