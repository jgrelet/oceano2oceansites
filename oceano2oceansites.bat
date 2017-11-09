pwd
SET "DRIVE=M:"
set "CRUISE=SARGASSE"
set "CONFIG=sargasse.toml"
set "PREFIX=S"
cd %DRIVE%\%CRUISE%\data-processing\CTD
pwd
%GOBIN%\oceano2oceansites.exe -c %DRIVE%\%CRUISE%\data-processing\%CONFIG% -r %DRIVE%\%CRUISE%\local\code_roscop.csv -e --files=%DRIVE%\%CRUISE%\data-processing\CTD\data\cnv\%PREFIX%*.cnv 
rem %DRIVE%\%CRUISE%\local\sbin\oceano2oceansites -c %DRIVE%\%CRUISE%\data-processing\%CONFIG% -r %DRIVE%\%CRUISE%\local\code_roscop.csv -e -a --files=%DRIVE%\%CRUISE%\data-processing\CTD\data\cnv\%PREFIX%*.cnv 
pause
