## main.exe

###### main1
x64dbg
```
0x00DF,478C: 55             ;OEP
```

###### main2
search mu.exe
```
0x004D,7E2A: 73 6D -> EB 6D ;disable mu.exe
```

###### main31
search config.ini read error
```
0x004D,8102: 74 19 -> EB 19 ;disable GameGuard
```
###### main32
search gg init error
```
0x004D,8145: 75 7C -> EB 7C ;disable GameGuard
```

###### main33
search ResourceGuard
```
0x004D,D639: 74 23 -> EB 23 ;disable GameGuard
```
###### main4
notepad++ Hex-editor plugin

###### main5
```
0x004D,7CE5: -> push ebp
0x004D,7CE6: -> mov ebp,esp
0x004D,7CE8: -> mov eax,1B60
```