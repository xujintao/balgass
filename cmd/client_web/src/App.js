import React from 'react'
import Nav from './nav/nav'
import Search from './search/search'
import Footer from './footer/footer'
import  './App.css'

class App extends React.Component {
  render() {
    return (
      <div>
        <Nav/>
        <Search/>
        <Footer/>
      </div>
    )
  }
}

export default App
