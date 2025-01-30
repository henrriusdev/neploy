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
  }),
});

export const { useCompleteInviteMutation } = authApi;
