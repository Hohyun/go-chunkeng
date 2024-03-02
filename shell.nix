let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-23.11";
  pkgs = import nixpkgs { config = {}; overlays = []; };
in

pkgs.mkShellNoCC {
  packages = with pkgs; [
    nodejs_20
    go
    cowsay
    lolcat
  ];

  GREETING = "Hello, Chunk Go Server Project! (go_1.21, nodejs_20)";

  shellHook = ''
    echo $GREETING | cowsay | lolcat
  '';
}
