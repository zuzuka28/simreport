import { FC, useState } from "react";
import { DocumentList } from "widget/DocumentList";
import { DocumentUploader } from "widget/DocumentUploader";
import { DocumentSimilarityMatchList } from "widget/DocumentSimilarityMatchList";
import { Popup } from "widget/Popup";
import { useDocuments, useSimilarityCheckDocuments } from "shared/hook";
import "./style.css";

export const SimilarityCheck: FC = () => {
  const [panelMinimized, setPanelMinimized] = useState(false);
  const [showPopup, setShowPopup] = useState(false);
  const {
    documents,
    loading: documents_loading,
    error: documents_error,
  } = useDocuments("");

  const [query, setQuery] = useState("");
  const {
    documents: matches,
    loading: matches_loading,
    error: matches_error,
  } = useSimilarityCheckDocuments(query);

  const togglePanel = () => {
    setPanelMinimized(!panelMinimized);
  };

  const openPopup = () => {
    setShowPopup(true);
  };

  const closePopup = () => {
    setShowPopup(false);
  };

  const handleDocumentSelect = (document_id: string) => {
    setQuery(document_id);
    setPanelMinimized(true);
    closePopup();
  };

  return (
    <>
      <div
        className={`document_manager_container ${panelMinimized ? "minimized" : ""}`}
      >
        {!panelMinimized && (
          <>
            <DocumentUploader onDocumentUploaded={handleDocumentSelect} />
            <button className="popup-button" onClick={openPopup}>
              Select Existing File
            </button>
          </>
        )}

        {panelMinimized && (
          <button className="toggle-button" onClick={togglePanel}>
            Expand Panel
          </button>
        )}
      </div>

      {showPopup && (
        <Popup onClose={closePopup}>
          <h2>Select a File</h2>
          <DocumentList
            documents={documents}
            loading={documents_loading}
            error={documents_error}
            onDocumentClick={handleDocumentSelect}
          />
        </Popup>
      )}

      <DocumentSimilarityMatchList
        documents={matches}
        loading={matches_loading}
        error={matches_error}
      />
    </>
  );
};
