import * as React from 'react'
import { useImperativeHandle, useRef, useState } from 'react'
import * as ReactDom from 'react-dom'

interface IChildHandle {
  addCount: () => void
}
interface IChildProps {}
const Child = React.forwardRef<IChildHandle, IChildProps>((_props, ref) => {
  const [count, setCount] = useState(0)
  const addCount = () => setCount(count => count + 1)
  // 子组件通过useImperativeHandle函数，将addCount函数添加到父组件中的ref.current中
  useImperativeHandle(ref, () => ({ addCount }), [])

  return (
    <div>
      {count}
      <button onClick={addCount}>I am child</button>
    </div>
  )
})

interface IApp2Props {}

const App: React.FC<IApp2Props> = () => {
  const childRef = useRef<IChildHandle>(null)
  const clickHandle = () => childRef.current?.addCount()

  return (
    <div>
      <Child ref={childRef} />
      <button onClick={clickHandle}>点击调用子组件方法</button>
    </div>
  )
}

const element = React.createElement(App)
ReactDom.render(element, document.querySelector('#app'))
