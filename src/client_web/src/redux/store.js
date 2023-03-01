import { createStore, combineReducers } from "redux";
import { composeWithDevTools } from "redux-devtools-extension";
import { searchReducer, loginReducer, modalReducer } from "./reducer";

const allReducer = combineReducers({
  search: searchReducer,
  modal: modalReducer,
  login: loginReducer,
});

export default createStore(allReducer, composeWithDevTools());
