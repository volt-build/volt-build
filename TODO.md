## TODO.md

#### Stuff I need to do.  

- [x] Lexer 
- [x] Language basics 
### [ ] DOCS!!!
- [x] command flags // probably includes the next one too. (for now, needs testing) 
- [ ] hook up more of the helpers to flags 
- [ ] more sensible command defaults.
- [ ] fix assignmet nodes 
- [ ] Update the really old styled flake


#### QoL additions that can wait a bit. 

- [x] Make it run on different goroutines (probably makes it faster, I think, makes me feel better atleast) 
explaination: 
    it spins up a goroutine every time a shell command is executed, well i dont expect people to be using
    small shell commands so fair thing, everything else which would already be in shell is builtin
    to the language.

# Fetch keyword 
how it should work: 
    fetch stuff from a code hosting place 

> which code hosting place? 
> should be able to mention in the fetch argument. 


