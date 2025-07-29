{
  description = "volt-build";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = import nixpkgs {inherit system;};
      lib = pkgs.lib;
    in {
      devShells.default = pkgs.mkShell {
        packages = [
          # lot of tools, use every single one of them.
          pkgs.go
          pkgs.golangci-lint
          pkgs.gotools
          pkgs.gopls
          pkgs.delve
          pkgs.gofumpt
        ];
      };
      packages.default = pkgs.buildGoModule {
        pname = "volt-build";
        version = "0.1.1";
        src = ./.;
        vendorHash = "sha256-Mq1r4dSHx8x5EiqXAecrcOXYIZiUwLF1xvaeYKQYqW0=";

        meta = {
          description = "A small build system to make running repetitive tasks easier without polluting PATH or making it hard to run commands with scripts.";
          license = lib.licenses.mit;
          maintainers = [];
          platforms = lib.platforms.unix;
        };
      };
      formatter = pkgs.alejandra;
    });
}
