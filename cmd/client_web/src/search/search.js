import React from 'react'
import axios from 'axios'
import './search.css'

class Search extends React.Component {
    state = {
        users: [],
        isFirst: true,
        isLoading: false,
        err: ''
    }

    refInput = React.createRef()
    handleSearch = (event)=>{
        event.preventDefault()
        this.setState({
            isFirst: false,
            isLoading: true
        })
        const keyword = this.refInput.current.value
        axios.get(`/api1/search/users?q=${keyword}`).then(
            response=>{
                this.setState({
                    users: response.data.items,
                    isLoading: false
                })
            },
            error=>{
                this.setState({
                    isLoading: false,
                    err: error.message
                })
            }
        )
    }
    render() {
        const {users,isFirst,isLoading,err} = this.state
        return (
            <div className="w">
                <div className="search-head">
                    <form action="/search" onSubmit={this.handleSearch}>
                        <input type="text" placeholder="我的发明可以让整个小区停电" name="keyword" ref={this.refInput}/>
                        <button>搜索</button>
                    </form>
                </div>
                <ul className="search-result">
                {
                    isFirst ? <h2>欢迎使用，输入关键字，随后点击搜索</h2> :
                    isLoading ? <h2>Loading...</h2> :
                    err ? <h2>{err}</h2> :
                    users.map((u)=>{
                        return (
                            <li key={u.id} className="card">
                                <a href={u.html_url} target="_blank" rel="noreferrer">
                                    <img src={u.avatar_url} style={{width:'100px'}} alt=""/>
                                </a>
                                <p className="card-text">{u.login}</p>
                            </li>
                        )
                    })
                }
                </ul>
            </div>
        )
    }
}

export default Search