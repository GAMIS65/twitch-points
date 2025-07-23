import useSWR, { type SWRConfiguration } from "swr";
import { buildApiUrl } from "../lib/swr-config";

// Base API hook with configurable options
export function useApi<T>(path: string | null, options?: SWRConfiguration) {
  const apiUrl = path ? buildApiUrl(path) : null;

  return useSWR<T>(apiUrl, options);
}

// Hook that doesn't revalidate on focus
export function useApiStatic<T>(
  path: string | null,
  options?: SWRConfiguration,
) {
  return useApi<T>(path, {
    revalidateOnFocus: false,
    ...options,
  });
}

export type Streamer = {
  username: string;
  twitch_id: string;
  profile_image_url: string;
  is_live: string;
};

export type TotalParticipants = {
  total_participants: number;
};

export type TotalEntries = {
  total_entries: number;
};

export type LeaderboardEntry = {
  total_redemptions: number;
  username: string;
};

export type RecentEntry = {
  message_id: string;
  redeemed_at: string;
  streamer_username: string;
  viewer_username: string;
};

export type User = {
  twitch_id: string;
  username: string;
  profile_image_url: string;
};

const API_ENDPOINTS = {
  STREAMERS: "/giveaway/streamers",
  PARTICIPANTS_COUNT: "/giveaway/participants-count",
  ENTRIES_COUNT: "/giveaway/entries-count",
  LEADERBOARD: "/giveaway/leaderboard",
  RECENT_ENTRIES: "/giveaway/recent-entries",
  ME: "/me",
} as const;

export function useStreamersStatic(options?: SWRConfiguration) {
  return useApiStatic<Streamer[]>(API_ENDPOINTS.STREAMERS, options);
}

export function useStreamers(options?: SWRConfiguration) {
  return useApi<Streamer[]>(API_ENDPOINTS.STREAMERS, options);
}

export function useTotalParticipantsStatic(options?: SWRConfiguration) {
  return useApiStatic<TotalParticipants>(
    API_ENDPOINTS.PARTICIPANTS_COUNT,
    options,
  );
}

export function useTotalParticipants(options?: SWRConfiguration) {
  return useApi<TotalParticipants>(API_ENDPOINTS.PARTICIPANTS_COUNT, options);
}

export function useTotalEntriesStatic(options?: SWRConfiguration) {
  return useApiStatic<TotalEntries>(API_ENDPOINTS.ENTRIES_COUNT, options);
}

export function useTotalEntries(options?: SWRConfiguration) {
  return useApi<TotalEntries>(API_ENDPOINTS.ENTRIES_COUNT, options);
}

export function useLeaderboardStatic(options?: SWRConfiguration) {
  return useApiStatic<LeaderboardEntry[]>(API_ENDPOINTS.LEADERBOARD, options);
}

export function useLeaderboard(options?: SWRConfiguration) {
  return useApi<LeaderboardEntry[]>(API_ENDPOINTS.LEADERBOARD, options);
}

export function useRecentEntries(options?: SWRConfiguration) {
  return useApi<RecentEntry[]>(API_ENDPOINTS.RECENT_ENTRIES, options);
}

// not ideal but should be fine for now
export function useUser(options?: SWRConfiguration) {
  return useApiStatic<User[]>(API_ENDPOINTS.ME, {
    shouldRetryOnError: false,
    ...options,
  });
}
