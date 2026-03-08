# ccc16s4

interval dp
if two rice balls same size, can combine.
2 rice balls 

base case: dp[i][i] = val[i]

n=400 - o(n^3) complexity?

transitions:

dp[1][n] = max riceball size given the whole range

dp[i][i] = val[i]
dp[i][i+1] = max(i, i+1) if val[i] != val[i+1] else val[i]*2
dp[l][r] = 

basically we can have either
dp[l][k] == dp[k+1][r] and dp[l][r] = sum(l, r)
dp[l][a] == dp[b][r] and dp[a+1][b-1] is one big ball, dp[l][r] = sum(l, r)



```cpp
for(int len=2; len<n; len++){
    for(int left=1; left<n; left++){
        
    }
}
```

iterate over length of ball, from 2 to N. ball of length 1 is just value of itself.

check if a ball is a ball by checking if its dp interval is equal to its sum

for every length, have a normal interval split point
    check that leftball is a ball, right ball is a ball
    and then if they are equal merge: sum[l, r]

if length is more than or equal to 3, you can have a middle ball.
    
    if middleball, leftball, and rightball all exist and leftball == rightball, dp[l][r] is just sum the sum from left to right.
