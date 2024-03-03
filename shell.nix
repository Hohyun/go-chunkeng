let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-23.11";
  pkgs = import nixpkgs { config = {}; overlays = []; };
in

pkgs.mkShellNoCC {
  packages = with pkgs; [
    redis
    nodejs_20
    go
    air
    cowsay
    lolcat
  ];

  GREETING = "Welcome to go(1.21), nodejs(20), redis develpment environment powered by Nix!";

  shellHook = ''
    echo $GREETING | cowsay | lolcat
  '';
}
