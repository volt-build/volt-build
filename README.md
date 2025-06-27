<h1 align="center"> mini-build </h1> 
<small>A build system I wrote because run-tasks.sh was too much effort</small> 


# Note: 
Project might be renamed and changed location of (to an org) 

<h3 align="left">A few examples on usage</h1>

- A task to format a directory with go files: 
```task
# Task declaration:
#     ┌─▶ 'fmt' is the name of the task
task fmt {
    # foreach loop:
    #        ┌─▶ glob pattern
    #        │        ┌─▶ loop variable (each matched file)
    foreach "./*.go" gofile {
        # shell command:
        #       ┌─────────────▶ command
        #       │        ┌───────▶ concatenation operator
        #       │        │    ┌──▶ variable reference
        shell "gofmt -w" ++ gofile
    }
}
```

- A task to lint a directory with go files: 
```task 
task lint {
    push "Linting..." 
    shell "golang-ci run ./..." 
    push "Done with exit code: " ++ $? 
}
```

- A task to build c files from `src/`: 
``` 
task buildc {
    output_dir = "bin" # example 
    shell "mkdir " ++ output # make sure it exists 
    cc = "your-compiler" 
    flags "-flags -for -your -compiler" 
    foreach "./src/*.c" cfile {
        cfile_with_o = cfile ++ ".o" 
        # a little cursed but it works pretty good: 
        ##      (             ONE VARIABLE              ) (var2) 
        compile cc ++ flags ++ " -o obj/" ++ cfile_with_o  cfile
        compile cc ++ " obj/*" ++ " -o " ++ output_dir
    }
}
```

Usage: 

- Put these in a build.mb file in the CWD and just mini-build -t `<TaskName>`! 


> This is was designed to be as simple as possible, but with no YAML/TOML/JSON/GNU make 
