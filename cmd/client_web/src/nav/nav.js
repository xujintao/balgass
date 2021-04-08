import React from 'react'
import './nav.css'
import logo from './images/logo.png'
import './fonts/iconfont.css'

class Nav extends React.Component {
    render() {
        return (
            <div className="nav">
                <div className="nav-link">
                    <ul>
                        <li>
                            <a href="index.html">
                                <img className="nav-link-logo" src={logo} alt=""/>首页
                            </a>
                        </li>
                        <li><a href="#">bug报告</a></li>
                        <li><a href="#">游戏下载</a></li>
                    </ul>
                </div>
                <div className="nav-search">
                    <form className="nav-search-form" action="">
                        <input className="nav-search-text" type="text" placeholder="我的发明可以让整个小区停电"/>
                        <button className="nav-search-btn iconfont icon-search"></button>
                    </form>
                </div>
                <div className="nav-user">
                    <div className="nav-user-login">
                        <ul>
                            <li>
                                <a href="#">
                                    <img className="nav-user-avatar" src="images/avatar.webp" alt=""/>
                                </a>
                            </li>
                            <li><a href="#">消息</a></li>
                            <li><a href="#">动态</a></li>
                            <li><a href="#">收藏</a></li>
                        </ul>
                    </div>
                    <div className="nav-user-logout">
                        <ul>
                            <li><a href="login.html">登录</a></li>
                            <li><a href="register.html">注册</a></li>
                        </ul>
                    </div>
                </div>
            </div>
        )
    }
}

export default Nav