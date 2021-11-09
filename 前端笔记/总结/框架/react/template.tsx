import * as React from 'react'
import * as ReactDom from 'react-dom'

interface IApp2Props {}

const App: React.FC<IApp2Props> = props => {
  const [val, setVal] = React.useState(0)

  return <div>{val}</div> 
}

const element = React.createElement(App)
ReactDom.render(element, document.querySelector('#app'))
