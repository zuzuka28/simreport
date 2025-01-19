import { FC, useState } from "react";
import { FileUploader } from "widget/FileUploader";
import { Dropdown } from "widget/Dropdown";
import client from "shared/api";
import "./style.css"

type DocumentUploadProps = {
  initialParentDocument?: string;
  initialGroups?: string[];
};

export const DocumentUploader: FC<DocumentUploadProps> = ({
  initialParentDocument,
  initialGroups = [],
}) => {
  const [selectedFileId, setSelectedFileId] = useState<string | null>(null);
  const [selectedGroups, setSelectedGroups] = useState<string[]>(initialGroups);

  const provideGroupIDs = async () => {
    const { items } = await client.attribute.fetchValues({
      attribute: "groupIDs",
    });

    const groups = items || [];

    return groups.map((doc: any) => ({ value: doc.id, label: doc.title }));
  };

  const handleFileUploaded = (fileId: string) => {
    setSelectedFileId(fileId);
  };

  const handleSubmit = async () => {
    if (!selectedFileId) {
      alert("Please complete all fields before submitting.");
      return;
    }

    try {
      await client.document.uploadDocument({
        fileID: selectedFileId,
        parentID: initialParentDocument,
        groupID: selectedGroups,
      });

      alert("Document uploaded successfully!");

      setSelectedFileId(null);
      setSelectedGroups(initialGroups);
    } catch (error) {
      console.error("Upload failed:", error);
      alert("Failed to upload document.");
    }
  };

  return (
    <div className="upload-container">
      <h2>Upload Document</h2>

      <div className="file-uploader">
        <label>Upload your file:</label>
        <FileUploader onDocumentUploaded={handleFileUploaded} />
      </div>

      {selectedFileId && (
        <p className="selected-file">Selected File ID: {selectedFileId}</p>
      )}

      <div className="dropdown">
        <label>Groups:</label>
        <Dropdown
          fetchOptions={provideGroupIDs}
          selectedValues={selectedGroups}
          onChange={setSelectedGroups}
          placeholder="Select Groups"
        />
      </div>

      <div className="submit-button">
        <button onClick={handleSubmit}>Upload Document</button>
      </div>
    </div>
  );
};
