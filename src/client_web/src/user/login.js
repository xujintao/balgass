import React from "react";
import { connect } from "react-redux";
import { modalPassword, modalJoin } from "../redux/action";
import { Breadcrumb, Form, Input, Checkbox, Button } from "antd";
import { MailOutlined, LockOutlined } from "@ant-design/icons";
import "./login.css";

function Login(props) {
  const handlePassword = () => {
    props.modalPassword();
  };

  const onFinish = (values) => {
    console.log(values);
  };

  const handleJoin = () => {
    props.modalJoin();
  };

  return (
    <>
      <div className="login-header">
        <Breadcrumb separator=">">
          <Breadcrumb.Item>用户</Breadcrumb.Item>
          <Breadcrumb.Item>登录</Breadcrumb.Item>
        </Breadcrumb>
      </div>
      <div className="login-form">
        <Form
          size="large"
          initialValues={{ remember: true }}
          onFinish={onFinish}
        >
          <Form.Item
            name="email"
            rules={[
              { type: "email", message: "无效的邮箱" },
              { required: true, message: "邮箱不能为空" },
            ]}
          >
            <Input prefix={<MailOutlined />} placeholder="邮箱" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: "密码不能为空" }]}
          >
            <Input
              prefix={<LockOutlined />}
              type="password"
              placeholder="社区密码"
            />
          </Form.Item>
          <Form.Item>
            <Form.Item name="remember" valuePropName="checked" noStyle>
              <Checkbox>记住我(不是自己的电脑不要勾选此项)</Checkbox>
            </Form.Item>
            <Button
              size="small"
              type="link"
              className="login-form-password"
              onClick={handlePassword}
            >
              忘记密码？
            </Button>
          </Form.Item>

          <Form.Item noStyle>
            <Button type="primary" htmlType="submit" block>
              登录
            </Button>
          </Form.Item>
        </Form>
      </div>
      <div className="login-join">
        <Button size="small" type="link" onClick={handleJoin}>
          没有账号，这里注册&gt;
        </Button>
      </div>
    </>
  );
}

export default connect(() => ({}), { modalPassword, modalJoin })(Login);
