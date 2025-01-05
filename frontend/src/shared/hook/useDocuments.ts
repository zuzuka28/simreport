import { useState, useEffect } from "react";
import client from "shared/api";
import { components } from "shared/api/api-types";

type Document = components["schemas"]["DocumentSummary"];

export const useDocuments = (query: string) => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchDocuments = async () => {
      setLoading(true);
      try {
        const response = await client.document.searchDocuments({ name: query });
        setDocuments(response.documents || []);
        setError(null);
      } catch (error) {
        setError("Failed to fetch documents.");
      } finally {
        setLoading(false);
      }
    };

    void fetchDocuments();
  }, [query]);

  return { documents, loading, error };
};

type DocumentSimilarityMatch = components["schemas"]["AnalyzedDocumentMatch"];

export const useSimilarityCheckDocuments = (id: string) => {
  const [documents, setDocuments] = useState<DocumentSimilarityMatch[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchDocuments = async () => {
      setLoading(true);
      try {
        const response = await client.analyze.findSimilarDocuments(id);
        setDocuments(response.documents || []);
        setError(null);
      } catch (error) {
        setError("Failed to fetch documents.");
      } finally {
        setLoading(false);
      }
    };

    void fetchDocuments();
  }, [id]);

  return { documents, loading, error };
};
