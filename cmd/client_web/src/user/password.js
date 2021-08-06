import React from "react";
import { connect } from "react-redux";
import { modalLogin } from "../redux/action";
import "./password.css";

class Password extends React.Component {
  state = {
    isVerifyEmail: true,
    isNewPassword: false,
    isPasswordResult: false,
    email: "",
    emailErr: "",
  };

  validateEmail = (event) => {
    const email = event.target.value;
    event.preventDefault();
    if (email == "") {
      this.setState({ emailErr: "请输入注册时用的邮箱" });
      return;
    }
    // validate
    this.setState({ email, emailErr: "" });
  };

  verifyEmail = (event) => {
    event.preventDefault();
    // this.props.modalPasswordNewPassword();
    this.setState({
      isVerifyEmail: false,
      isNewPassword: true,
      isPasswordResult: false,
    });
  };

  newPassword = (event) => {
    event.preventDefault();
    // this.props.modalPasswordResult();
    this.setState({
      isVerifyEmail: false,
      isNewPassword: false,
      isPasswordResult: true,
    });
  };

  modalLogin = (event) => {
    event.preventDefault();
    this.props.modalLogin();
  };

  render() {
    const { email, emailErr, isVerifyEmail, isNewPassword, isPasswordResult } =
      this.state;
    return (
      <div className="password">
        <h1 className="password-title">忘记密码</h1>
        <div className="password-steps">
          <ul>
            <li>
              <a
                className={`password-step-verify-email ${
                  isVerifyEmail ? "password-step-active" : ""
                }`}
              >
                确认账号
              </a>
            </li>
            <li>
              <a
                className={`password-step-new-password ${
                  isNewPassword ? "password-step-active" : ""
                }`}
              >
                重置密码
              </a>
            </li>
            <li>
              <a
                className={`password-step-result ${
                  isPasswordResult ? "password-step-active" : ""
                }`}
              >
                重置成功
              </a>
            </li>
          </ul>
        </div>
        {isVerifyEmail ? (
          <form className="password-form-verify-email">
            <ul>
              <li>
                <input
                  className={`password-form-input ${
                    emailErr ? "input-border-red" : ""
                  }`}
                  type="text"
                  placeholder="请输入绑定的邮箱"
                  onChange={this.validateEmail}
                  // onBlur={this.validateEmail}
                />
                {emailErr ? (
                  <p className="password-verify-email-message">注意邮箱格式</p>
                ) : null}
              </li>
              <li>
                <button
                  className={`password-form-submit ${
                    emailErr || !email ? "" : "btn-active"
                  }`}
                  disabled={emailErr ? true : false}
                  onClick={this.verifyEmail}
                >
                  确认
                </button>
              </li>
            </ul>
          </form>
        ) : isNewPassword ? (
          <form className="password-form-new-password" action="">
            <ul>
              <li>
                <input
                  className="password-form-input"
                  type="password"
                  placeholder="新密码（6-16个字符组成，区分大小写）"
                />
                <p className="password-error-message">密码不能少于6个字符</p>
                <span className="password-safe">安全</span>
              </li>
              <li>
                <input
                  className="password-form-input"
                  type="password"
                  placeholder="确认密码"
                />
              </li>
              <li>
                <input
                  className="password-form-input"
                  type="text"
                  placeholder="请输入邮件验证码"
                />
                <button className="password-form-btn">点击获取</button>
              </li>
              <li>
                <button
                  className="password-form-submit"
                  onClick={this.newPassword}
                >
                  确认修改
                </button>
              </li>
            </ul>
          </form>
        ) : isPasswordResult ? (
          <h4 className="password-form-result">更改密码成功，请牢记新密码</h4>
        ) : null}
        {isVerifyEmail || isNewPassword ? (
          <div className="password-login">
            <a onClick={this.modalLogin}>已有账号，直接登录&gt;</a>
          </div>
        ) : null}
      </div>
    );
  }
}

export default connect(() => ({}), {
  modalLogin,
})(Password);
