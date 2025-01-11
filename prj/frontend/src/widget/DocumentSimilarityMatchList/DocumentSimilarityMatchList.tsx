import { FC } from "react";
import { components } from "shared/api/api-types";
import "./style.css";

type DocumentSimilarityMatch = components["schemas"]["AnalyzedDocumentMatch"];

type DocumentSimilarityMatchListProps = {
    documents: DocumentSimilarityMatch[];
    loading: boolean;
    error: string | null;
    onDocumentClick?: (id: string) => void;
};

export const DocumentSimilarityMatchList: FC<
    DocumentSimilarityMatchListProps
> = ({ documents, loading, error, onDocumentClick }) => {
    if (onDocumentClick == undefined) {
        onDocumentClick = (_: string) => { }; // eslint-disable-line
    }

    if (loading) {
        return <p className="loading">Loading...</p>;
    }

    if (error) {
        return <p className="error">{error}</p>;
    }

    return (
        <div className="document-list">
            {documents.length > 0 ? (
                <ul className="styled-list">
                    {documents.map((doc) => (
                        <li key={doc.id} className="styled-item">
                            <DocumentSimilarityMatchListItem
                                doc={doc}
                                onClick={onDocumentClick}
                            />
                        </li>
                    ))}
                </ul>
            ) : (
                <p>No documents found</p>
            )}
        </div>
    );
};

export const DocumentSimilarityMatchListItem: FC<{
    doc: DocumentSimilarityMatch;
    onClick: (id: string) => void;
}> = ({ doc, onClick }) => {
    return (
        <div className="document-card" onClick={() => onClick(doc.id || "")}>
            <h3 className="doc-title">{doc.id}</h3>
            <p className="doc-description">{doc.rate}</p>
        </div>
    );
};
