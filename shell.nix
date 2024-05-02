{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  name = "go-dev";

  buildInputs = with pkgs; [
    go
    gopls
  ];
}
