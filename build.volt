task build { 
    push "Building.."
	shell "mkdir -p ./build/" 
	shell "echo '*' >> ./build/.gitignore"
    compile "." "go build -o build/"
    push "Done with exit code: " ++ $?
}


task something requires fmt {
	shell "cat main.go" 
}

task fmt {
    push "starting formatting.... "
    foreach "./*.go" gofile {
        push "Formatting: " ++ gofile
		# Run 2 formatters 
        shell "gofumpt -w " ++ gofile
		shell "goimports -w " ++ gofile 
    }
    push "Done with exit code: "  ++ $?
}

task termlint requires fmt {
    push "Linting code"
	compile "./..." "golangci-lint run"
    push "Done with exit code: "  ++ $?
}

task build_system_test {
	var = "variable string dude" 
	push var 
	if 1 { # boolean condition 
		push "true" 
	}
}
