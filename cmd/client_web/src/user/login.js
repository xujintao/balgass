import React from "react";
import { connect } from "react-redux";
import {
  modalJoin,
  modalPasswordValidateEmail,
  modalClose,
} from "../redux/action";
// import { Link } from "react-router-dom";
import "./login.css";

class Login extends React.Component {
  state = { email: "", emailErr: "", password: "", passwordErr: "" };

  checkEmail = (event) => {
    const email = event.target.value;
    if (email == "") {
      this.setState({ emailErr: "请输入注册时用的邮箱" });
      return;
    }
    // validate
    this.setState({ email, emailErr: "" });
  };

  checkPassword = (event) => {
    const password = event.target.value;
    if (password == "") {
      this.setState({ passwordErr: "喵，你没输入密码么？" });
      return;
    }
    this.setState({ password, passwordErr: "" });
  };

  resetPassword = (event) => {
    event.preventDefault();
    this.props.modalPasswordValidateEmail();
  };

  modalJoin = (event) => {
    event.preventDefault();
    this.props.modalJoin();
  };

  login = (event) => {
    event.preventDefault();
    // validate
    this.props.modalClose();
  };

  render() {
    const { emailErr, passwordErr } = this.state;
    return (
      <div className="login">
        <h1 className="login-title">登录</h1>
        <form className="login-form" action="">
          <ul>
            <li>
              <input
                className={`login-form-input ${
                  emailErr ? "input-border-red" : null
                }`}
                type="text"
                placeholder="邮箱"
                onChange={this.checkEmail}
              />
              {emailErr ? <p className="error-message">{emailErr}</p> : null}
            </li>
            <li>
              <input
                className={`login-form-input ${
                  passwordErr ? "input-border-red" : null
                }`}
                type="password"
                placeholder="社区密码（6-16个字符组成，区分大小写）"
                onChange={this.checkPassword}
              />
              {passwordErr ? (
                <p className="error-message">{passwordErr}</p>
              ) : null}
            </li>
            <li className="login-form-remember">
              <div className="left">
                <input type="checkbox" />
                记住我
                <span>(不是自己的电脑不要勾选此项)</span>
              </div>
              <div className="right">
                {/* <Link to="/password">忘记密码？</Link> */}
                <a onClick={this.resetPassword}>忘记密码？</a>
              </div>
            </li>
            <li className="login-form-submit">
              <button className="login-form-submit-btn" onClick={this.login}>
                登录
              </button>
              {/* <Link className="login-form-join" to="/join">没有账号，这里注册&gt;</Link> */}
              <a className="login-form-join" onClick={this.modalJoin}>
                没有账号，这里注册&gt;
              </a>
            </li>
          </ul>
        </form>
      </div>
    );
  }
}

export default connect(() => ({}), {
  modalJoin,
  modalPasswordValidateEmail,
  modalClose,
})(Login);
