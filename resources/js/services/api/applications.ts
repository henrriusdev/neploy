import { baseApi } from "./api";

export const applications = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getAll: builder.query({
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
    create: builder.mutation({
      query: ({ appName, description }: { appName: string; description: string }) => ({
        url: "applications",
        method: "POST",
        body: { appName, description },
      })
    }),
    deploy: builder.mutation({
      query: ({ appId, repoUrl, branch }: { appId: string; repoUrl: string; branch: string }) => ({
        url: `applications/${appId}/deploy`,
        method: "POST",
        body: { repoUrl, branch },
      })
    }),
    upload: builder.mutation({
      query: ({ appId, file }: { appId: string; file: File }) => ({
        url: `applications/${appId}/upload`,
        method: "POST",
        body: file,
        headers: {
          "Content-Type": "multipart/form-data",
        },
      })
    }),

    start: builder.mutation({
      query: ({ appId }: { appId: string }) => ({
        url: `applications/${appId}/start`,
        method: "POST",
      })
    }),
    stop: builder.mutation({
      query: ({ appId }: { appId: string }) => ({
        url: `applications/${appId}/stop`,
        method: "POST",
      })
    }),
    delete: builder.mutation({
      query: ({ appId }: { appId: string }) => ({
        url: `applications/${appId}`,
        method: "DELETE",
      })
    }),
  }),
});