import { type FC } from "react";
import {
  createBrowserRouter,
  RouterProvider,
  RouteObject,
} from "react-router-dom";
import { Layout } from "../layout/Layout";

const routes: RouteObject = {
  element: <Layout />,
  errorElement: <div>Page 404</div>,
  children: [
    {
      element: (
        <div>
          <b>Coming soon!</b>
        </div>
      ),
    },
  ],
};

export const AppRouterProvider: FC = () => {
  return <RouterProvider router={createBrowserRouter([routes])} />;
};
