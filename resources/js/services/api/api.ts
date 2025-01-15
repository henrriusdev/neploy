import { RootState } from "@/store";
import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";

const baseQuery = fetchBaseQuery({
  baseUrl: "/",
});

export const baseApi = createApi({
  reducerPath: "base-api",
  tagTypes: [
    "onboarding",
    "auth",
    "accounts",
    "transactions",
    "overview",
    "categories",
    "budget",
    "budgetStep",
    "enforcement",
  ],
  baseQuery: baseQuery,
  endpoints: () => ({}),
});
