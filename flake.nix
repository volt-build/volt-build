{
  # fixed the wrong desc 
  description = "Go dev shell for direnv";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }: 
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.golangci-lint
            pkgs.go-tools 
            pkgs.gopls 
          ];

          shellHook = ''
            if [ -z "$IN_ZSH" ]; then
              export IN_ZSH=1
              exec zsh -i
            fi
          '';
        };
      });
}

