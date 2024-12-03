# Advent of Code 2024

![Advent of Code](https://img.shields.io/badge/Advent%20of%20Code-2024-brightgreen)
![Language](https://img.shields.io/badge/Language-Go-blue)
![Days Completed](https://img.shields.io/badge/Days%20Completed-3-orange)

## About Advent of Code

[Advent of Code](https://adventofcode.com/) is an annual event where participants solve daily programming puzzles, each
released in the form of a two-part challenge. The challenges cover various topics, and participants often use the
opportunity to sharpen their problem-solving and coding skills.

## About My Approach

I was aiming to:

* utilize my knowledge of algorithms and data structures;
* make solutions as general as I could;
* prioritize code readability over how fast I can write the solution — I didn't try to compete on global leaderboard.

### Day 2: Red-Nosed Reports

In part 2 my idea was to iterate solely over differences of adjacent levels in reports, and look for anomalies. In case
of an anomaly is found, we should remove some level, and simply reuse solution from part 1 on the new report. In
general, when we encounter an anomaly, it consists of one or two diffs, and therefore at most three levels might
contribute to it. I simply tro to remove each level, and if any of the resulting reports is safe, then we found a level
that is safe to remove.

We can do it, because we can't skip an anomaly once we found it. If we wouldn't fix it, we'd have an unsafe report in
the beginning.

Solution runs in `O(N)` time of input data, and in the worst case processes all the input 4 times. The worst case is an
anomaly in the beginning that can be fixed by removing 3rd level in the corresponding report.
