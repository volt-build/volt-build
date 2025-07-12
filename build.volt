task build { 
    push "Building.."
    compile "." "go build -o target/"
    push "Done with exit code: " ++ $?
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

task termlint require fmt {
    push "Linting code"
	compile "./..." "golangci-lint run"
    push "Done with exit code: "  ++ $?
}

