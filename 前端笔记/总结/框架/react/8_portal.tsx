import * as React from 'react'
import { useEffect, useState } from 'react'
import * as ReactDom from 'react-dom'

// 如果一个子组件的真实DOM结构必须渲染到当前组件外，但又想保留这两者的父子关系，就可以用Protals。
const Modal: React.FC<{}> = props => {
  const container = document.createElement('div')
  document.body.appendChild(container)

  useEffect(() => {
    return () => {
      document.removeChild(container)
    }
  }, [])

  return ReactDom.createPortal(props.children, container)
}

const App: React.FC<{}> = () => {
  const [count, setCount] = useState(0)
  const add = () => setCount(count => count + 1)
  return (
    <div onClick={add}>
      <Modal>加一{count}</Modal>
    </div>
  )
}

const element = React.createElement(App)
ReactDom.render(element, document.querySelector('#app'))
