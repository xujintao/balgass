import {
  SEARCH_LOADING,
  SEARCH_USERS,
  SEARCH_ERR,
  MODAL_JOIN,
  MODAL_LOGIN,
  MODAL_PASSWORD,
  MODAL_CLOSE,
} from "./const";

const initState = { isFirst: true, isLoading: false, users: [], err: "" };
export function searchReducer(preState = initState, action) {
  const { type, data } = action;
  switch (type) {
    case SEARCH_LOADING:
      return { isFirst: false, isLoading: true, users: [], err: "" };
    case SEARCH_USERS:
      return { isFirst: false, isLoading: false, users: data, err: "" };
    case SEARCH_ERR:
      return { isFirst: false, isLoading: false, users: [], err: data };
    default:
      return preState;
  }
}

const initModal = {
  isModalJoin: false,
  isModalLogin: false,
  isModalPassword: false,
};
export function modalReducer(preState = initModal, action) {
  const { type } = action;
  switch (type) {
    case MODAL_JOIN:
      return {
        isModalJoin: true,
        isModalLogin: false,
        isModalPassword: false,
      };
    case MODAL_LOGIN:
      return {
        isModalJoin: false,
        isModalLogin: true,
        isModalPassword: false,
      };
    case MODAL_PASSWORD:
      return {
        isModalJoin: false,
        isModalLogin: false,
        isModalPassword: true,
      };
    case MODAL_CLOSE:
      return initModal;
    default:
      return preState;
  }
}

const initLogin = { isLogin: false };
export function loginReducer(preState = initLogin, action) {
  const { type } = action;
  switch (type) {
    default:
      return preState;
  }
}
