import { CreateRoleRequest, UpdateRoleRequest } from "@/types";
import { baseApi } from "./api";

export const roles = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getRoles: builder.query<any, void>({
      query: () => "/roles",
      providesTags: ["roles"],
    }),
    getRole: builder.query<any, { id: string }>({
      query: ({ id }) => `/roles/${id}`,
      providesTags: ["roles"],
    }),
    createRole: builder.mutation<any, CreateRoleRequest>({
      query: ({ name, description, icon, color }) => ({
        url: "/roles",
        method: "POST",
        body: { name, description, icon, color },
      }),
      invalidatesTags: ["roles"],
    }),
    updateRole: builder.mutation<any, UpdateRoleRequest>({
      query: ({ id, name, description, icon, color }) => ({
        url: `/roles/${id}`,
        method: "PATCH",
        body: { name, description, icon, color },
      }),
      invalidatesTags: ["roles"],
    }),
    deleteRole: builder.mutation<any, { id: string }>({
      query: ({ id }) => ({
        url: `/roles/${id}`,
        method: "DELETE",
      }),
      invalidatesTags: ["roles"],
    }),
  }),
});

export const {
  useGetRolesQuery,
  useGetRoleQuery,
  useCreateRoleMutation,
  useUpdateRoleMutation,
  useDeleteRoleMutation,
} = roles;
