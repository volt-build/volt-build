{
  description = "mini-build: A small build system I wrote for myself to run repetitive tasks easily";

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
          pkgs.go
          pkgs.golangci-lint
          pkgs.go-tools
          pkgs.gopls
          pkgs.delve
        ];
      };
      packages.default = pkgs.buildGoModule {
        pname = "mini-build";
        version = "0.1.0";
        src = ./.;
        vendorHash = null;

        meta = {
          description = "A small build system to make running repetitive tasks easier without polluting PATH or making it hard to run commands with scripts."; 
          license = lib.licenses.mit; 
          maintainers = [ ];
          platforms = lib.platforms.unix; 
        }; 
      };
      formatter = pkgs.alejandra;
    });
}
