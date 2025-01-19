import { FC, useState } from "react";
import { SearchBar } from "widget/SearchBar";
import { DocumentList } from "widget/DocumentList";
import "./style.css";
import { useDocuments } from "shared/hook";
import client from "shared/api";
import { Dropdown } from "widget/Dropdown";

export const DocumentManager: FC = () => {
  const [query, setQuery] = useState<{ groupID: string[], name: string }>({
    groupID: [],
    name: "",
  });

  const { data, loading, error } = useDocuments(query);

  const handleSearch = () => { };

  const provideGroupIDs = async () => {
    const { items } = await client.attribute.fetchValues({
      attribute: "groupID",
    });
    const groups = items || [];
    return groups.map((doc: any) => ({ value: doc.id, label: doc.title }));
  };

  return (
    <>
      <div className="document_manager_container">
        <div className="action-section">
          <SearchBar
            query={query.name}
            onQueryChange={(v) => { setQuery({ ...query, name: v }) }}
            onSearch={handleSearch}
          />
          <div className="document_manager_dropdown">
            <Dropdown
              fetchOptions={provideGroupIDs}
              selectedValues={query.groupID}
              onChange={(v) => { setQuery({ ...query, groupID: v }) }}
              placeholder="Groups"
            />
          </div>
        </div>
        <DocumentList documents={data?.documents || []} loading={loading} error={error} />
      </div>
    </>
  );
};
