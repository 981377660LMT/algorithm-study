import * as React from 'react'
import * as ReactDom from 'react-dom'

interface IApp2Props {}

const App: React.FC<IApp2Props> = props => {
  const [val, setVal] = React.useState(0)

  React.useEffect(() => {
    setVal(val => val + 1)
    console.log(val) // 闭包 输出0
    setVal(val => val + 1)
    console.log(val) // 闭包 输出0
    setTimeout(() => {
      setVal(val => val + 1)
      console.log(val) // 闭包 输出0
      setVal(val => val + 1)
      console.log(val) // 闭包 输出0
    }, 0)
  }, [])

  return <div>{val}</div> // 4
}

const element = React.createElement(App)
ReactDom.render(element, document.querySelector('#app'))
