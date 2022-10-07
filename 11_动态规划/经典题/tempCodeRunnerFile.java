 for (int d = 0; d < 2; d++) {
          for (int c = 0; c < 2; c++) {
            for (int e = 0; e < 5; e++) {
              for (int t = 0; t < 2; t++) {
                for (int o = 0; o < 3; o++) {
                  for (int l = 0; l < 4; l++) {
                    for (int h = 0; h < 2; h++) {
                      dp[c][d][e][t][o][l][h] = ndp[c][d][e][t][o][l][h];
                    }
                  }
                }
              }
            }
          }
        }
      }