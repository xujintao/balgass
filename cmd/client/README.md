## 去壳方式得到的main1.exe

##### main_dump
x64dbg载入main，使用Scylla插件找到OEP，然后dump
```
0x00DF,478C: 55             ;OEP, dump, 那么dump的原理是什么？和打补丁有什么区别？
```
##### main_dump_SCY
```
一方面，因为壳的导入表把IAT指定到.rsrc段的0x099E,917C ~ 0x099E,9B3F位置  
另一方面，壳的导入表在kernel32.dll中多加了一个API叫VirtualQueryEx  

修复步骤  
在main基础上使用x64dbg自带的Scylla插件对main_dump的IAT进行修复(Fix Dump)  
恢复的本质是修复导入表把IAT指定到.rdata段的0x00de,9000 ~ 0x00de,99BF位置  
但是Scylla把符号表也复制了一份放在SCY段，有这个必要吗？
然后Scylla一致性很差，总是会在IAT尾部多分析出一些乱七八糟的import，虽然不影响使用但给人一种不靠谱的感觉
```

##### main1
~~notepad++ Hex-editor plugin~~
```
replace "connect.muchina.com" with "192.168.0.100"
```

载入dll
```
内存窗口匹配15个CC的空间，这样的坑大概有2500多个，我们使用0x0114开头的空间
0x01140B54坑 硬编码 "IGC.dll"
0x01141792坑 硬编码 LoadLibraryA
...
```
## 非去壳方式得到的main2.exe
```
1，利用壳加密没考虑页对齐的漏洞，直接在.rsrc尾部很顺利的打上载入dll的补丁
2，在0x0B2BE910地址(memcpy后的0x00DF490F)处hook住jmp 0x00DF478C来执行我们的补丁

这样很完美，相当于把dll的载入融入壳的代码，但每次调试都要手动stepover解密函数后再启用断点
```

## 以下所有hook都移动到IGC.dll中了  
search mu.exe
```
0x004D,7E2A: 73 6D -> EB 6D ;jae to jmp, disable mu.exe
```

search config.ini read error
```
0x004D,8102: 74 19 -> EB 19 ;je to jmp, disable GameGuard
```

search gg init error
```
0x004D,8145: 75 7C -> EB 7C ;jne to jmp, disable GameGuard
```

search ResourceGuard
```
0x004D,D639: 74 23 -> EB 23 ;je to jmp, disable GameGuard
```

disable log encode
```
0x00B3,8653: push ebp -> ret ;disable log encode
```

fix "r6602 floating point support not loaded"
```
0x00E1,A4AC: push 2 -> ret ;让它直接ret，这样靠谱吗？
```

~~instead obfuscation with normal code, 这个不行，还不能动壳的代码~~
```
0x004D,7CE5: -> push ebp
0x004D,7CE6: -> mov ebp,esp
0x004D,7CE8: -> mov eax,1B60
```

disable version match
```
0x0ABD,696B: je 0x006C76BA -> jmp 0x006C76BA
```

## build
```
cd ~/github.com/xujintao/balgass/cmd/client
GOOS=windows go build -x github.com/xujintao/balgass/cmd/client &> build.out
```

## breakpoints
```
// f00439420 -> SendPacket
// C1
conditional breakpoints: word(edx+2)==31f3 || word(edx+2)==31f3
conditional tracing:     prefix: {mem;4@edx}
```

```
// igc.dll dllmain.cpp:409
conditional tracing: tick: (x){eax}ms, (u){u:eax}ms, (u){u:eax/3E8}s
```
