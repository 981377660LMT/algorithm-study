interface Route {
  id: number
  url: string
  path: string
  component: string
  name: string
  icon: string
  children: Route[] | null
}

type RoutesArray = (Pick<Route, 'path' | 'name' | 'children'> & {
  component: () => Promise<unknown>
})[]

const initRoutes = (routes: Route[]): RoutesArray => {
  const newRouter: RoutesArray = []
  const dfs = (root: Route) => {
    const { path, name, children } = root
    const component = loadRoutes(root.component)
    newRouter.push({ path, name, component, children })
    root.children && root.children.forEach(dfs)
  }
  routes.forEach(dfs)
  return newRouter
}

const loadRoutes = (component: string) => () => import(`${component}`)

const routes: Route[] = [
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
]

console.dir(initRoutes(routes))
export default 1
