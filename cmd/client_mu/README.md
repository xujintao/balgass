## inject mu.dll
```
// 0x0043503B: jmp to load mu.dll
0043503B | E9 A0740100              | jmp mu.44C4E0                                 |

// 0x0044C4E0: load mu.dll
0044C4E0 | 60                       | pushad                                        | load mu.dll
0044C4E1 | 9C                       | pushfd                                        |
0044C4E2 | 68 D0C44400              | push mu.44C4D0                                | 44C4D0:"mu.dll"
0044C4E7 | FF15 88D24400            | call dword ptr ds:[<&LoadLibraryA>]           |
0044C4ED | 85C0                     | test eax,eax                                  |
0044C4EF | 74 07                    | je mu.44C4F8                                  |
0044C4F1 | 9D                       | popfd                                         |
0044C4F2 | 61                       | popad                                         |
0044C4F3 | E9 C089FEFF              | jmp mu.434EB8                                 |
0044C4F8 | FF15 18D34400            | call dword ptr ds:[<&GetLastError>]           |
0044C4FE | 50                       | push eax                                      |
0044C4FF | E8 1C000000              | call mu.44C520                                |
0044C504 | 58                       | pop eax                                       |
0044C505 | 6A 00                    | push 0                                        |
0044C507 | FF15 9CD14400            | call dword ptr ds:[<&ExitProcess>]            |
0044C50D | C3                       | ret                                           |

// f0044C520printErr
0044C520 | 8BFF                     | mov edi,edi                                   |
0044C522 | 55                       | push ebp                                      |
0044C523 | 8BEC                     | mov ebp,esp                                   |
0044C525 | 83EC 20                  | sub esp,20                                    |
0044C528 | 6A 01                    | push 1                                        |
0044C52A | E8 3A5EFFFF              | call mu.442369                                |
0044C52F | 59                       | pop ecx                                       |
0044C530 | E8 DA00FFFF              | call mu.43C60F                                |
0044C535 | FF75 08                  | push dword ptr ss:[ebp+8]                     |
0044C538 | 68 10C54400              | push mu.44C510                                | 44C510:"getlasterror %d"
0044C53D | 8D7D E0                  | lea edi,dword ptr ss:[ebp-20]                 |
0044C540 | 57                       | push edi                                      |
0044C541 | E8 756AFEFF              | call mu.432FBB                                |
0044C546 | 83C4 0C                  | add esp,C                                     |
0044C549 | 6A 00                    | push 0                                        |
0044C54B | 6A 00                    | push 0                                        |
0044C54D | 8D75 E0                  | lea esi,dword ptr ss:[ebp-20]                 |
0044C550 | 56                       | push esi                                      |
0044C551 | 6A 00                    | push 0                                        |
0044C553 | FF15 24D54400            | call dword ptr ds:[<&MessageBoxA>]            |
0044C559 | C9                       | leave                                         |
0044C55A | C3                       | ret                                           |

// origin
0044C4D0  6D 75 2E 64 6C 6C 00 00 00 00 00 00 00 00 00 00  mu.dll..........  
0044C4E0  60 9C 68 D0 C4 44 00 FF 15 88 D2 44 00 85 C0 74  `.hÐÄD.ÿ..ÒD..Àt  
0044C4F0  07 9D 61 E9 C0 89 FE FF FF 15 18 D3 44 00 50 E8  ..aéÀ.þÿÿ..ÓD.Pè  
0044C500  1C 00 00 00 58 6A 00 FF 15 9C D1 44 00 C3 00 00  ....Xj.ÿ..ÑD.Ã..  
0044C510  67 65 74 6C 61 73 74 65 72 72 6F 72 20 25 64 00  getlasterror %d.  
0044C520  8B FF 55 8B EC 83 EC 20 6A 01 E8 3A 5E FF FF 59  .ÿU.ì.ì j.è:^ÿÿY  
0044C530  E8 DA 00 FF FF FF 75 08 68 10 C5 44 00 8D 7D E0  èÚ.ÿÿÿu.h.ÅD..}à  
0044C540  57 E8 75 6A FE FF 83 C4 0C 6A 00 6A 00 8D 75 E0  Wèòfþÿ.Ä.j.j..uà  
0044C550  56 6A 00 FF 15 24 D5 44 00 C9 C3 00 00 00 00 00  Vj.ÿ.$ÕD.ÉÃ.....  

// stream
6D 75 2E 64 6C 6C 00 00 00 00 00 00 00 00 00 00 60 9C 68 D0 C4 44 00 FF 15 88 D2 44 00 85 C0 74 07 9D 61 E9 C0 89 FE FF FF 15 18 D3 44 00 50 E8 1C 00 00 00 58 6A 00 FF 15 9C D1 44 00 C3 00 00 67 65 74 6C 61 73 74 65 72 72 6F 72 20 25 64 00 8B FF 55 8B EC 83 EC 20 6A 01 E8 3A 5E FF FF 59 E8 DA 00 FF FF FF 75 08 68 10 C5 44 00 8D 7D E0 57 E8 F2 66 FE FF 83 C4 0C 6A 00 6A 00 8D 75 E0 56 6A 00 FF 15 24 D5 44 00 C9 C3 00 00 00 00 00
```