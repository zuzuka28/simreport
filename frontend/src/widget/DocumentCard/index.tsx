import { FC } from "react";
import { components } from "shared/api/api-types";

type Document = components["schemas"]["DocumentSummary"];

type DocumentCardParams = {
  document: Document;
};

export const DocumentCard: FC<DocumentCardParams> = ({ document }) => (
  <div className="document-card">
    <h2>{document.name}</h2>
    <p>Last updated: {document.lastUpdated}</p>
    <button onClick={() => alert("Download functionality coming soon!")}>
      Download
    </button>
  </div>
);
