import React from "react";
import { connect } from "react-redux";
import {
  modalPasswordValidateEmail,
  modalPasswordNewPassword,
  modalPasswordResult,
  modalClose,
} from "../redux/action";
import "./password.css";

class Password extends React.Component {
  render() {
    const { isModalValidateEmail, isModalNewPassword, isModalPasswordResult } =
      this.props.modal;
    return (
      <div className="password">
        <h1 className="title">忘记密码</h1>
        <div className="steps">
          <ul>
            <li>
              <a className={`a-verify ${isModalValidateEmail ? "active" : ""}`}>
                确认账号
              </a>
            </li>
            <li>
              <a className={`a-set ${isModalNewPassword ? "active" : ""}`}>
                重置密码
              </a>
            </li>
            <li>
              <a
                className={`a-result ${isModalPasswordResult ? "active" : ""}`}
              >
                重置成功
              </a>
            </li>
          </ul>
        </div>
        {isModalValidateEmail ? (
          <form className="form-verify">
            <input
              className="form-input"
              type="text"
              placeholder="请输入绑定的邮箱"
            />
            <p className="verify-message"></p>
            <button className="form-submit">确认</button>
          </form>
        ) : isModalNewPassword ? (
          <form className="form-set" action="">
            <ul>
              <li>
                <input
                  className="form-input"
                  type="password"
                  placeholder="新密码（6-16个字符组成，区分大小写）"
                />
                <p className="error-message">密码不能少于6个字符</p>
                <span className="safe">安全</span>
              </li>
              <li>
                <input
                  className="form-input"
                  type="password"
                  placeholder="确认密码"
                />
              </li>
              <li>
                <input
                  className="form-input"
                  type="text"
                  placeholder="请输入邮件验证码"
                />
                <button className="form-btn">点击获取</button>
              </li>
              <li>
                <button className="form-submit">确认修改</button>
              </li>
            </ul>
          </form>
        ) : isModalPasswordResult ? (
          <h4 className="form-result">更改密码成功，请牢记新密码</h4>
        ) : null}
      </div>
    );
  }
}

export default connect((state) => ({ modal: state.modal }), {
  modalPasswordValidateEmail,
  modalPasswordNewPassword,
  modalPasswordResult,
  modalClose,
})(Password);
