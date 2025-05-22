{ lib, pkgs }: 
pkgs.buildGoModule {
    pname = "mini-build";
    version = "0.1.0";

    src = pkgs.fetchFromGitHub {
      owner = "randomdude16671";
      repo = "mini-build";
      rev = "main";
      sha256 = "MK+vzffho2YprwuBiGP9vF2a6GaGQYCPTsu0u6HBgec=";
    };
    vendorHash = null;
    meta = with lib; {
      description = "A super small build system written mainly for personal use.";
      homepage = "https://github.com/randomdude16671/mini-build";
      platforms = platforms.linux;
      license = licenses.mit;
    };
  }
