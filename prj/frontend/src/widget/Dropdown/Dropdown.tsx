import { FC, useState, useEffect, useRef } from "react";
import "./style.css";

type Option = { value: string; label: string };

type DropdownProps = {
  options?: Option[];
  fetchOptions?: () => Promise<Option[]>;
  selectedValues: string[];
  onChange: (selected: string[]) => void;
  placeholder?: string;
  maxSelections?: number;
  searchPlaceholder?: string;
};

export const Dropdown: FC<DropdownProps> = ({
  options = [],
  fetchOptions,
  selectedValues,
  onChange,
  placeholder = "Select...",
  searchPlaceholder = "Search...",
  maxSelections,
}) => {
  const dropdownRef = useRef<HTMLDivElement>(null);

  const [isOpen, setIsOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [hasLoadedOptions, setHasLoadedOptions] = useState(false);

  const [searchQuery, setSearchQuery] = useState("");
  const [filteredOptions, setFilteredOptions] = useState<Option[]>(options);

  const handleDropdownToggle = () => {
    setIsOpen((prev) => {
      const nextState = !prev;
      if (nextState && fetchOptions && !hasLoadedOptions) {
        loadOptions();
      }

      return nextState;
    });
  };

  const loadOptions = async () => {
    if (fetchOptions) {
      setIsLoading(true);

      try {
        options = await fetchOptions();
        setHasLoadedOptions(true);
      } catch (error) {
        console.error("Failed to fetch options:", error);
      } finally {
        setIsLoading(false);
      }
    }
  };

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const query = e.target.value;
    setSearchQuery(query);

    setFilteredOptions(
      options.filter((option) =>
        option.label.toLowerCase().includes(query.toLowerCase())
      )
    );
  };

  const handleSelect = (value: string) => {
    if (selectedValues.includes(value)) {
      onChange(selectedValues.filter((v) => v !== value));
    } else if (!maxSelections || selectedValues.length < maxSelections) {
      onChange([...selectedValues, value]);
    }
  };

  const handleOutsideClick = (e: MouseEvent) => {
    if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
      setIsOpen(false);
    }
  };

  useEffect(() => {
    document.addEventListener("mousedown", handleOutsideClick);
    return () => {
      document.removeEventListener("mousedown", handleOutsideClick);
    };
  }, []);

  return (
    <div className="dropdown" ref={dropdownRef}>
      <div
        className="dropdown-placeholder"
        onClick={handleDropdownToggle}
        role="button"
        tabIndex={0}
        aria-expanded={isOpen}
      >
        {selectedValues.length > 0
          ? selectedValues
            .map((val) => options.find((opt) => opt.value === val)?.label)
            .join(", ")
          : placeholder}
      </div>
      {isOpen && (
        <div className="dropdown-menu">
          <input
            type="text"
            placeholder={searchPlaceholder}
            value={searchQuery}
            onChange={handleSearchChange}
            className="dropdown-search"
            aria-label="Search"
          />
          {isLoading ? (
            <div className="dropdown-loading">Loading...</div>
          ) : (
            <ul className="dropdown-list" role="listbox">
              {filteredOptions.map((option) => (
                <li
                  key={option.value}
                  onClick={() => handleSelect(option.value)}
                  className={`dropdown-item ${selectedValues.includes(option.value) ? "selected" : ""
                    }`}
                  role="option"
                  aria-selected={selectedValues.includes(option.value)}
                >
                  {option.label}
                </li>
              ))}
              {filteredOptions.length === 0 && (
                <li className="dropdown-no-options">No options found</li>
              )}
            </ul>
          )}
        </div>
      )}
    </div>
  );
};
