# Contributing to mini-build 

First of all, I want to thank you for your interest in contributing to mini-build! 
Just as a heads up: This project is fairly low level (With things like assembly, DAG's, and complex multi-threading), Its not very beginner friendly. 
I would love to see any contributions tho! 

Requirements: 

Go 1.23+ installed 
Fasm (Not required unless you're contributing to the code generation sector of mini-build with target x86_64-linux) 
(* optional but recommended for the exact dev environment I use: direnv && nix-direnv *) 

## What to work on? 

Things that are welcome: 

- New graph generation with [graphviz](https://graphviz.org/)
- Help with compilation 
- Improvements to the DAG (package `arena`) 
- `SnapShot()` function for the `Arena` struct 
- Fixing of things in [TODO.md](TODO.md), or fixing of the TODO's in the code
- Better Error Messages (Very very appreciated if you do this!, This feature is very important for any kind of language) 
- Docs/Docs Website help (Different repo) 

Things that are not: 

- making of GUI for it in the main repo (Other repos are completely your wish) 
- making random web API's in the source code (highly discouraged) 
- inclusion of `cgo` ( I'm not against Cgo/C, I just don't like using it because it causes a lot of problems in my debugger setup)
- making flags with no purpose 

# **Code Style:**
### Conventions: 

- Use gofumpt and goimports when you can (Integration of those into your editor will be preferred)
- Tabs? Spaces? -- The honest answer is: I don't really care as long as it appears as 4 spaces in the end.
- Every `Future` must only share its data between one other future or no other future 
- Chain `Arena`s only when its a huge task which would take 5~10+ futures to complete 
- Please just make a notice when you're adding a dependency in go.mod (And if you can, update the vendorHash in flake.nix) 


> If I'm wrong ever in the source code. Please feel free to point it out! I'd love to learn! 
