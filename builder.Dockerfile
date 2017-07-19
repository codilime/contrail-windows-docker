# escape=`

FROM microsoft/windowsservercore

SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProressPreference = 'SilentlyContinue';"]

ENV DOWNLOADS_PATH C:\downloads

ENV GIT_VERSION 2.13.3

RUN New-Item -Type Directory $Env:DOWNLOADS_PATH | Out-Null;

# Install git
RUN $gitDownloadPath = 'https://www.nuget.org/api/v2/package/GitForWindows/' + $Env:GIT_VERSION; `
    (New-Object System.Net.WebClient).DownloadFile($gitDownloadPath, $Env:DOWNLOADS_PATH + '\git.zip');
RUN $gitSetupPath = $Env:DOWNLOADS_PATH + '\git.zip'; `
    Expand-Archive $gitSetupPath C:\git-tmp; `
    New-Item -Type Directory C:\git | Out-Null; `
    Move-Item C:\git-tmp\tools\* C:\git\.; `
    Remove-Item -Recurse -Force C:\git-tmp;
RUN setx /M PATH $('C:\git\cmd;C:\git\usr\bin;'+$Env:PATH) | Out-Null;

# Install GCC + other stuff
RUN $gccDownloadPath = 'https://raw.githubusercontent.com/jhowardmsft/docker-tdmgcc/master/gcc.zip'; `
    (New-Object System.Net.WebClient).DownloadFile($gccDownloadPath, $Env:DOWNLOADS_PATH + '\gcc.zip');
RUN $runtimeDownloadPath = 'https://raw.githubusercontent.com/jhowardmsft/docker-tdmgcc/master/runtime.zip'; `
    (New-Object System.Net.WebClient).DownloadFile($runtimeDownloadPath, $Env:DOWNLOADS_PATH + '\runtime.zip');
RUN $binutilsDownloadPath = 'https://raw.githubusercontent.com/jhowardmsft/docker-tdmgcc/master/binutils.zip'; `
    (New-Object System.Net.WebClient).DownloadFile($binutilsDownloadPath, $Env:DOWNLOADS_PATH + '\binutils.zip');
RUN Expand-Archive -Path $($Env:DOWNLOADS_PATH + '\gcc.zip') -Destination C:\gcc -Force;
RUN Expand-Archive -Path $($Env:DOWNLOADS_PATH + '\runtime.zip') -Destination C:\gcc -Force;
RUN Expand-Archive -Path $($Env:DOWNLOADS_PATH + '\binutils.zip') -Destination C:\gcc -Force;
RUN setx /M PATH $($Env:PATH+'C:\gcc\bin;') | Out-Null;

# Install Go
# According to: http://www.wadewegner.com/2014/12/easy-go-programming-setup-for-windows/
ENV GO_VERSION 1.8.3
RUN $goDownloadPath = 'https://golang.org/dl/go' + $Env:GO_VERSION + '.windows-amd64.zip'; `
    (New-Object System.Net.WebClient).DownloadFile($goDownloadPath, $Env:DOWNLOADS_PATH + '\go.zip');
RUN $goSetupPath = $Env:DOWNLOADS_PATH + '\go.zip'; `
    Expand-Archive -Path $goSetupPath -Destination C:\;
RUN setx /M PATH $($Env:PATH+'C:\go\bin;') | Out-Null;

ENV GOROOT C:\go
RUN New-Item -Type Directory $Env:GOPATH | Out-Null;

ENTRYPOINT ["powershell.exe"]