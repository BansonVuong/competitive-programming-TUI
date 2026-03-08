# codeforces1201c

check -> is this median possible with k ops?

how to check median:

1 1 1 1 1 -> 1
1 1 1 1 2 -> 1
1 1 1 2 2 -> 1
1 1 2 2 2 -> 2
1 1 2 2 3 -> 3

essentially indice n/2 and below: ignore

add to indice n/2+1 until its equal to n/2+2, etc.

(2 minutes)

i didnt necessarily get the recognition because i was told this was binary search on answer 

i set bsearch right to 1e9 which was wrong, 1e9 is k, but its possible for it to be up to 4e10 for ans

**Claude's take on 1201C:**

**Got on your own:** The check function logic (only need to raise indices n/2 and above), the implementation, clean binary search structure. 20 min, no coaching needed beyond problem selection.

**What was handed to you:** The pattern recognition — you were told it was binary search on answer. On contest day, you need to see "maximize the median" and think binary search yourself. The trigger phrase is there ("maximize X" where X can be checked monotonically).

**Bug caught:** Right bound of 1e9 instead of ~4e10. This is the same bounds issue from CSES 1085. **Check constraints before setting bounds** — you've now made this mistake twice.
