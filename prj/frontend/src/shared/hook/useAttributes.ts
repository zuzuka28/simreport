import client from "shared/api";
import { components } from "shared/api/document/api-types";
import { createDataHook } from "./createDataHook";

type AttributeResult = components["responses"]["AttributeResult"]["content"]["application/json"];

type AttributeQuery = components["schemas"]["AttributeRequest"];

export const useAttributes = createDataHook<AttributeQuery, AttributeResult>(
   (q: AttributeQuery) => { return client.attribute.fetchValues(q)} ,
)