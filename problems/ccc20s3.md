# ccc20s3 

okay first glance.. a distinct permutation of N is just a window of length N with the same frequency count as N

so its just a frequency counter and sliding window

(3 minutes)

debugging lessons:
substring is beginning index and length. python string ops strikes again.
learned how to hash: just multiply everything by 131 and add the next bit
also that hash is expensive time wise and actually adds an o(N) time to complexity. would have seen it if i didnt abstract it away in a function.
also that substring is equally expensive.. sigh
learned sliding window where you just add the start and remove the end as you slide along the string.

RUBBER DUCK DEBUGGING WORKS WONDERS

**Claude's read on CCC '20 S3:**

**Got on your own:** The entire algorithm — frequency matching, sliding window, dedup for distinct permutations. Approach identified in 3 minutes. That's fast and correct.

**What killed you:** Implementation choices, not algorithmic thinking. You reached for a 200K×26 prefix sum array when a 26-element sliding frequency array was the natural tool. Then you layered on `substr` and `hsh` calls without thinking about their cost — O(N) each, called O(N) times = O(N²). You abstracted away the expense inside a function and lost visibility into it.

**Specific bugs:**
1. `substr(start, length)` not `substr(start, end)` — Python muscle memory
2. `substr` placed above the `if(good)` check twice after being told to move it
3. MLE from storing full strings in an unordered_map
4. TLE from not recognizing that string operations are O(N), not O(1)

**Pattern learned:** Sliding window with a `matches` counter for frequency equality. This is cleaner and faster than prefix sums for fixed-window frequency problems. You already knew sliding window from '22 S4 but didn't recognize it in this wrapper.

**One habit: before writing a line inside a loop, ask "what's the complexity of this operation × how many times does it run?"** Would've caught both the substr and the prefix sum issue instantly.
