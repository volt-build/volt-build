{ stdenv, fetchFromGitHub, pkgs}: 
stdenv.mkDerivation {
  pname = "mini-build" ;
  version = "0.1.0";
  src = fetchFromGitHub {
    owner = "randomdude16671"; 
    repo = "mini-build"; 
    rev = "main"; 
    sha256 = "119sm4hq9mjms7r0ca3y2c0klh4rrv6kypmpcfvh4v6cwhsb7lhm"; 
  };

  buildInputs = [
    pkgs.go 
  ]; 
  
  buildPhase = '' 
    go build -o $out/target/ . 
  '';
  installPhase = '' 
    cp $out/target/* $out/bin/ 
  '';

  meta = {
    description = "A super small build system I wrote myself"; 
    homepage = "https://github.com/randomdude16671/mini-build"; 
    license = stdenv.lib.licenses.mit; 
    maintainers = with stdenv.lib.maintainers; []; 
    platforms = stdenv.lib.platform.all; 
  };
}
