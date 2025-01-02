import { FC } from "react";
import { components } from "shared/api/api-types";
import "./style.css";

type Document = components["schemas"]["DocumentSummary"];

type DocumentListProps = {
  documents: Document[];
  loading: boolean;
  error: string | null;
};

export const DocumentList: FC<DocumentListProps> = ({
  documents,
  loading,
  error,
}) => {
  if (loading) {
    return <p className="loading">Loading...</p>;
  }

  if (error) {
    return <p className="error">{error}</p>;
  }

  return (
    <div className="document-list">
      {documents.length > 0 ? (
        <ul>
          {documents.map((doc) => (
            <li key={doc.id}>{doc.name}</li>
          ))}
        </ul>
      ) : (
        <p>No documents found</p>
      )}
    </div>
  );
};
