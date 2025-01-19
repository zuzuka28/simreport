import { FC, useEffect, useState } from "react";
import "./style.css"

type TaggerProps = {
    initialValues: string[],
    suggestions?: string[],
    fetchSuggestions?: (text: string) => Promise<string[]>,
    onBubblesChange: (bubbles: string[]) => void,
}

export const Tagger: FC<TaggerProps> = ({ initialValues = [], onBubblesChange, suggestions = [], fetchSuggestions }) => {
    const [inputValue, setInputValue] = useState('');
    const [items, setItems] = useState(initialValues);
    const [filteredSuggestions, setFilteredSuggestions] = useState<string[]>([]);

    useEffect(() => {
        const loadSuggestion = async () => {
            if (fetchSuggestions) {
                await fetchSuggestions(inputValue).then((result) => setFilteredSuggestions(result));
            } else {
                setFilteredSuggestions(
                    suggestions.filter((suggestion) =>
                        suggestion.toLowerCase().includes(inputValue.toLowerCase())
                    )
                );
            }
        }

        loadSuggestion()

    }, [inputValue]);

    const handleInputChange = (event) => {
        setInputValue(event.target.value);
    };

    const handleKeyDown = (event) => {
        if (event.key === 'Enter' && inputValue.trim() !== '') {
            const updatedItems = [...items, inputValue];
            setItems(updatedItems);
            setInputValue('');
            if (onBubblesChange) {
                onBubblesChange(updatedItems);
            }
        } else if (event.key === 'Tab' && filteredSuggestions.length > 0) {
            event.preventDefault();
            const updatedItems = [...items, filteredSuggestions[0]];
            setItems(updatedItems);
            setInputValue('');
            if (onBubblesChange) {
                onBubblesChange(updatedItems);
            }
        }
    };

    const handleBubbleClick = (index: number) => {
        const updatedItems = items.filter((_, i) => i !== index);
        setItems(updatedItems);
        if (onBubblesChange) {
            onBubblesChange(updatedItems);
        }
    };

    return (
        <div className="tagger_container">
            <input
                type="text"
                value={inputValue}
                onChange={handleInputChange}
                onKeyDown={handleKeyDown}
                placeholder="Enter text and press Enter"
                className="input"
            />
            {filteredSuggestions.length > 0 && (
                <ul className="tagger_suggestion-list">
                    {filteredSuggestions.map((suggestion, index) => (
                        <li key={index} className="suggestion-item">
                            {suggestion}
                        </li>
                    ))}
                </ul>
            )}
            <div className="tagger_bubble-container">
                {items.map((item, index) => (
                    <div
                        key={index}
                        className="bubble"
                        onClick={() => handleBubbleClick(index)}
                    >
                        {item}
                    </div>
                ))}
            </div>
        </div>
    );
};