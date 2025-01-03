import React from "react";
import { Link } from "react-router-dom";
import "./style.css";

interface MenuItem {
    label: string;
    path?: string;
    onClick?: () => void;
    icon?: string;
}

interface ToolbarProps {
    items: MenuItem[];
}

export const Toolbar: React.FC<ToolbarProps> = ({ items }) => {
    return (
        <nav className="toolbar">
            <ul className="toolbar-menu">
                {items.map((item, index) => (
                    <li key={index} className="toolbar-item">
                        {item.path ? (
                            <Link to={item.path} className="toolbar-link">
                                {item.icon && (
                                    <img
                                        src={item.icon}
                                        alt={item.label}
                                        className="toolbar-icon"
                                    />
                                )}
                                <span className="toolbar-label">{item.label}</span>
                            </Link>
                        ) : (
                            <button
                                className="toolbar-link"
                                onClick={item.onClick}
                                style={{
                                    background: "none",
                                    border: "none",
                                    cursor: "pointer",
                                    padding: 0,
                                }}
                            >
                                {item.icon && (
                                    <img
                                        src={item.icon}
                                        alt={item.label}
                                        className="toolbar-icon"
                                    />
                                )}
                                <span className="toolbar-label">{item.label}</span>
                            </button>
                        )}
                    </li>
                ))}
            </ul>
        </nav>
    );
};
