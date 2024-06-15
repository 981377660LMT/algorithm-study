/* eslint-disable no-shadow */

function transform<From extends { children?: From[] | null | undefined }, To extends object>(
  from: Iterable<From>,
  mapFrom: (from: From) => To
): To[] {
  const res: To[] = []
  for (const f of from) dfs(f)
  return res

  function dfs(cur: From): void {
    const to = mapFrom(cur)
    res.push(to)
    cur.children && cur.children.forEach(dfs)
  }
}

if (require.main === module) {
  interface Route {
    id: number
    url: string
    path: string
    component: string
    name: string
    icon: string
    children: Route[] | null
  }

  type To = Pick<Route, 'path' | 'name' | 'children'> & {
    component: () => Promise<unknown>
  }

  const res = transform<Route, To>(
    [
      {
        id: 1,
        url: '/',
        path: '/home',
        component: 'Home',
        name: '员工资料',
        icon: 'dsadsa',
        children: [
          {
            id: 7,
            url: '/sdjajdi',
            path: '/home/son',
            component: 'Homeson',
            name: '首页子路由',
            icon: 'dsadsa',
            children: null
          }
        ]
      }
    ],
    from => {
      const { path, name, children } = from
      const component = () => import(`${from.component}`)
      return { path, name, component, children }
    }
  )

  console.dir(res)
}

export {}
