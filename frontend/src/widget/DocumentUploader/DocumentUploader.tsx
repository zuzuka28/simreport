import { useState, DragEvent, ChangeEvent, FormEvent, FC } from "react";
import "./style.css";
import client from "shared/api";

type DocumentUploaderProps = {
  onDocumentUploaded?: (documentID: string) => void;
};

export const DocumentUploader: FC<DocumentUploaderProps> = ({
  onDocumentUploaded,
}) => {
  const [files, setFiles] = useState<File[]>([]);
  const [dragActive, setDragActive] = useState<boolean>(false);

  if (onDocumentUploaded == undefined) {
    onDocumentUploaded = (_: string) => { }; // eslint-disable-line
  }

  const handleDragEnter = (e: DragEvent<HTMLFormElement>): void => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(true);
  };

  const handleDragLeave = (e: DragEvent<HTMLFormElement>): void => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
  };

  const handleDragOver = (e: DragEvent<HTMLFormElement>): void => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = (e: DragEvent<HTMLFormElement>): void => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) {
      const droppedFiles = Array.from(e.dataTransfer.files);
      setFiles((prevFiles) => [...prevFiles, ...droppedFiles]);
      e.dataTransfer.clearData();
    }
  };

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>): void => {
    if (e.target.files) {
      const selectedFiles = Array.from(e.target.files);
      setFiles((prevFiles) => [...prevFiles, ...selectedFiles]);
      e.target.value = "";
    }
  };

  const handleFileRemove = (index: number): void => {
    setFiles((prevFiles) => prevFiles.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    console.log("Files submitted:", files);

    for (const file of files) {
      const formData = new FormData();
      formData.append("document", file);

      try {
        const response = await client.document.uploadDocument(formData);
        onDocumentUploaded(response.documentID || "");
      } catch (error) {
        alert("Failed to upload the document.");
      }
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      onDragEnter={handleDragEnter}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      className={`file-upload-form ${dragActive ? "drag-active" : ""}`}
    >
      <p>Drag and drop your files here, or</p>
      <input
        type="file"
        multiple
        onChange={handleFileChange}
        style={{ display: "none" }}
        id="fileInput"
      />
      <label htmlFor="fileInput" className="file-upload-label">
        Choose files
      </label>
      {files.length > 0 && (
        <div className="file-list">
          <p>Files:</p>
          <ul>
            {files.map((file, index) => (
              <li key={index} className="file-item">
                {file.name}
                <button
                  type="button"
                  className="file-remove-button"
                  onClick={() => handleFileRemove(index)}
                >
                  &times;
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}
      <button type="submit" className="file-upload-button">
        Upload
      </button>
    </form>
  );
};
