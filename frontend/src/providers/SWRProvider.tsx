import type { ReactNode } from "react";
import { SWRConfig } from "swr";
import { swrOptions } from "../lib/swr-config";

interface SWRProviderProps {
  children: ReactNode;
}

export function SWRProvider({ children }: SWRProviderProps) {
  return <SWRConfig value={swrOptions}>{children}</SWRConfig>;
}
