= n; j++ {
				dist[s][j] += dist[t][j]
				dist[j][s] += dist[j][t]
			}