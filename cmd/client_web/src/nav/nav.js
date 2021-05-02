import React from "react";
import { Link, withRouter } from "react-router-dom";
import { connect } from "react-redux";
import { searchLoading, searchUsers, searchErr } from "../redux/action";
import axios from "axios";
import "./nav.css";
import logo from "./images/logo.png";
import "./fonts/iconfont.css";

class Nav extends React.Component {
  refInput = React.createRef();
  search = (event) => {
    event.preventDefault();
    this.props.searchLoading();
    const keyword = this.refInput.current.value;
    axios.get(`/api1/search/users?q=${keyword}`).then(
      (response) => this.props.searchUsers(response.data.items),
      (error) => this.props.searchErr(error.message)
    );
    this.props.history.push(`/search?q=${keyword}`);
  };
  render() {
    return (
      <div className="nav">
        <div className="nav-link">
          <ul>
            <li>
              <Link to="/">
                <img className="nav-link-logo" src={logo} alt="" />
                首页
              </Link>
            </li>
            <li>
              <Link to="/bugs">bug报告</Link>
            </li>
            <li>
              <Link to="/download">游戏下载</Link>
            </li>
          </ul>
        </div>
        <div className="nav-search">
          <form className="nav-search-form" action="/search">
            <input
              className="nav-search-text"
              type="text"
              placeholder="我的发明可以让整个小区停电"
              name="q"
              ref={this.refInput}
            />
            <button
              className="nav-search-btn iconfont icon-search"
              onClick={this.search}
            />
          </form>
        </div>
        <div className="nav-user">
          <div className="nav-user-login">
            <ul>
              <li>
                <a href="#">
                  <img
                    className="nav-user-avatar"
                    src="images/avatar.webp"
                    alt=""
                  />
                </a>
              </li>
              <li>
                <a href="#">消息</a>
              </li>
              <li>
                <a href="#">动态</a>
              </li>
              <li>
                <a href="#">收藏</a>
              </li>
            </ul>
          </div>
          <div className="nav-user-logout">
            <ul>
              <li>
                <a href="login.html">登录</a>
              </li>
              <li>
                <a href="register.html">注册</a>
              </li>
            </ul>
          </div>
        </div>
      </div>
    );
  }
}

export default withRouter(
  connect((state) => ({ search: state.search }), {
    searchLoading,
    searchUsers,
    searchErr,
  })(Nav)
);
