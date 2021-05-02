import { SEARCH_LOADING, SEARCH_USERS, SEARCH_ERR } from "./const";

export const searchLoading = (data) => ({ type: SEARCH_LOADING, data });
export const searchUsers = (data) => ({ type: SEARCH_USERS, data });
export const searchErr = (data) => ({ type: SEARCH_ERR, data });
