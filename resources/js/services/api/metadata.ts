import { MetadataRequest } from "@/types/request";
import { baseApi } from "./api";

export const metadata = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    updateMetadata: builder.mutation({
      query: ({ data }: { data: MetadataRequest }) => ({
        url: "metadata",
        method: "PATCH",
        body: data,
      }),
    }),
  }),
});

export const { useUpdateMetadataMutation } = metadata;
