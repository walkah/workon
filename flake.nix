{
  description = "Manage tmux for what you work on.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    {
      overlays.default = final: prev: {
        inherit (self.packages.${prev.system}) workon;
      };
    } // flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.workon = pkgs.buildGoModule {
          pname = "workon";
          version = "0.2.3";
          src = ./.;
          vendorHash = "sha256-s2HCpKQ718q6L17wnJWf9yHkrcB1LZ6D185XX8MNt1Q=";
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
