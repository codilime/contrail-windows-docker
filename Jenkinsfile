#!/usr/bin/env groovy



def powershell(script, returnStatus) { 
  pshell = "powershell.exe"
  base64script = script.bytes.encodeBase64().toString()
  bat(script: "${pshell} -NoProfile  -NonInteractive -NoLogo -ExecutionPolicy Bypass -EncodedCommand ${base64script}", returnStatus: $returnStatus)
}

def gotool(tool, args) {
  
}

