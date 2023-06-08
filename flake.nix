{
  description = "Manage tmux for what you work on.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/release-23.05";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.workon = pkgs.buildGoModule {
          pname = "workon";
          version = "0.2.2";
          src = ./.;
          vendorSha256 = "sha256-+EFL3cry1hFqVSWxobU6+V/30jbejft8kM5RXgroTxM=";
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
          buildInputs = with pkgs; [ go gopls ];
        };
      }
    );
}
