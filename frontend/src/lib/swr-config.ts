class FetchError extends Error {
  info?: any;
  status?: number;

  constructor(message: string, status?: number, info?: any) {
    super(message);
    this.name = "FetchError";
    this.status = status;
    this.info = info;
  }
}

export const defaultFetcher = async (url: string) => {
  const response = await fetch(url);

  if (!response.ok) {
    const info = await response.json().catch(() => null);
    throw new FetchError(
      "An error occurred while fetching the data.",
      response.status,
      info,
    );
  }

  return response.json();
};

export const swrOptions = {
  fetcher: defaultFetcher,
  revalidateOnFocus: true,
  revalidateOnReconnect: true,
  refreshInterval: 0,
  shouldRetryOnError: true,
  dedupingInterval: 2000,
  errorRetryCount: 3,
  errorRetryInterval: 5000,
};

export const buildApiUrl = (path: string) => {
  const baseUrl = import.meta.env.VITE_BACKEND_URL || "";
  return `${baseUrl}${path.startsWith("/") ? path : `/${path}`}`;
};
