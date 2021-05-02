import React from "react";
import { connect } from "react-redux";
import "./search.css";

class Search extends React.Component {
  render() {
    console.log(this.props);
    const { users, isFirst, isLoading, err } = this.props.search;
    return (
      <div className="w">
        <ul className="search-result">
          {isFirst ? (
            <h2>欢迎使用，输入关键字，随后点击搜索</h2>
          ) : isLoading ? (
            <h2>Loading...</h2>
          ) : err ? (
            <h2>{err}</h2>
          ) : users ? (
            users.map((u) => {
              return (
                <li key={u.id} className="card">
                  <a href={u.html_url} target="_blank" rel="noreferrer">
                    <img src={u.avatar_url} style={{ width: "100px" }} alt="" />
                  </a>
                  <p className="card-text">{u.login}</p>
                </li>
              );
            })
          ) : (
            {}
          )}
        </ul>
      </div>
    );
  }
}

export default connect((state) => ({ search: state.search }))(Search);
