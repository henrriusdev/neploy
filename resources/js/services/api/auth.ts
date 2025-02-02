import { CompleteInviteRequest } from "@/types/request";
import { baseApi } from "./api";
import { User } from "@/types/common";

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    completeInvite: builder.mutation<void, CompleteInviteRequest>({
      query: (data) => ({
        url: "users/complete-invite",
        method: "POST",
        body: data,
      }),
    }),
    login: builder.mutation<User, { email: string; password: string }>({
      query: ({ email, password }) => ({
        url: "login",
        method: "POST",
        body: { email, password },
      }),
      transformErrorResponse: (response: any) => {
        if (response.status === 303) {
          return { data: null, meta: { location: response.headers.get('location') } };
        }
        return response;
      },
    }),
  }),
});

export const { useCompleteInviteMutation, useLoginMutation } = authApi;
