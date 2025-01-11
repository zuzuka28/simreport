import { FC, ReactNode, useEffect, useRef } from "react";
import "./style.css";

type PopupProps = {
  onClose: () => void;
  children: ReactNode[];
};

export const Popup: FC<PopupProps> = ({ onClose, children }) => {
  const popupRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (popupRef.current && !popupRef.current.contains(e.target as Node)) {
        onClose();
      }
    };

    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [onClose]);

  return (
    <div className="popup" ref={popupRef}>
      <button className="popup-close" onClick={onClose}>
        &times;
      </button>
      <div className="popup-content">{children}</div>
    </div>
  );
};
