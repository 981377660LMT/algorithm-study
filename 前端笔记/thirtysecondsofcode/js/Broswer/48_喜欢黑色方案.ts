// 在适当的媒体查询中使用 Window.matchMedia ()检查用户的颜色方案首选项。
const prefersDarkColorScheme = () => window?.matchMedia('(prefers-color-scheme: dark)')?.matches
prefersDarkColorScheme() // true
