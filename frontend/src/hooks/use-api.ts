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
    }[]
  >("/giveaway/streamers", options);
}
