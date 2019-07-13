## main.exe

##### main1
x64dbg
```
0x00DF,478C: 55             ;OEP
```

##### main2
search mu.exe
```
0x004D,7E2A: 73 6D -> EB 6D ;jae to jmp, disable mu.exe
```

##### main31
search config.ini read error
```
0x004D,8102: 74 19 -> EB 19 ;je to jmp, disable GameGuard
```
##### main32
search gg init error
```
0x004D,8145: 75 7C -> EB 7C ;jne to jmp, disable GameGuard
```

##### main33
search ResourceGuard
```
0x004D,D639: 74 23 -> EB 23 ;je to jmp, disable GameGuard
```
##### main4
notepad++ Hex-editor plugin

##### main5
instead obfuscation with normal code
```
0x004D,7CE5: -> push ebp
0x004D,7CE6: -> mov ebp,esp
0x004D,7CE8: -> mov eax,1B60
```

##### main6
一方面，因为壳的导入表把IAT指定到.rsrc段的0x099E,917C ~ 0x099E,9B3F位置  
另一方面，壳的导入表在kernel32.dll中多加了一个API叫VirtualQueryEx  

修复步骤  
在main基础上使用x64dbg自带的Scylla插件对main5的IAT进行修复(Fix Dump)  
恢复的本质是修复导入表把IAT指定到.rdata段的0x00de,9000 ~ 0x00de,99BF位置  
但是Scylla把符号表也复制了一份放在SCY段，有这个必要吗？
然后Scylla一致性很差，总是会在IAT尾部多分析出一些乱七八糟的import，虽然不影响使用但给人一种不靠谱的感觉

##### main7
r6602 floating point support not loaded
```
0x00E1,A4AC: push 2 -> ret
让它直接ret，这样靠谱吗？
```

##### main8
disable log encode
```
0x00B3,8653: push ebp -> ret
```