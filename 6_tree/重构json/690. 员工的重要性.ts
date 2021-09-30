class Employee {
  id: number
  importance: number
  subordinates: number[]
  constructor(id: number = 0, importance: number = 0, subordinates: number[] = []) {
    this.id = id
    this.importance = importance
    this.subordinates = subordinates
  }
}
/**
 *
 * @param employees  它包含了员工 唯一的 id ，重要度 和 直系下属的 id 。
 * @param id
 * 输入一个公司的所有员工信息，以及单个员工 id ，返回这个员工和他所有下属的重要度之和。
 */
function getImportance(employees: Employee[], id: number): number {
  const idToEmployee = new Map<number, Employee>()
  employees.forEach(employee => idToEmployee.set(employee.id, employee))

  let res = 0
  const dfs = (id: number) => {
    const employee = idToEmployee.get(id)!
    const { importance, subordinates } = employee
    res += importance
    for (const child of subordinates) {
      dfs(child)
    }
  }
  dfs(id)
  return res
}

console.log(
  getImportance([new Employee(1, 5, [2, 3]), new Employee(2, 3, []), new Employee(3, 3, [])], 1)
)
export {}
