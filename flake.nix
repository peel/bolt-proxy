{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.05";
    devenv.url = "github:cachix/devenv";
    flake-utils.url = "github:numtide/flake-utils";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = {
    self,
    nixpkgs,
    devenv,
    flake-utils,
    ...
  } @ inputs: let
    pname = "bolt-proxy";
  in
    nixpkgs.lib.attrsets.recursiveUpdate
    (flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
          config.allowUnsupportedSystem = true;
        };
      in rec {
        packages = {
          devenv-up = devShell.config.procfileScript;
          default = pkgs.buildGoModule {
            inherit pname;
            version = self.shortRev or "${self.lastModifiedDate}-dirty";
            src = self;
            vendorHash = "sha256-SAN7JGtENRK8B0HnaV9dhY3qy+z1CS8UABi8t9pSCjM=";
            doCheck = false;
          };
          docker = pkgs.dockerTools.buildImage {
            name = "peelsky/${pname}";
            created = "now";
            tag = "latest";
            copyToRoot = self.packages."${system}".default;
            config.Entrypoint = ["/bin/${pname}"];
            config.Cmd = ["-debug"];
          };
        };
        devShell = devenv.lib.mkShell {
          inherit inputs pkgs;
          modules = [
            rec {
              scripts.docker-build.exec = ''
                nix build .#docker
                docker load < ./result
              '';
              scripts.go-watch.exec = ''
                find . -type f -name "*.go" -o -name "*.hcl" 2> /dev/null | ${pkgs.entr}/bin/entr -r $@
              '';
              scripts.go-test.exec = ''
                go test -v ./... $@
              '';
              scripts.go-run.exec = ''
                go run . $@
              '';
              packages = [pkgs.alejandra];
              difftastic.enable = true;
              languages.go.enable = true;
              enterShell = let
                names = pkgs.lib.mapAttrsToList (name: _: name) scripts;
              in ''
                echo "Available commands:"
                echo ""
                echo "devenv"
                ${
                  pkgs.lib.strings.concatMapStrings (script: "echo ${script}\n") names
                }
                echo ""
              '';
            }
          ];
        };
      }
    ))
    {
      packages.aarch64-darwin.docker = nixpkgs.legacyPackages.aarch64-darwin.dockerTools.buildImage {
        name = "peelsky/${pname}";
        created = "now";
        tag = "latest";
        copyToRoot = self.packages."aarch64-darwin".default.overrideAttrs (old:
          old
          // {
            GOOS = "linux";
            GOARCH = "amd64";
            CGO_ENABLED = 0;
          });
        config.Entrypoint = ["/bin/linux_amd64/${pname}"];
        config.Cmd = ["serve"];
      };
    };
}
