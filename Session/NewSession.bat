if "%1"=="" then exit

echo creating new session: session\%1\

mkdir session\%1
mkdir session\%1\design
mkdir session\%1\archives

copy design\Examples\SessionCurrent.bpd session\%1\design
