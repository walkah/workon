{
  description = "Manage tmux for what you work on.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    {
      overlays.default = final: prev: {
        inherit (self.packages.${prev.system}) workon;
      };
    }
    // flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version = "0.4.0";
      in
      {
        packages.workon = pkgs.buildGoModule {
          pname = "workon";
          version = version;
          src = ./.;
          vendorHash = "sha256-p3bwgNUbUsgJfANIOxG/+GXeh2/OddxYs0EK3JfnD4s=";

          ldflags = [
            "-s -w -X github.com/walkah/workon/cmd.version=${version}"
          ];

          nativeBuildInputs = with pkgs; [
            tmux
            installShellFiles
          ];

          postInstall = ''
            for shell in bash fish zsh; do
              $out/bin/workon completion $shell > workon.$shell
              installShellCompletion --$shell workon.$shell
            done
          '';
        };

        packages.default = self.packages.${system}.workon;

        devShells.default = pkgs.mkShell {
          name = "workon";
          buildInputs = with pkgs; [
            cobra-cli
            go
            gopls
          ];
        };
      }
    );
}
