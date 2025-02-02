import { MetadataRequest } from "@/types/request";
import { baseApi } from "./api";

export const metadata = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getMetadata: builder.query<MetadataRequest, void>({
      query: () => ({
        url: "metadata",
        method: "GET",
      }),
      providesTags: ["metadata"],
    }),
    updateMetadata: builder.mutation({
      query: ({ data }: { data: MetadataRequest }) => ({
        url: "metadata",
        method: "PATCH",
        body: data,
      }),
      invalidatesTags: ["metadata"],
    }),
  }),
});

export const { useUpdateMetadataMutation } = metadata;
