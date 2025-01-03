import { type FC } from "react";
import {
    createBrowserRouter,
    RouterProvider,
    RouteObject,
} from "react-router-dom";
import { MainPage } from "page/MainPage";
import { Toolbar } from "widget/Toolbar/Toolbar";
import { Layout } from "app/Layout/Layout";

const routes: RouteObject = {
    path: "/",
    element: (
        <Layout
            toolbar={
                <Toolbar
                    items={[
                        {
                            label: "files",
                            path: "/files",
                        },
                        {
                            label: "similarity",
                            path: "/similarity",
                        },
                    ]}
                />
            }
        />
    ),
    errorElement: <div>Page 404</div>,
    children: [
        {
            path: "files",
            element: <MainPage />,
            errorElement: <div>Page 404</div>,
        },
        {
            path: "similarity",
            element: <MainPage />,
            errorElement: <div>Page 404</div>,
        },
    ],
};

export const AppRouterProvider: FC = () => {
    return <RouterProvider router={createBrowserRouter([routes])} />;
};
