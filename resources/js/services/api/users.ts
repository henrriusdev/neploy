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
  }),
});