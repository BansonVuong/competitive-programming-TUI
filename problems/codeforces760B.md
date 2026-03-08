# codeforces760B

n beds, m pillows

same bsearch check logic.

this was highkey implementation hell

Claude's take

**What you got on your own:** The core algorithm — binary search on Frodo's pillow count, greedy "triangle" distribution tapering by 1 each side, subtract the excess when the triangle extends past the row ends. That's the hard part of this problem, and you nailed it.

**What you needed help with:** Three bugs, all implementation-level:
1. Variable name shadowing (`m`) — classic C++ gotcha
2. Negative × negative in the subtraction formula sneaking past your `max(0, ...)` guard
3. Overflow from an unnecessarily large binary search bound

None of these were conceptual gaps — you understood the problem and the approach. Just the kind of fiddly details that bite everyone in competitive programming when you're coding fast.
