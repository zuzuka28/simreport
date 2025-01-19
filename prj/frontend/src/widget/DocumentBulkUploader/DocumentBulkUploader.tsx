import { FC, useState } from "react";
import { FileUploader } from "widget/FileUploader";
import client from "shared/api";
import "./style.css"
import { Tagger } from "widget/Tagger";

type DocumentBulkUploadProps = {
};

export const DocumentBulkUploader: FC<DocumentBulkUploadProps> = ({
}) => {
  const [selectedFileIds, setSelectedFileId] = useState<string[]>([]);
  const [selectedGroups, setSelectedGroups] = useState<string[]>([]);

  const provideGroupIDs = async () => {
    const { items } = await client.attribute.fetchValues({
      attribute: "groupIDs",
    });

    const groups = items || [];

    return groups.map((doc: any) => (doc.value));
  };

  const handleFileUploaded = (id: string) => {
    if (selectedFileIds.indexOf(id) === -1) {
      setSelectedFileId((prevFiles) => [...prevFiles, id]);
    }
  }

  const handleSubmit = async () => {
    for (const id of selectedFileIds) {
      try {
        await client.document.uploadDocument({
          fileID: id,
          parentID: "",
          groupID: selectedGroups,
        });
      } catch (error) {
        console.error("Upload failed:", error);
        alert("Failed to upload document.");
      }
    }

    setSelectedFileId([]);
    setSelectedGroups([]);
  };

  return (
    <div className="upload-container">
      <h2>Upload Document</h2>

      <div className="file-uploader">
        <label>Upload your file:</label>
        <FileUploader onDocumentUploaded={handleFileUploaded} />
      </div>

      {selectedFileIds && (
        <div className="selected-file">Selected File IDs: {selectedFileIds.map((v) => (<p>{v}</p>))}</div>
      )}

      <div className="dropdown">
        <label>Groups:</label>
        <Tagger
          initialValues={selectedGroups}
          onBubblesChange={setSelectedGroups}
          fetchSuggestions={provideGroupIDs}
        />
      </div>

      <div className="submit-button">
        <button onClick={handleSubmit}>Upload Document</button>
      </div>
    </div>
  );
};
