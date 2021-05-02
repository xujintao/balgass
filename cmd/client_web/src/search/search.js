import React from "react";
import qs from "querystring";
import axios from "axios";
import "./search.css";

class Search extends React.Component {
  state = {
    isFirst: true,
    isLoading: false,
    users: [],
    err: "",
  };

  componentDidMount() {
    this.setState({
      isFirst: false,
      isLoading: true,
    });
    const { search } = this.props.location;
    const { q: keyword } = qs.parse(search.slice(1));
    axios.get(`/api1/search/users?q=${keyword}`).then(
      (response) => {
        this.setState({
          isLoading: false,
          users: response.data.items,
        });
      },
      (error) => {
        this.setState({
          isLoading: false,
          err: error.message,
        });
      }
    );
  }

  render() {
    const { users, isFirst, isLoading, err } = this.state;
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

export default Search;
