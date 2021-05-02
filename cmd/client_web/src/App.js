import React from "react";
import { Route, BrowserRouter, Switch, Redirect } from "react-router-dom";
import Nav from "./nav/nav";
import Home from "./home/home";
import Search from "./search/search";
import Footer from "./footer/footer";
import "./App.css";

class App extends React.Component {
  render() {
    return (
      <BrowserRouter>
        <Nav />
        <Switch>
          {/* https://stackoverflow.com/questions/49162311/react-difference-between-route-exact-path-and-route-path */}
          <Route exact path="/" component={Home} />
          <Route path="/search" component={Search} />
          {/* <Route path="/bugs" component={Bugs}/> */}
          {/* <Route path="/download" component={Download}/> */}
          <Redirect to="/" />
        </Switch>
        <Footer />
      </BrowserRouter>
    );
  }
}

export default App;
