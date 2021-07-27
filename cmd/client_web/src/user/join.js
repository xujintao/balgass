import React from "react";
import { connect } from "react-redux";
import { modalLogin, modalClose } from "../redux/action";
// import { Link } from "react-router-dom";
import "./join.css";

class Join extends React.Component {
  modalLogin = (event) => {
    event.preventDefault();
    this.props.modalLogin();
  };

  render() {
    return (
      <div className="join">
        <h1 className="join-title">注册</h1>
        <form className="join-form" action="">
          <ul>
            <li>
              <input
                className="join-form-input"
                type="text"
                placeholder="社区昵称"
              />
              <p className="error-message">请告诉我你的昵称吧</p>
            </li>
            <li>
              <input
                className="join-form-input"
                type="password"
                placeholder="社区密码（6-16个字符组成，区分大小写）"
              />
              <p className="error-message">密码不能少于6个字符</p>
              <span className="safe">安全</span>
            </li>
            <li>
              <input
                className="join-form-input"
                type="text"
                placeholder="邮箱（用于验证）"
              />
            </li>
            <li>
              <input
                className="join-form-input"
                type="text"
                placeholder="请输入验证码"
              />
              <button className="join-form-btn">点击获取</button>
            </li>
            <li className="join-form-agree">
              <input type="checkbox" />
              我已经同意有关条款
            </li>
            <li>
              <button className="join-form-submit">注册</button>
            </li>
            <li className="join-form-login">
              {/* <Link to="/login">已有账号，直接登录&gt;</Link> */}
              <a onClick={this.modalLogin}>已有账号，直接登录&gt;</a>
            </li>
          </ul>
        </form>
      </div>
    );
  }
}

export default connect(() => ({}), { modalLogin, modalClose })(Join);
