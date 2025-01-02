import { type FC } from "react";
import {
    createBrowserRouter,
    RouterProvider,
    RouteObject,
} from "react-router-dom";
import { MainPage } from "page/MainPage";

const routes: RouteObject = {
    path: "/",
    element: <MainPage />,
    errorElement: <div>Page 404</div>,
};

export const AppRouterProvider: FC = () => {
    return <RouterProvider router={createBrowserRouter([routes])} />;
};
