import { User } from "@/types/common";
import { baseApi } from "./api";

export const users = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    completeInvite: builder.mutation({
      query: ({ token, userData }: { token: string; userData: User }) => ({
        url: "users/complete-invite",
        method: "POST",
        body: { token, ...userData },
      }),
    }),
    updateProfile: builder.mutation({
      query: (profileData) => ({
        url: "users/profile/update",
        method: "PUT",
        body: profileData,
      }),
    }),

    updatePassword: builder.mutation({
      query: (passwordData) => ({
        url: "users/profile/update-password",
        method: "PUT",
        body: passwordData,
      }),
    }),
    updateUserTechStacks: builder.mutation<void, { userId: string; techIds: string[] }>({
      query: ({ userId, techIds }) => ({
        url: '/users/update-techstacks',
        method: 'PUT',
        body: { userId, techIds },
      }),
    }),
  }),
});

export const {useUpdateUserTechStacksMutation, useUpdateProfileMutation, useUpdatePasswordMutation} = users;