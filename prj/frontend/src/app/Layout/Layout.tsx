import { FC, ReactElement } from "react";
import { Outlet } from "react-router-dom";
import "./style.css";

type LayoutProps = {
    toolbar?: ReactElement;
};

export const Layout: FC<LayoutProps> = ({ toolbar }) => {
    return (
        <>
            <div className="container">
                {toolbar}

                <div className="content_container">
                    <Outlet />
                </div>
            </div>
        </>
    );
};
