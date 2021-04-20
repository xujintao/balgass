import React from 'react'
import {Route,BrowserRouter,Switch,Redirect} from 'react-router-dom'
import Nav from './nav/nav'
import Main from './main/main'
import Search from './search/search'
import Footer from './footer/footer'
import  './App.css'

class App extends React.Component {
  render() {
    return (
      <BrowserRouter>
        <Nav/>
        <Switch>
          <Route path="/main" component={Main}/>
          <Route path="/search" component={Search}/>
          {/* <Route path="/bugs" component={Bugs}/> */}
          {/* <Route path="/download" component={Download}/> */}
          <Redirect to="/main"/>
        </Switch>
        <Footer/>
      </BrowserRouter>
    )
  }
}

export default App
