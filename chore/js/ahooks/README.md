你可以按照以下步骤在本地创建一个 React + TypeScript 项目并集成 ahooks：

1. **创建 React + TypeScript 项目**  
   在终端执行：

   ```sh
   npm create vite@latest my-ahooks-demo -- --template react-ts
   ```

2. **进入项目目录**

   ```sh
   cd my-ahooks-demo
   ```

3. **安装 ahooks**

   ```sh
   npm install react@18 react-dom@18
   npm install ahooks
   ```

4. **启动项目**

   ```sh
   npm run dev
   ```

5. **在代码中使用 ahooks**  
   编辑 `src/App.tsx`，引入并使用 ahooks 提供的 hook，例如：

   ```tsx
   import { useRequest } from 'ahooks'

   function App() {
     const { data, error, loading } = useRequest(() => fetch('https://api.github.com').then(res => res.json()))

     if (loading) return <div>Loading...</div>
     if (error) return <div>Error!</div>
     return <pre>{JSON.stringify(data, null, 2)}</pre>
   }

   export default App
   ```
