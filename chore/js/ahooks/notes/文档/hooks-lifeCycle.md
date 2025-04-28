## useMount

```tsx
import { useMount, useBoolean } from 'ahooks'
import { message } from 'antd'
import React from 'react'

const MyComponent = () => {
  useMount(() => {
    message.info('mount')
  })

  return <div>Hello World</div>
}

export default () => {
  const [state, { toggle }] = useBoolean(false)

  return (
    <>
      <button type="button" onClick={toggle}>
        {state ? 'unmount' : 'mount'}
      </button>
      {state && <MyComponent />}
    </>
  )
}
```

## useUnmount

## useUnmountedRef

获取当前组件是否已经卸载的 Hook。

```tsx
import { useBoolean, useUnmountedRef } from 'ahooks'
import { message } from 'antd'
import React, { useEffect } from 'react'

const MyComponent = () => {
  const unmountedRef = useUnmountedRef()
  useEffect(() => {
    setTimeout(() => {
      if (!unmountedRef.current) {
        message.info('component is alive')
      }
    }, 3000)
  }, [])

  return <p>Hello World!</p>
}

export default () => {
  const [state, { toggle }] = useBoolean(true)

  return (
    <>
      <button type="button" onClick={toggle}>
        {state ? 'unmount' : 'mount'}
      </button>
      {state && <MyComponent />}
    </>
  )
}
```
