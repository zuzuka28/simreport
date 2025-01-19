import { FC, useState } from "react";
import { DocumentList } from "widget/DocumentList";
import { FileUploader } from "widget/FileUploader";
import { DocumentSimilarityMatchList } from "widget/DocumentSimilarityMatchList";
import { Popup } from "widget/Popup";
import { useDocuments, useSimilarityCheckDocuments } from "shared/hook";
import "./style.css";

export const SimilarityCheck: FC = () => {
  const [panelMinimized, setPanelMinimized] = useState(false);
  const [showPopup, setShowPopup] = useState(false);

  const [dquery] = useState<{ groupID: string[], name: string }>({
    groupID: [],
    name: "",
  });

  const { data: documentsData, loading: documentsLoading, error: documentsError } = useDocuments(dquery);

  const [query, setQuery] = useState("");
  const {
    data: matchesData,
    loading: matchesLoading,
    error: matchesError,
  } = useSimilarityCheckDocuments(query);

  let matches = matchesData?.documents || []

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
            <FileUploader onDocumentUploaded={handleDocumentSelect} />
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
            documents={documentsData?.documents || []}
            loading={documentsLoading}
            error={documentsError}
            onDocumentClick={handleDocumentSelect}
          />
        </Popup>
      )}

      <DocumentSimilarityMatchList
        documents={matches}
        loading={matchesLoading}
        error={matchesError}
      />
    </>
  );
};
