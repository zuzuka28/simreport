import { type FC } from "react";
import {
  createBrowserRouter,
  RouterProvider,
  RouteObject,
} from "react-router-dom";
import { Toolbar } from "widget/Toolbar/Toolbar";
import { Layout } from "app/Layout/Layout";
import { DocumentManager } from "page/DocumentManager";
import { SimilarityCheck } from "page/SimilarityCheck/SimilarityCheck";
import { DocumentUpload } from "page/DocumentUpload/DocumentUpload";
import { DocumentBulkUpload } from "page/DocumentBulkUpload";

const routes: RouteObject = {
  path: "/",
  element: (
    <Layout
      toolbar={
        <Toolbar
          items={[
            {
              label: "upload",
              path: "/upload",
            },
            {
              label: "bulk_upload",
              path: "/bulk_upload",
            },
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
      path: "upload",
      element: <DocumentUpload />,
      errorElement: <div>Page 404</div>,
    },
    {
      path: "bulk_upload",
      element: <DocumentBulkUpload />,
      errorElement: <div>Page 404</div>,
    },
    {
      path: "files",
      element: <DocumentManager />,
      errorElement: <div>Page 404</div>,
    },
    {
      path: "similarity",
      element: <SimilarityCheck />,
      errorElement: <div>Page 404</div>,
    },
  ],
};

export const AppRouterProvider: FC = () => {
  return <RouterProvider router={createBrowserRouter([routes])} />;
};
