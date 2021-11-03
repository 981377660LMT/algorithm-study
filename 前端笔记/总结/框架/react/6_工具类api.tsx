import * as React from 'react'
import * as ReactDom from 'react-dom'

interface IApp2Props {}

const App: React.FC<IApp2Props> = props => {
  const [val, setVal] = React.useState(0)
  const a = React.createElement('a', { id: 1 })
  console.log(React.cloneElement(a, { id: 12 }))
  console.log(React.isValidElement(a))

  return <div>{val}</div> // 4
}

const element = React.createElement(App)
ReactDom.render(element, document.querySelector('#app'))
