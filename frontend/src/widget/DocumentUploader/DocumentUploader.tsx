import { FC, useState } from "react";
import client from "shared/api";
import "./style.css";

type DocumentUploaderProps = {
    onDocumentUploaded?: (documentID: string) => void;
};

export const DocumentUploader: FC<DocumentUploaderProps> = ({
    onDocumentUploaded,
}) => {
    const [file, setFile] = useState<File | null>(null);

    if (onDocumentUploaded == undefined) {
        onDocumentUploaded = (_: string) => { }; // eslint-disable-line
    }

    const handleUpload = async () => {
        if (file) {
            const formData = new FormData();
            formData.append("document", file);

            try {
                const response = await client.document.uploadDocument(formData);
                onDocumentUploaded(response.documentID || "");
                alert("Document uploaded!");
            } catch (error) {
                alert("Failed to upload the document.");
            }
        }
    };

    return (
        <div className="upload-section">
            <input
                type="file"
                onChange={(e) => setFile(e.target.files ? e.target.files[0] : null)}
            />
            <button onClick={handleUpload} disabled={!file}>
                Upload Document
            </button>
        </div>
    );
};
