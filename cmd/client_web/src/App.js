import React from "react";
import { Route, BrowserRouter, Switch, Redirect } from "react-router-dom";
import { connect } from "react-redux";
import { modalClose } from "./redux/action";
import { Layout, Modal } from "antd";
import Nav from "./nav/nav";
import Home from "./home/home";
import Search from "./search/search";
import Join from "./user/join";
import Login from "./user/login";
import Password from "./user/password";
import "./App.less";

const { Header, Content, Footer } = Layout;

class App extends React.Component {
  handleCancel = () => {
    this.props.modalClose();
  };

  render() {
    const { isModalJoin, isModalLogin, isModalPassword } = this.props.modal;
    return (
      <BrowserRouter>
        <Layout>
          <Header className="header">
            <Nav />
          </Header>
          <Content className="main">
            <Switch>
              {/* https://stackoverflow.com/questions/49162311/react-difference-between-route-exact-path-and-route-path */}
              <Route exact path="/" component={Home} />
              {/* https://stackoverflow.com/questions/50667609/react-router-component-not-updating-on-url-search-param-change */}
              <Route path="/search" component={Search} />
              {/* <Route path="/bugs" component={Bugs}/> */}
              {/* <Route path="/download" component={Download}/> */}
              {/* <Route path="/join" component={Join} /> */}
              {/* <Route path="/login" component={Login} /> */}
              {/* <Route path="/password" component={Password} /> */}
              <Redirect to="/" />
            </Switch>
          </Content>
          <Footer></Footer>
        </Layout>
        <Modal
          visible={isModalJoin || isModalLogin || isModalPassword}
          onCancel={this.handleCancel}
          footer={null}
          maskClosable={false}
          destroyOnClose={true}
        >
          {isModalJoin ? (
            <Join />
          ) : isModalLogin ? (
            <Login />
          ) : isModalPassword ? (
            <Password />
          ) : null}
        </Modal>
      </BrowserRouter>
    );
  }
}

export default connect(
  (state) => ({
    modal: state.modal,
  }),
  { modalClose }
)(App);
