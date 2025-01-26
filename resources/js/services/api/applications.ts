import { Application } from "@/types/common";
import { baseApi } from "./api";

export const applications = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAllApplications: builder.query<Application[], void>({
      query: () => ({
        url: "applications",
        method: "GET",
      }),
    }),
    loadBranches: builder.query({
      query: ({repoUrl}: {repoUrl: string}) => ({
        url: "applications/branches",
        method: "GET",
        body: {repoUrl},
      }),
    }),
    createApplication: builder.mutation({
      query: ({ appName, description }: { appName: string; description: string }) => ({
        url: "applications",
        method: "POST",
        body: { appName, description },
      })
    }),
    deployApplication: builder.mutation({
      query: ({ appId, repoUrl, branch }: { appId: string; repoUrl: string; branch: string }) => ({
        url: `applications/${appId}/deploy`,
        method: "POST",
        body: { repoUrl, branch },
      })
    }),
    uploadApplication: builder.mutation({
      query: ({ appId, file }: { appId: string; file: File }) => ({
        url: `applications/${appId}/upload`,
        method: "POST",
        body: file,
        headers: {
          "Content-Type": "multipart/form-data",
        },
      })
    }),
    deleteApplication: builder.mutation({
      query: ({ appId }: { appId: string }) => ({
        url: `applications/${appId}`,
        method: "DELETE",
      })
    }),
    startApplication: builder.mutation({
      query: ({ appId }: { appId: string }) => ({
        url: `applications/${appId}/start`,
        method: "POST",
      })
    }),
    stopApplication: builder.mutation({
      query: ({ appId }: { appId: string }) => ({
        url: `applications/${appId}/stop`,
        method: "POST",
      })
    }),
  }),
});

export const { useGetAllApplicationsQuery, useLoadBranchesQuery, useCreateApplicationMutation, useDeployApplicationMutation, useUploadApplicationMutation, useStartApplicationMutation, useStopApplicationMutation, useDeleteApplicationMutation } = applications;