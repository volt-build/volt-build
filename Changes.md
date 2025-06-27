## Commit 2025-05-31: Big optimizations 

#### Big Changes: 
- Optimize compile statements to work with wait groups and make them non blocking: 
    They run on a different OS thread now without blocking I/O with blocking channels in spawnCompile 
- Look at $0 (first positional argument in the command) to determine the path of the file -- Doesn't work completely 
- Rename `evaluateCompile` to `spawnCompile`. 
- Use env var for path

<small>TODO's in (TODO.md)[TODO.md]</small>
<small>This is a new file to hold commits which can't be described in a single sentence</small>

## Commit on 2025-06-11: Compiler Directory 

#### This directory will implement: 

- A compiler which depends on codegen directory 
- Parallel compilation 
- A new language system and a semantic analyzer 
- Possibly making the language statically typed (Similar to go itself) 

## Commit on 2025-06-27: Compiler parser update and filename changed to build.mb instead of Taskfile

#### NOTE: Theres going to be a new scheduler-like thing in a new GitHub organisation for complex compilation tasks 
