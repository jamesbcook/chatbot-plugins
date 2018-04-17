# VirusTotal Plugin

## Details

* This plugin query the VirusTotal api and return what signatures are available for that hash.
  * ```export CHATBOT_VIRUSTOTAL={APIKEY}```

```
CMD: /virustotal
Help: /virustotal {sha256 of file}
```

### Example

```
/virustotal bd2c2cf0631d881ed382817afcce2b093f4e412ffb170a719e2762f250abfea4
-------------------------
VirusTotal Detection Results
Total Detected 6
Kaspersky  not-a-virus:HEUR:RiskTool.Win32.ProcHack.gen
ZoneAlarm  not-a-virus:HEUR:RiskTool.Win32.ProcHack.gen
ALYac      Misc.Riskware.ProcessHacker
Jiangmin   RiskTool.ProcHack.b
Panda      HackingTool/ProcHack
CAT-QuickHeal HackTool.ProcHacker.S906087
```
