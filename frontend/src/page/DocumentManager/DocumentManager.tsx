import { FC, useState } from "react";
import { SearchBar } from "widget/SearchBar";
import { DocumentList } from "widget/DocumentList";
import { DocumentUploader } from "widget/DocumentUploader";
import "./style.css";
import { useDocuments } from "shared/hook";

export const DocumentManager: FC = () => {
  const [query, setQuery] = useState("");
  const { documents, loading, error } = useDocuments(query);

  const handleSearch = () => {
    if (query.trim() === "") {
      alert("Please enter a search term.");
    }
  };

  return (
    <>
      <div className="document_manager_container">
        <DocumentUploader />
        <div className="action-section">
          <SearchBar
            query={query}
            onQueryChange={setQuery}
            onSearch={handleSearch}
          />
        </div>
        <DocumentList documents={documents} loading={loading} error={error} />
      </div>
    </>
  );
};
