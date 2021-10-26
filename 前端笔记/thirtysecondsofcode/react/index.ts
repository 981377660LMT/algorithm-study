import React from 'react'
import ReactDOM from 'react-dom'
import { UseStateWithLable } from './src/hook/1_useStateWithLable'

const virtualDom = React.createElement(UseStateWithLable)
ReactDOM.render(virtualDom, document.getElementById('app'))
