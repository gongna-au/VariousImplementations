package TSP

//动态规划
import "fmt"

func DP(dist [][]int, st int) []int {
	const inf int = 1e9
	dp := make([][]int, 1<<len(dist))
	for i := range dp {
		dp[i] = make([]int, len(dist))
		for j := range dp[i] {
			dp[i][j] = inf
		}
	}
	fmt.Println(dp)
	dp[1<<st][st] = 0 // 多个起点的话就设置多个 dp[1<<st[i]][st[i]] = 0
	for s, dr := range dp {
		// 利用位运算快速求出 s 中 1 的位置 i，以及 s 中 0 的位置 j（通过 s 的补集中的 1 的位置求出）
		for ss := uint(s); ss > 0; ss &= ss - 1 {
			i := bits.TrailingZeros(ss)
			for t, lb := len(dp)-1^s, 0; t > 0; t ^= lb {
				lb = t & -t
				j := bits.TrailingZeros(uint(lb))
				dp[s|lb][j] = min(dp[s|lb][j], dr[i]+dist[i][j])
			}
		}
	}
	return []int{}
}
