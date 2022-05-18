{
  description = "Manage tmux for what you work on.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/release-21.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.workon = pkgs.buildGoModule {
          name = "workon";
          src = self;
          vendorSha256 = "sha256-Jao45EKn3rybLa8Hi6xvE8ByPESz7Tvx4b8JehTCWww=";
        };

        defaultPackage = self.packages.${system}.workon;

        devShell = pkgs.mkShell {
          name = "workon";
          buildInputs = with pkgs; [ go gopls ];
        };
      }
    );
}
