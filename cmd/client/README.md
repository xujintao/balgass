## 去壳方式得到的main1.exe

#### main_dump
x64dbg载入main，使用Scylla插件找到OEP，然后dump
```
0x00DF,478C: 55             ;OEP, dump, 那么dump的原理是什么？和打补丁有什么区别？
```
#### main_dump_SCY
```
一方面，因为壳的导入表把IAT指定到.rsrc段的0x099E,917C ~ 0x099E,9B3F位置  
另一方面，壳的导入表在kernel32.dll中多加了一个API叫VirtualQueryEx  

修复步骤  
在main基础上使用x64dbg自带的Scylla插件对main_dump的IAT进行修复(Fix Dump)  
恢复的本质是修复导入表把IAT指定到.rdata段的0x00de,9000 ~ 0x00de,99BF位置  
但是Scylla把符号表也复制了一份放在SCY段，有这个必要吗？
然后Scylla一致性很差，总是会在IAT尾部多分析出一些乱七八糟的import，虽然不影响使用但给人一种不靠谱的感觉
```

#### main1.exe
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

```
// with heapInit and flsInit
0B2BE980 | 60                       | pushad                                                        |
0B2BE981 | 9C                       | pushfd                                                        |
0B2BE982 | 68 70E92B0B              | push main.exe.B2BE970                                         | B2BE970:"main.dll"
0B2BE987 | FF15 D0921401            | call dword ptr ds:[<&LoadLibraryA>]                           | 0x011492D0
0B2BE98D | 85C0                     | test eax,eax                                                  |
0B2BE98F | 74 07                    | je main.exe.B2BE998                                           |
0B2BE991 | 9D                       | popfd                                                         |
0B2BE992 | 61                       | popad                                                         |
0B2BE993 | E9 F45DB3F5              | jmp main.exe.DF478C                                           |
0B2BE998 | FF15 0C921401            | call dword ptr ds:[<&GetLastError>]                           | 0x0114920C
0B2BE99E | 50                       | push eax                                                      |
0B2BE99F | E8 1C000000              | call main.exe.B2BE9C0                                         |
0B2BE9A4 | 58                       | pop eax                                                       |
0B2BE9A5 | 6A 00                    | push 0                                                        |
0B2BE9A7 | FF15 00921401            | call dword ptr ds:[<&ExitProcess>]                            | 0x01149200
0B2BE9AD | C3                       | ret                                                           |

0B2BE9C0 | 8BFF                     | mov edi,edi                                                   |
0B2BE9C2 | 55                       | push ebp                                                      |
0B2BE9C3 | 8BEC                     | mov ebp,esp                                                   |
0B2BE9C5 | 83EC 20                  | sub esp,20                                                    |
0B2BE9C8 | 6A 01                    | push 1                                                        | heapInit
0B2BE9CA | E8 A506B4F5              | call main.exe.DFF074                                          |
0B2BE9CF | 59                       | pop ecx                                                       |
0B2BE9D0 | E8 D7DBB3F5              | call main.exe.DFC5AC                                          | flsInit
0B2BE9D5 | FF75 08                  | push dword ptr ss:[ebp+8]                                     |
0B2BE9D8 | 68 B0E92B0B              | push main.exe.B2BE9B0                                         | B2BE9B0:"getlasterror %d"
0B2BE9DD | 8D7D E0                  | lea edi,dword ptr ss:[ebp-20]                                 |
0B2BE9E0 | 57                       | push edi                                                      |
0B2BE9E1 | E8 9497B2F5              | call main.exe.DE817A                                          |
0B2BE9E6 | 83C4 0C                  | add esp,C                                                     |
0B2BE9E9 | 6A 00                    | push 0                                                        |
0B2BE9EB | 6A 00                    | push 0                                                        |
0B2BE9ED | 8D75 E0                  | lea esi,dword ptr ss:[ebp-20]                                 |
0B2BE9F0 | 56                       | push esi                                                      |
0B2BE9F1 | 6A 00                    | push 0                                                        |
0B2BE9F3 | FF15 3C981401            | call dword ptr ds:[<&MessageBoxA>]                            | 0x0114983C
0B2BE9F9 | C9                       | leave                                                         |
0B2BE9FA | C3                       | ret                                                           |

// jmp 0x0B2BE980
0B2BE911  6C A0 4C 0A 8B FF 55 8B EC 53 56 8B 75 08 8B 46  l L..ÿU.ìSV.u..F  

// origin
0B2BE970  6D 61 69 6E 2E 64 6C 6C 00 00 00 00 00 00 00 00  main.dll........  
0B2BE980  60 9C 68 70 E9 2B 0B FF 15 D0 92 14 01 85 C0 74  `.hpé+.ÿ.Ð....Àt  
0B2BE990  07 9D 61 E9 F4 5D B3 F5 FF 15 0C 92 14 01 50 E8  ..aéô]³õÿ.....Pè  
0B2BE9A0  1C 00 00 00 58 6A 00 FF 15 00 92 14 01 C3 00 00  ....Xj.ÿ.....Ã..  
0B2BE9B0  67 65 74 6C 61 73 74 65 72 72 6F 72 20 25 64 00  getlasterror %d.  
0B2BE9C0  8B FF 55 8B EC 83 EC 20 6A 01 E8 A5 06 B4 F5 59  .ÿU.ì.ì j.è¥.´õY  
0B2BE9D0  E8 D7 DB B3 F5 FF 75 08 68 B0 E9 2B 0B 8D 7D E0  è×Û³õÿu.h°é+..}à  
0B2BE9E0  57 E8 94 97 B2 F5 83 C4 0C 6A 00 6A 00 8D 75 E0  Wè..²õ.Ä.j.j..uà  
0B2BE9F0  56 6A 00 FF 15 3C 98 14 01 C9 C3 00 00 00 00 00  Vj.ÿ.<...ÉÃ.....  

// stream
6D 61 69 6E 2E 64 6C 6C 00 00 00 00 00 00 00 00 60 9C 68 70 E9 2B 0B FF 15 D0 92 14 01 85 C0 74 07 9D 61 E9 F4 5D B3 F5 FF 15 0C 92 14 01 50 E8 1C 00 00 00 58 6A 00 FF 15 00 92 14 01 C3 00 00 67 65 74 6C 61 73 74 65 72 72 6F 72 20 25 64 00 8B FF 55 8B EC 83 EC 20 6A 01 E8 A5 06 B4 F5 59 E8 D7 DB B3 F5 FF 75 08 68 B0 E9 2B 0B 8D 7D E0 57 E8 94 97 B2 F5 83 C4 0C 6A 00 6A 00 8D 75 E0 56 6A 00 FF 15 3C 98 14 01 C9 C3 00 00 00 00 00
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

## breakpoints and logs
#### do13: 0x004CCD8F
```
// call dword ptr ds:[eax+30]
conditional breakpoints: 0
conditional tracing: do13: {ecx}.f{dword(eax+30)}
```

#### igc.dll, giocp.cpp:129
```
// movzx eax,byte ptr ds:[edx+ecx]
conditional breakpoints: word(edx+2)==31f3 || word(edx+2)==31f3
conditional tracing:     prefix: {mem;4@edx}
```

#### igc.dll dllmain.cpp:409
```
// sub eax,dword ptr ds:[unsigned long dwAntiHackTime]
conditional breakpoints: 0
conditional tracing: tick: (x){eax}ms, (u){u:eax}ms, (u){u:eax/3E8}s
```

#### net.go 0x004397E9
```
// mov dword ptr ss:[ebp-10],ecx
conditional breakpoints: word(dword(ebp+8))==0bc3 || word(dword(ebp+8))==13c3
conditional tracing: net log: {mem;3@dword(ebp+8)}
```

#### api.send, for some inline scene
```
// push esi
// word(dword(ebp+c))==0bc3 || word(dword(ebp+c))==13c3
```

## GFX edit
#### MainFrame: show sd/ag and position tweak
```as3
// scripts/__Packages/MainFrame
class MainFrame extends Common.Base.MUComponent
{
    // line 10:
    function onLoadUI()
    {
        // line 45:
        this.tfSD = this.mcMainFrame.tfSD;
        this.SetSD(0,0); // this.tfSD._visible = false;
        this.tfAG = this.mcMainFrame.tfAG;
        this.SetAG(0,0); // this.tfAG._visible = false;
    }

    // line 113:
    function onAddEvent()
    {
        // line 132:
        // this.btnClearSD.addEventListener("rollOver",this,"onOverSD");
        // this.btnClearSD.addEventListener("rollOut",this,"onOutSD");
        // this.btnClearAG.addEventListener("rollOver",this,"onOverAG");
        // this.btnClearAG.addEventListener("rollOut",this,"onOutAG");
    }
}

// tweak
DefineSprite(209).PlaceObject3(206).matrix.translateY = 570 // 600
DefineSprite(209).PlaceObject3(204).matrix.translateY = 570 // 600
```

#### PartyFrame: show party hp/mp value and name position tweak
```
// hp
DefineEditText(72).bounds.Ymax = 240 // 80
DefineEditText(72).hasText = true // false
DefineEditText(72).fontHeight = 200 // 20
DefineEditText(72).textColor = [216,216,216,255] // [171,171,171,255]
DefineEditText(72).align = 2 // 0
DefineSprite(74).PlaceObject2(72).Depth = 3 // 1
DefineSprite(74).PlaceObject2(72).matrix.nTranslateBits = 8 // 7
DefineSprite(74).PlaceObject2(72).matrix.translateX = -100 // 40
DefineSprite(74).PlaceObject2(72).matrix.translateY = -100 // 40
DefineSprite(74).PlaceObject2(72).placeFlagHasName = true // false
DefineSprite(74).PlaceObject2(72).name = tfHP

// mp
DefineEditText(69).bounds.Ymax = 240 // 80
DefineEditText(69).fontHeight = 200 // 20
DefineEditText(69).textColor = [224,224,224,255] // [171,171,171,255]
DefineEditText(69).align = 2 // 0
DefineSprite(71).PlaceObject2(69).Depth = 3 // 1
DefineSprite(71).PlaceObject2(69).matrix.nTranslateBits = 8 // 7
DefineSprite(71).PlaceObject2(69).matrix.translateX = -100 // 40
DefineSprite(71).PlaceObject2(69).matrix.translateY = -60 // 40
DefineSprite(71).PlaceObject2(69).placeFlagHasName = true // false
DefineSprite(71).PlaceObject2(69).name = tfMP

// scripts/__Packages/PartyFrame
class PartyFrame extends Common.Base.MUComponent
{
    function ClearPartyInfo(iIndex)
    {
        this.amcPartyMember[iIndex]._visible = false;
        this.atfChannel[iIndex].text = "";
        this.atfName[iIndex].text = "";
        this.abtnLeave[iIndex]._visible = false;
        this.amcCrown[iIndex].gotoAndStop(Common.Global.MUDefines.MATCHING_PARTY_NORMAL_MEMBER_FRAME);
        this.amcStateLayer[iIndex].gotoAndStop(1);
        this.apbHPBar[iIndex].minimum = 0;
        this.apbHPBar[iIndex].maximum = 0;
        this.apbHPBar[iIndex].position = 0;
        this.apbHPBar[iIndex].tfHP.text = ""; // add
        this.apbMpBar[iIndex].minimum = 0;
        this.apbMpBar[iIndex].maximum = 0;
        this.apbMpBar[iIndex].position = 0;
        this.apbMpBar[iIndex].tfMP.text = ""; // add
        this.RemoveAllBuff(iIndex);
    }
    function SetPartyHPValue(iIndex, iHP, iHPMax)
    {
        this.apbHPBar[iIndex].tfHP.text = String(iHP) + " / " + String(iHPMax); // add
    }
    function SetPartyMPValue(iIndex, iMP, iMPMax)
    {
        this.apbMpBar[iIndex].tfMP.text = String(iMP) + " / " + String(iMPMax); // add
    }
}

// tweak
DefineEditText(63).bounds.Ymax = 320 // 280
DefineEditText(64).bounds.Ymax = 320 // 280
DefineSprite(89).frame1.PlaceObject2(65).matrix.translateY = 60 // 90
```

#### PetFrame: name position tweak
```
DefineEditText(20).bounds.Ymax = 320 // 280
DefineEditText(21).bounds.Ymax = 320 // 280
DefineSprite(27).frame1.PlaceObject2(22).matrix.translateY = 60 // 90
```

#### reference
```
http://forum.ragezone.com/f82/convert-ozp-files-png-1074563/
https://forum.xentax.com/viewtopic.php?t=12157
http://forum.ragezone.com/f508/development-decrypt-files-ozd-ozg-1124869/
http://forum.ragezone.com/f197/release-igcn-season-9-5-a-1125721/index58.html
https://en.wikipedia.org/wiki/Scaleform_GFx
https://gtamods.com/wiki/.gfx
https://github.com/jindrapetrik/jpexs-decompiler
https://www.adobe.com/content/dam/acom/en/devnet/pdf/swf-file-format-spec.pdf
```