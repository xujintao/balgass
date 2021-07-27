import {
  SEARCH_LOADING,
  SEARCH_USERS,
  SEARCH_ERR,
  MODAL_JOIN,
  MODAL_LOGIN,
  MODAL_PASSWORD_VALIDATE_EMAIL,
  MODAL_PASSWORD_NEW_PASSWORD,
  MODAL_PASSWORD_RESULT,
  MODAL_CLOSE,
} from "./const";

export const searchLoading = (data) => ({ type: SEARCH_LOADING, data });
export const searchUsers = (data) => ({ type: SEARCH_USERS, data });
export const searchErr = (data) => ({ type: SEARCH_ERR, data });

export const modalJoin = (data) => ({ type: MODAL_JOIN, data });
export const modalLogin = (data) => ({ type: MODAL_LOGIN, data });
export const modalPasswordValidateEmail = (data) => ({
  type: MODAL_PASSWORD_VALIDATE_EMAIL,
  data,
});
export const modalPasswordNewPassword = (data) => ({
  type: MODAL_PASSWORD_NEW_PASSWORD,
  data,
});
export const modalPasswordResult = (data) => ({
  type: MODAL_PASSWORD_RESULT,
  data,
});
export const modalClose = (data) => ({
  type: MODAL_CLOSE,
  data,
});
