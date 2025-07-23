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

// Static version that doesn't revalidate on focus
export function useStreamersStatic(options?: SWRConfiguration) {
  return useApiStatic<
    {
      username: string;
      twitch_id: string;
      profile_image_url: string;
      is_live: string;
    }[]
  >("/giveaway/streamers", options);
}

export function useStreamers(options?: SWRConfiguration) {
  return useApi<
    {
      username: string;
      twitch_id: string;
      profile_image_url: string;
      is_live: string;
    }[]
  >("/giveaway/streamers", options);
}

export function useTotalParticipantsStatic(options?: SWRConfiguration) {
  return useApiStatic<{
    total_participants: number;
  }>("/giveaway/participants-count", options);
}

export function useTotalParticipants(options?: SWRConfiguration) {
  return useApi<{
    total_participants: number;
  }>("/giveaway/participants-count", options);
}

export function useTotalEntriesStatic(options?: SWRConfiguration) {
  return useApiStatic<{
    total_entries: number;
  }>("/giveaway/entries-count", options);
}

export function useTotalEntries(options?: SWRConfiguration) {
  return useApi<{
    total_entries: number;
  }>("/giveaway/entries-count", options);
}

export function useLeaderboardStatic(options?: SWRConfiguration) {
  return useApiStatic<
    {
      total_redemptions: number;
      username: string;
    }[]
  >("/giveaway/leaderboard", options);
}

export function useLeaderboard(options?: SWRConfiguration) {
  return useApi<
    {
      total_redemptions: number;
      username: string;
    }[]
  >("/giveaway/leaderboard", options);
}

export function useRecentEntries(options?: SWRConfiguration) {
  return useApi<
    {
      message_id: string;
      redeemed_at: string;
      streamer_username: string;
      viewer_username: string;
    }[]
  >("/giveaway/recent-entries", options);
}

// not ideal but should be fine for now
export function useUser(options?: SWRConfiguration) {
  return useApiStatic<
    {
      twitch_id: string;
      username: string;
      profile_image_url: string;
    }[]
  >("/me", { shouldRetryOnError: false, ...options });
}
