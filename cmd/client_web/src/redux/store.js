import { createStore, combineReducers } from "redux";
import { composeWithDevTools } from "redux-devtools-extension";
import { searchReducer } from "./reducer";

const allReducer = combineReducers({
  search: searchReducer,
});

export default createStore(allReducer, composeWithDevTools());
