import { FC } from "react";
import "./style.css";

type SearchBarProps = {
  query: string;
  onQueryChange: (query: string) => void;
  onSearch: () => void;
};

export const SearchBar: FC<SearchBarProps> = ({
  query,
  onQueryChange,
  onSearch,
}) => {
  return (
    <div className="search-bar">
      <input
        type="text"
        value={query}
        onChange={(e) => onQueryChange(e.target.value)}
        placeholder="Search documents..."
      />
      <button onClick={onSearch}>Search</button>
    </div>
  );
};
