pwd
SET "DRIVE=data/CTD/test"
set "CRUISE=TEST"
set "CONFIG=oceano2oceansites.toml"
set "ROSCOP=roscop/code_roscop.csv"
set "PREFIX=csp"
pwd
%GOBIN%/oceano2oceansites -c %CONFIG% -r %ROSCOP% -e --files=%DRIVE%/%PREFIX%*.cnv 
pause
