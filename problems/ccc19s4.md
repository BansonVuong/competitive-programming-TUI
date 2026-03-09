[Judge](https://dmoj.ca/problem/ccc19s4)

# CCC 19 S4

first instinct - maximizing and minimizing, dp?

dp[attraction][day], solution is dp[n][n/k+1]

working through samples: essentially if n%k == 0, the solution is already predetermined. nothing complicated required.

if n%k != 0, that means there is a block of length n%k we can play around with. we can shift this block around, inserting it between blocks of length 3 to figure out exactly which is most optimal

naive: shift each special block and then recompute

compute blocks of length k twice:once from beginning, once from end. store sum in psa.

then for each position of the special block, just sum up the left side, right side, and find the value for the middle block

okay so none of that was true - i was assuming blocks of length k and then one block of length n%k, but that’s just not true.

it’s still dp.

i spent too long on this so I just got spoonfed the dp transition - dp by endpoint

#days = ceil(n/k)

each day, the endpoint is in between [days*k + r-1, days * k + k-1] (0 indexed)

days*k - standard # of attractions in that day

+r-1 - minimum day length, subtract 1
+k-1 - maximum day length, subtract 1

okay i completely gave up and let AI spoonfeed me the dp transition because this was taking entirely too long. I need to move on.
