import * as React from 'react'
import { useState } from 'react'
import * as ReactDom from 'react-dom'

interface IApp2Props {}

const App: React.FC<IApp2Props> = props => {
  const { value: userName, reset: resetUserName, bind: bindUserName } = useInput('')
  const { value: password, reset: resetPassword, bind: bindPassword } = useInput('')

  const sumbitHanle = (event: React.ChangeEvent<HTMLFormElement>) => {
    event.preventDefault()
    console.log(userName, password)
    resetUserName()
    resetPassword()
  }

  return (
    <form onSubmit={sumbitHanle}>
      <label htmlFor="userName">userName:</label>
      <input id="userName" type="text" {...bindUserName} />
      <label htmlFor="password">password:</label>
      <input id="password" type="password" {...bindPassword} />
      <input type="submit" value="login" />
    </form>
  )
}

function useInput(initValue: string) {
  const [value, setValue] = useState(initValue)

  const reset = () => setValue(initValue)
  const bind = {
    value,
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => setValue(event.target.value),
  }

  return { value, reset, bind }
}

const element = React.createElement(App)
ReactDom.render(element, document.querySelector('#app'))
