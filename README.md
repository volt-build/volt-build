# mini-build

Its a small build system I wrote myself. 

### Writing build files

- Syntax: 

```task

task build {
    # this is a comment 
    shell "# command to execute with sh -c over here" 
    foreach "./*go" { 
        push "found a go file" 
    }
    push "Something to print out stdout" 
}

exec build  # make sure to execute the task at the end. 
```


