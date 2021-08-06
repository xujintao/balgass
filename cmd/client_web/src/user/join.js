import React from "react";
import { connect } from "react-redux";
import { modalLogin } from "../redux/action";
import { Breadcrumb, Form, Input, Checkbox, Button, Row, Col } from "antd";
import "./join.css";

function Join(props) {
  // const tailFormItemLayout = {
  //   wrapperCol: {
  //     xs: {
  //       span: 24,
  //       offset: 0,
  //     },
  //     sm: {
  //       span: 16,
  //       // offset: 8,
  //     },
  //   },
  // };

  const handleFinish = (values) => {
    console.log(values);
  };

  const handleLogin = () => {
    props.modalLogin();
  };

  return (
    <>
      <div className="join-header">
        <Breadcrumb separator=">">
          <Breadcrumb.Item>用户</Breadcrumb.Item>
          <Breadcrumb.Item>注册</Breadcrumb.Item>
        </Breadcrumb>
      </div>
      <div className="join-form">
        <Form size="large" onFinish={handleFinish}>
          <Form.Item
            name="nickname"
            rules={[
              { required: true, whitespace: true, message: "昵称不能为空" },
            ]}
          >
            <Input placeholder="社区昵称" />
          </Form.Item>

          <Form.Item
            name="password"
            hasFeedback
            rules={[{ required: true, message: "密码不能为空" }]}
          >
            {/* <Input type="password" placeholder="社区密码" /> */}
            <Input.Password placeholder="社区密码（6-16个字符组成，区分大小写）" />
          </Form.Item>

          <Form.Item
            name="confirm"
            dependencies={["password"]}
            hasFeedback
            rules={[
              { required: true, message: "密码不能为空" },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue("password") === value) {
                    return Promise.resolve();
                  }
                  return Promise.reject(
                    new Error(
                      "The two passwords that you entered do not match!"
                    )
                  );
                },
              }),
            ]}
          >
            <Input.Password placeholder="确认社区密码" />
          </Form.Item>

          <Form.Item
            name="email"
            rules={[
              { type: "email", message: "无效的邮箱" },
              { required: true, message: "邮箱不能为空" },
            ]}
          >
            <Input placeholder="邮箱（用于验证）" />
          </Form.Item>

          <Form.Item>
            <Row gutter={8}>
              <Col span={12}>
                <Form.Item
                  name="captcha"
                  noStyle
                  rules={[{ required: true, message: "验证码不能为空" }]}
                >
                  <Input placeholder="验证码" />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Button>获取验证码</Button>
              </Col>
            </Row>
          </Form.Item>

          <Form.Item
            name="agreement"
            valuePropName="checked"
            // {...tailFormItemLayout}
            rules={[
              {
                validator: (_, value) =>
                  value
                    ? Promise.resolve()
                    : Promise.reject(new Error("Should accept agreement")),
              },
            ]}
          >
            <Checkbox>
              我已经同意有关条款
              <Button
                size="small"
                type="link"
                target="_blank"
                href="https://www.bilibili.com/protocal/licence.html"
              >
                《服务协议》
              </Button>
              和
              <Button
                size="small"
                type="link"
                target="_blank"
                href="https://www.bilibili.com/blackboard/privacy-pc.html"
              >
                《隐私政策》
              </Button>
            </Checkbox>
          </Form.Item>

          <Form.Item noStyle>
            <Button type="primary" htmlType="submit" block>
              注册
            </Button>
          </Form.Item>
        </Form>
      </div>
      <div className="join-login">
        <Button size="small" type="link" onClick={handleLogin}>
          已有账号，直接登录&gt;
        </Button>
      </div>
    </>
  );
}

export default connect(() => ({}), { modalLogin })(Join);
