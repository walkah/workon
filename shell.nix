let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs { };
in pkgs.mkShell {
  name = "workon";
  buildInputs = with pkgs; [ go ];
}
