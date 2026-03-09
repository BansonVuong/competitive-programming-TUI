[Judge](https://dmoj.ca/problem/ccc24s3)

# CCC24S3

swipe left, swipe right

l, r -> interval dp?

okay
given an array with distinct sections a, b, c, d

the original array must still have the same order of sections a, b, c, d

essentially its 2 pointer, one on each array

1 2 3 4 5 2 6 4
1 1 2 4 4 5 2 4

advance first one pointer until you reach a new section

1 2 3 4 5 2 6 4
1 1 2 4 4 5 2 4
    ^

now search in original essay for section's number

1 2 3 4 5 2 6 4
^
1 1 2 4 4 5 2 4
    ^
store. 1 swipe right to 2

next step, advance to end of 2 and search
1 2 3 4 5 2 6 4
  ^
1 1 2 4 4 5 2 4
      ^
store. 2 swipe right to 3 (how to order?)

next step
1 2 3 4 5 2 6 4
      ^
1 1 2 4 4 5 2 4
          ^
store. swipe right to 5

next step
1 2 3 4 5 2 6 4
        ^
1 1 2 4 4 5 2 4
            ^
store. swipe right to 5

etc.

impossible:

1 1 1 2 2 2
      ^
2 2 2 1 1 1
      ^
2 swipe left to 1

1 1 1 2 2 2
           ^
2 2 2 1 1 1
           ^
first pointer reached end. impossible.

ordering:
right swipes, go reverse order (furthest right first)
left swipes, forwards (furthest left first)

(done in ~10 minutes)

## **Claude's read on CCC '24 S3 — Swipe:**

**Got on your own:** The full algorithm in 10 minutes. Two-pointer matching of contiguous blocks between A and B, determine swipe direction based on whether B's block extends past A's block. YES/NO detection. Banked 7 points for subtask 1 quickly.

**What killed you:** Three directional bugs, all the same root cause — you wrote comparisons and orderings from the wrong perspective every single time:
1. `end < ed` instead of `ed < end` — condition for "B extends past A on the right" was flipped
2. `{beginning, start}` instead of `{start, beginning}` — left swipe pair had l > r, violating constraint
3. Sort order backwards — right swipes need rightmost-first (descending) to avoid corrupting source values, left swipes need leftmost-first (ascending). You had them swapped.

**What was NOT a problem:** The algorithm, the YES/NO logic, the block-matching structure, the two-pointer advancement. All correct from the start.

**Spoonfed:** All three bugs were identified for you. You didn't catch any of them through your own tracing.

**Recurring pattern:** This is the same failure mode as every other problem today. You write directional logic (comparisons, pair orderings, sort orders) by gut feeling instead of tracing one concrete example through the line. Three seconds of "if beginning=3 and start=1, which goes first?" would have caught all three bugs instantly.

**Takeaway for Wednesday:** Every `<`, `>`, `{a,b}`, and sort comparator you write — pause, pick two concrete numbers, trace that single line. Not the whole algorithm. Just that line. This one habit would have saved you 30+ minutes today across two problems.
