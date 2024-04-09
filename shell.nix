let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-23.11";
  pkgs = import nixpkgs { config = {}; overlays = []; };
in

pkgs.mkShellNoCC {
  packages = with pkgs; [
    go
    air
    cowsay
    lolcat
  ];

  GREETING = "Welcome to go(1.21), develpment environment powered by Nix!";

  shellHook = ''
    echo $GREETING | cowsay | lolcat
  '';
}
