import { Application } from "@/types/common";
import { baseApi } from "./api";

export const applications = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAllApplications: builder.query<Application[], void>({
      query: () => ({
        url: "applications",
        method: "GET",
      }),
      providesTags: (result) =>
        result
          ? [
              ...result.map(({ id }) => ({ type: 'applications' as const, id })),
              { type: 'applications', id: 'LIST' },
            ]
          : [{ type: 'applications', id: 'LIST' }],
    }),
    loadBranches: builder.query({
      query: ({repoUrl}: {repoUrl: string}) => ({
        url: "applications/branches",
        method: "POST",
        body: { repoUrl },
      }),
    }),
    createApplication: builder.mutation({
      query: ({ appName, description }: { appName: string; description: string }) => ({
        url: "applications",
        method: "POST",
        body: { appName, description },
      }),
      invalidatesTags: [{ type: 'applications', id: 'LIST' }],
    }),
    deployApplication: builder.mutation({
      query: ({ appId, repoUrl, branch }: { appId: string; repoUrl: string; branch: string }) => ({
        url: `applications/${appId}/deploy`,
        method: "POST",
        body: { repoUrl, branch },
      }),
      invalidatesTags: (result, error, { appId }) => [{ type: 'applications', id: appId }],
    }),
    uploadApplication: builder.mutation({
      query: ({ appId, file }: { appId: string; file: File }) => {
        const formData = new FormData();
        formData.append('file', file);
        return {
          url: `applications/${appId}/upload`,
          method: "POST",
          body: formData,
          // Don't set Content-Type header, browser will set it with boundary
          formData: true,
        };
      },
      invalidatesTags: (result, error, { appId }) => [{ type: 'applications', id: appId }],
    }),
    deleteApplication: builder.mutation({
      query: ({ appId }: { appId: string }) => ({
        url: `applications/${appId}`,
        method: "DELETE",
      }),
      invalidatesTags: (result, error, { appId }) => [
        { type: 'applications', id: appId },
        { type: 'applications', id: 'LIST' }
      ],
    }),
    startApplication: builder.mutation({
      query: ({ appId }: { appId: string }) => ({
        url: `applications/${appId}/start`,
        method: "POST",
      }),
      invalidatesTags: (result, error, { appId }) => [{ type: 'applications', id: appId }],
    }),
    stopApplication: builder.mutation({
      query: ({ appId }: { appId: string }) => ({
        url: `applications/${appId}/stop`,
        method: "POST",
      }),
      invalidatesTags: (result, error, { appId }) => [{ type: 'applications', id: appId }],
    }),
  }),
});

export const {
  useGetAllApplicationsQuery,
  useLoadBranchesQuery,
  useCreateApplicationMutation,
  useDeployApplicationMutation,
  useUploadApplicationMutation,
  useStartApplicationMutation,
  useStopApplicationMutation,
  useDeleteApplicationMutation,
} = applications;