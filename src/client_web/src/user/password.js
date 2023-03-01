import React from "react";
import { connect } from "react-redux";
import { modalLogin } from "../redux/action";
import {
  Breadcrumb,
  Steps,
  Form,
  Input,
  Button,
  Result,
  Modal,
  Row,
  Col,
} from "antd";
import "./password.css";

const { Step } = Steps;

class Password extends React.Component {
  state = {
    code: false,
    current: 0,
  };

  handleChangeStep = (current) => {
    if (current === 2) {
      current = 3;
    }
    // console.log(current);
    this.setState({ current });
  };

  handleFinishVerifyEmail = (value) => {
    console.log(value);
    this.setState({ code: true });
  };

  handleFinishFailedVerifyEmail = (err) => {
    console.log(err);
  };

  handleLogin = () => {
    this.props.modalLogin();
  };

  handleCancel = () => {
    this.setState({ code: false, current: 1 });
  };

  render() {
    const { current } = this.state;
    return (
      <>
        <div className="password-header">
          <Breadcrumb separator=">">
            <Breadcrumb.Item>用户</Breadcrumb.Item>
            <Breadcrumb.Item>忘记密码</Breadcrumb.Item>
          </Breadcrumb>
        </div>

        <div className="password-steps">
          <Steps current={current} onChange={this.handleChangeStep}>
            <Step title="确认账号" />
            <Step title="重置密码" />
            <Step title="重置成功" />
          </Steps>
        </div>

        <div className="password-form">
          {current === 0 ? (
            <Form
              size="large"
              onFinish={this.handleFinishVerifyEmail}
              onFinishFailed={this.handleFinishFailedVerifyEmail}
              // validateMessages={this.validateMessages}
            >
              <Form.Item
                name="email"
                rules={[
                  { type: "email", message: "无效的邮箱" },
                  { required: true, message: "邮箱不能为空" },
                ]}
              >
                <Input placeholder="请输入绑定的邮箱" />
              </Form.Item>
              <Form.Item noStyle>
                <Button type="primary" htmlType="submit" block>
                  确认
                </Button>
              </Form.Item>
            </Form>
          ) : current === 1 ? (
            <Form size="large" validateMessages={this.validateMessages}>
              <Form.Item
                name="password"
                hasFeedback
                rules={[{ required: true, message: "新密码不能为空" }]}
              >
                <Input.Password placeholder="新密码（6-16个字符组成，区分大小写）" />
              </Form.Item>
              <Form.Item
                name="confirm"
                dependencies={["password"]}
                hasFeedback
                rules={[
                  { required: true, message: "确认密码不能为空" },
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
                <Input.Password placeholder="确认密码" />
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
              <Form.Item noStyle>
                <Button type="primary" htmlType="submit" block>
                  确认
                </Button>
              </Form.Item>
            </Form>
          ) : current === 2 || current === 3 ? (
            <Result status="success" title="更改密码成功，请牢记新密码" />
          ) : null}
        </div>
        {current === 0 || current === 1 ? (
          <div className="password-login">
            <Button size="small" type="link" onClick={this.handleLogin}>
              已有账号，直接登录&gt;
            </Button>
          </div>
        ) : null}
        <Modal
          visible={this.state.code}
          onCancel={this.handleCancel}
          footer={null}
          maskClosable={false}
          destroyOnClose={true}
        >
          <div></div>
        </Modal>
      </>
    );
  }
}

export default connect(() => ({}), { modalLogin })(Password);
