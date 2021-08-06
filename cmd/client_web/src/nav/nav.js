import React from "react";
import { Link, withRouter } from "react-router-dom";
import { connect } from "react-redux";
import {
  searchLoading,
  searchUsers,
  searchErr,
  modalLogin,
} from "../redux/action";
import axios from "axios";
import qs from "querystring";
import { Input, Space, Button } from "antd";
import "./nav.css";
import logo from "./images/logo.png";

class Nav extends React.Component {
  state = { keyword: "" };

  componentDidMount() {
    const { pathname, search } = this.props.location;
    if (pathname === "/search") {
      this.props.searchLoading();
      if (search.match(/\?q=/)) {
        const { q: keyword } = qs.parse(search.slice(1));
        this.setState({ keyword });
        axios.get(`/api1/search/users?q=${keyword}`).then(
          (response) => this.props.searchUsers(response.data.items),
          (error) => this.props.searchErr(error.message)
        );
      } else {
        this.props.searchUsers([]);
      }
    }
  }

  handleSearch = (value) => {
    this.props.searchLoading();
    const keyword = value;
    axios.get(`/api1/search/users?q=${keyword}`).then(
      (response) => this.props.searchUsers(response.data.items),
      (error) => this.props.searchErr(error.message)
    );
    this.props.history.push(`/search?q=${keyword}`);
  };

  handleChange = (event) => {
    this.setState({ keyword: event.target.value });
  };

  handleClick = () => {
    this.props.modalLogin();
  };

  render() {
    // console.log("@", this.state.keyword);
    return (
      <div className="nav">
        <div className="nav-link">
          <Space size={12} align="center">
            <Link to="/">
              <img className="nav-link-logo" src={logo} alt="" />
              首页
            </Link>
            <Link to="/issues">bug报告</Link>
            <Link to="/download">游戏下载</Link>
          </Space>
        </div>
        <div className="nav-search">
          <Input.Search
            placeholder="我的发明可以让整个小区停电"
            onSearch={this.handleSearch}
            enterButton
            value={this.state.keyword}
            onChange={this.handleChange}
          />
        </div>
        <div className="nav-user">
          <Space size={12}>
            <Button onClick={this.handleClick}>登录</Button>
          </Space>
        </div>
      </div>
    );
  }
}

export default withRouter(
  connect(
    (state) => ({
      login: state.login,
    }),
    {
      searchLoading,
      searchUsers,
      searchErr,
      modalLogin,
    }
  )(Nav)
);
