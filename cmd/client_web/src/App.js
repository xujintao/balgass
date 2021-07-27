import React from "react";
import { Route, BrowserRouter, Switch, Redirect } from "react-router-dom";
import Nav from "./nav/nav";
import Home from "./home/home";
import Search from "./search/search";
import Footer from "./footer/footer";
import { connect } from "react-redux";
import { modalClose } from "./redux/action";
import Join from "./user/join";
import Login from "./user/login";
import Password from "./user/password";
import "./App.css";
import "./fonts/iconfont.css";

class App extends React.Component {
  closeAuth = (event) => {
    event.preventDefault();
    this.props.modalClose();
  };

  render() {
    const { isModalJoin, isModalLogin, isModalPassword } = this.props.modal;
    return (
      <div>
        <div className="view-box">
          <BrowserRouter>
            <Nav />
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
            <Footer />
          </BrowserRouter>
        </div>
        {isModalJoin || isModalLogin || isModalPassword ? (
          <div className="modal-box">
            <div className="auth-box">
              {isModalJoin ? (
                <Join />
              ) : isModalLogin ? (
                <Login />
              ) : isModalPassword ? (
                <Password />
              ) : null}
              <button
                className="auth-close-btn iconfont icon-close"
                onClick={this.closeAuth}
              />
            </div>
          </div>
        ) : null}
      </div>
    );
  }
}

export default connect((state) => ({ modal: state.modal }), { modalClose })(
  App
);
