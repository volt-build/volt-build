## TODO.md

#### Stuff I need to do.  

- [x] Lexer 
- [x] Language basics 
- [x]  Fix some bugs (I just caused those) 
- [ ] fetch keyword to fetch stuff from websites (see [this](https://github.com/justachillguy57/Taskr/blob/master/README.md#fetch-keyword)) 
- [ ] command flags
- [ ] more sensible command defaults.


#### QoL additions that can wait a bit. 

- [x] Make it run on different threads for speed.
explaination: 
    it spins up a goroutine every time a shell command is executed, well i dont expect people to be using
    small shell commands so fair thing, everything else which would already be in shell is builtin
    to the language.

# Fetch keyword 
how it should work: 
    fetch stuff from a code hosting place 

> which code hosting place? 
> should be able to mention in the fetch argument. 


