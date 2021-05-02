import { SEARCH_LOADING, SEARCH_USERS, SEARCH_ERR } from "./const";

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
