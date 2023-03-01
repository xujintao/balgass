import {
  SEARCH_LOADING,
  SEARCH_USERS,
  SEARCH_ERR,
  MODAL_JOIN,
  MODAL_LOGIN,
  MODAL_PASSWORD,
  MODAL_CLOSE,
} from "./const";

export const searchLoading = (data) => ({ type: SEARCH_LOADING, data });
export const searchUsers = (data) => ({ type: SEARCH_USERS, data });
export const searchErr = (data) => ({ type: SEARCH_ERR, data });

export const modalJoin = (data) => ({ type: MODAL_JOIN, data });
export const modalLogin = (data) => ({ type: MODAL_LOGIN, data });
export const modalPassword = (data) => ({ type: MODAL_PASSWORD, data });
export const modalClose = (data) => ({ type: MODAL_CLOSE, data });
