
import { useState, useEffect } from "react";

type FetchFunction<Q, R> = (query: Q) => Promise<  R  | undefined>;

interface UseDataHookResult<R> {
  data?: R;
  loading: boolean;
  error: string | null;
}

export function createDataHook<Q, R>(
  fetchFunction: FetchFunction<Q, R>,
) {
  return (query: Q): UseDataHookResult<R> => {
    const [data, setData] = useState<R>();
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
      const fetchData = async () => {
        setLoading(true);

        try {
          const response = await fetchFunction(query);
          if (response) {
            setData(response);
            setError(null);
          }
        } catch (err) {
          console.error(err)
          setError("can't load data");
        } finally {
          setLoading(false);
        }
      };

      void fetchData();
    }, [query]);

    return { data, loading, error };
  };
}