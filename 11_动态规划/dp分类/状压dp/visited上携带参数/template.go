package main

func main() {
	var nums []int
	const INF int = 1e18
	n := len(nums)
	mask := 1<<n - 1

	memo := make([]int, 1<<n)
	for i := range memo {
		memo[i] = -1
	}
	hash := func(visited int) int {
		return visited
	}

	var dfs func(visited int) int
	dfs = func(visited int) int {
		if visited == mask {
			return 0
		}
		hash_ := hash(visited)
		if memo[hash_] != -1 {
			return memo[hash_]
		}
		res := INF
		for next := 0; next < n; next++ {
			if visited&(1<<next) == 0 {
				res = min(res, dfs(visited|(1<<next))+1)
			}
		}
		memo[hash_] = res
		return res
	}

}
