## Commit 2025-05-31: Big optimizations 

#### Big Changes: 
    - Optimize compile statements to work with wait groups and make them non blocking: 
        They run on a different OS thread now without blocking I/O with blocking channels in spawnCompile 
    - Look at $0 (first positional argument in the command) to determine the path of the file 
    - Rename `evaluateCompile` to `spawnCompile`. 

<small>TODO's in (TODO.md)[TODO.md]</small>
<small>This is a new file to hold commits which can't be described in a single sentence</small>
