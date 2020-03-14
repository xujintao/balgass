package main

import (
	"fmt"
	"os"
)

// dialogEx base dialog
type dialogEx struct{}

// dlgcore.cpp
func f004158B6mfcAfxDlgProc(hWnd uintptr, msg uint, wParam uint, lParam uint) uintptr {
	switch msg {
	case 0x110: // WM_INITDIALOG
		// CDialog* pDlg = DYNAMIC_DOWNCAST(CDialog, CWnd::FromHandlePermanent(hWnd));
		// if (pDlg != NULL)
		// 	return pDlg->OnInitDialog(); // f0040CA40OnInitDialog
		// else
		// 	return 1;
	}
	return 0
}

// wincore.cpp
func f0041B33BmfcAfxWndProc(hWnd uintptr, msg uint, wParam, lParam uint) uintptr {
	// if msg == WM_QUERYAFXWNDPROC {
	// 	return 1
	// }
	// CWnd* pWnd = CWnd::FromHandlePermanent(hWnd);
	// if (pWnd == NULL || pWnd->m_hWnd != hWnd)
	// 	return ::DefWindowProc(hWnd, msg, wParam, lParam)
	// return AfxCallWndProc(pWnd, hWnd, msg, wParam, lParam) // -> Cwnd::WindowProc -> Cwnd::OnWndMsg
	return 0
}

// wincore.cpp
func f0041C495mfcAfxCbtFilterHook() {
	// ...
	// dll.user32.SetWindowLongA(f0041B33BmfcAfxWndProc)
}

// f00416045, dlgcore.cpp, CDialog::DoModal
func (d *dialogEx) DoModal() uintptr {
	// Load(0x00400000, FindResourceA(0x004000000, 0x66, 0x05))
	// d.f00415BBBpreModal() // CDialog::PreModal()
	func() {
		// f0041C6E8mfcAfxHookWindowCreate(Cwnd* pWnd)
		func() {
			// dll.user32.SetWindowsHookExA(5, f0041C495mfcAfxCbtFilterHook, 0, dll.user32.GetCurrentThreadId())
		}() // this
	}()

	// ...

	// inline: Create and run dialog indirect
	// d.f00415E8FcreateDlgIndirect()
	func() {
		// ...
		// f00425FA7()
		// f0041C6E8mfcAfxHookWindowCreate()
		// dll.user32.CreateDialogIndirectParamA(0x00400000, &temp, 0, f004158B6mfcAfxDlgProc, 0) // 回调Proc
		// f00401A30()
		// f0041AC7D()
	}()
	// d.f004163F3GetStyle()
	// d.f0041A74DRunModalLoop(4) // MLF_SHOWNIDLE
	func(x int) {
		// msg := struct{}{}
		// for dll.user32.PeekMessageA(&msg, 0, 0, 0, 0) == true {
		// 	// f00420F56()
		// 	func() {
		// 		// f0041786A()
		// 		msg := f0041726B()
		// 		if dll.user32.GetMessageA(&msg, 0, 0, 0) {
		// 			if f00420D41handle(&msg) == false {
		// 				dll.user32.TranslateMessage(&msg)
		// 				dll.user32.DispatchMessage(&msg)
		// 			}
		// 		}
		// 	}()
		// }
	}(4)
	return 0
}

func (d *dialogEx) OnInitDialog() bool {
	return true
}

var v004505F8 = "CDialog"
var v00461DF0stack int = 0xC58267AF
var v0046327C int = 0
var v00463280IP string
var v004632B4port uint16
var v0048F4A0 func()

type t0018FC8C struct {
	m108 *os.File // 0x00462010
}

func (t *t0018FC8C) f00403F50open(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return f
}

func (t *t0018FC8C) f00403FC0fread(f *os.File, data []byte) int {
	n, err := f.Read(data)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return n
}

func (t *t0018FC8C) f00403FB0EOF() bool {
	// f004331E4feof(t.m108)
	return false // true: eof
}

func (t *t0018FC8C) f00403F80close() {
	t.m108.Close()
}

type t0018FDB0 struct {
	// partitions
	partitions        []*partition
	m0CpartitionsHead *partition // head
	m10partitionsTail *partition // tail
	m14partitionsSafe *partition // safe
}

func (t *t0018FDB0) f0040ADF0(p *partition) {
	t.partitions = append(t.partitions, p)
}

type partition struct {
	m00Index int
	m04Num   int
	m08Site  string
	m3AName  string
	m7AIP    string
	mBCPort  int
}

func (p *partition) f00405E00() {
	// partition section initiation
}

type button struct {
}

func (b *button) f0040F770() {
	// b.f0041987A()
}

type block struct {
	m00 uint
	m04 uint
	m08 uint
	m0C uint
	m10 uint
	m14 uint
	m18 uint
	m1C uint
	m20 uint
	m24 uint8
	m28 uint
	m2C uint
	m30 uint
	m34 uint
	m38 uint
	m3C uint
	m40 uint
	m44 uint
	m48 uint
	m4C uint
	m50 uint
	m54 uint
	m58 uint
	m68 uint
	m74 uint
	m78 uint
	m84 uint
	m88 uint
}

func (b *block) f00417837() uint {
	// v00A6DE08 := v0048F8C8.f0042566A(f0041547C)
	// if v00A6DE08 == 0 {
	// 	// 0x00415460
	// }
	// if v00A6DE08.f04 != 0 {
	// 	return 1
	// }
	var ret uint
	// ret = v0048F8C4.f00425146(f00417808)
	// if ret == 0 {
	// 	// 0x00415460
	// }
	return ret // 0x00A6DF18
}

func (b *block) f0041987A() {
	// b.f0041DE41()
	func() {
		// b.f00417837
		b.m1C = b.f00417837()
		b.m04 = 1
		b.m08 = 0
		b.m0C = 0
		b.m10 = 0
		b.m14 = 1
		b.m18 = 0
	}()
	b.m00 = 0x004510AC
	b.m20 = 0
	b.m24 = 0
	b.m28 = 0
	b.m2C = 0
	b.m30 = 0x0045101C
	b.m34 = 0x00451090
	b.m38 = 0
	b.m3C = 0
	b.m40 = 0
	b.m44 = 0
	b.m48 = 0
	b.m4C = 0
	b.m50 = 0
}

func (b *block) f0040DA60() {
	b.f0041987A()
	b.m00 = 0x0044ECF4
}

func (b *block) f0040F670() {
	b.m00 = 0x0044F568
}

// muDlg dialog
type muDlg struct {
	dialogEx
	blocks [10]block // 0x00 0x98 0xEC
	// ...
	pathName       [256]uint8 // 0x3F96C
	major          uint8      // 0x3FAFC
	minor          uint8      // 0x3FAFD
	patch          uint8      // 0x3FAFE
	name           [256]uint8 // 0x44928, "奇迹("
	dir            [260]uint8 // 0x44A28
	block2         [100]block // 0x44B40 0x44B94 0x44BEC 0x44C60 0x44CD4
	lParam         uint       // 0x44FE4
	m44FE8partiSum int        // 0x44FE8
	m44FEC         struct{}   // 0x44FEC, 0x0018FC74
	m45004         t0018FC8C  // 0x45004, 0x0018FC8C
	m45110         struct{}   // 0x45110, 0x0018FD98
	m45128         t0018FDB0  // 0x45128, 0x0018FDB0
}

func (d *muDlg) f0040AEC0(x uint) {
	// d.blocks[0].f0040E510(x, 0x66)
	func(x, y uint) {
		// d.f00415984(x, y)
		func(x, y uint) {
			d.blocks[0].f0041987A()
			d.blocks[0].m00 = 0x0045061C
			// f00433360memset(d.f54[:], 0, 32)
			d.blocks[0].m68 = y
			d.blocks[0].m54 = x
			d.blocks[0].m58 = x & 0x0000FFFF
		}(x, y)
		d.blocks[0].m00 = 0x0044F28C
		// d.f0040E430()
		func() {
			d.blocks[0].m74 = 1
			d.blocks[0].m78 = 0
			d.blocks[0].m84 = 0
			// v0048F4A0 = dll.user32.SetLayeredWindowsAttributes
		}()
	}(x, 0x66)
	d.blocks[0].m00 = 0x0044E654
	d.blocks[0].m88 = 0x0044E620
	d.blocks[1].f0041987A() // 0x98
	d.blocks[2].f0041987A() // 0xEC
	v0014AAC4state := 0
	v0014AAC4state = 2
	// v0018A75C.f00422CDD()
	v0014AAC4state = 3
	// v0014A75C.f00 = 0x0044E1FC
	// v0018A75C.f00422D3E(1)
	v0014AAC4state = 4
	d.block2[0].f0041987A() // 0x44B40
	d.block2[0].m00 = 0x0044E4AC
	v0014AAC4state = 5
	d.block2[1].f0041987A() // 0x44B94
	d.block2[1].m00 = 0x0044E4AC
	v0014AAC4state = 6
	// f0041E335().f0C()
	v0014AAC4state = 7
	d.block2[2].f0041987A() // 0x44BEC
	d.block2[2].m00 = 0x0044DCAC
	v0014AAC4state = 8
	d.block2[3].f0041987A() // 0x44C60
	d.block2[3].m00 = 0x0044DCAC
	v0014AAC4state = 9
	d.block2[4].f0041987A() // 0x44CD4
	d.block2[4].m00 = 0x0044DCAC
	v0014AAC4state = 10
	d.block2[5].f0041987A() // 0x44D48
	d.block2[5].m00 = 0x0044DCAC
	v0014AAC4state = 11
	d.block2[6].f0041987A() // 0x44D48
	d.block2[6].m00 = 0x0044DCAC
	v0014AAC4state = 12
	d.block2[7].f0041987A() // 0x44E30
	d.block2[7].m00 = 0x0044DCAC
	v0014AAC4state = 13
	d.block2[8].f0041987A() // 0x44EA4
	d.block2[8].m00 = 0x0044DCAC
	v0014AAC4state = 14
	d.block2[9].f0040DA60() // 0x44F18
	v0014AAC4state = 15
	d.block2[10].f0040F670() // 0x44F84
	v0014AAC4state = 16
	d.block2[11].f0040F670() // 0x44FAC
	v0014AAC4state = 17
	d.block2[12].f0040F670() // 0x44FEC
	v0014AAC4state = 18
	d.block2[13].f0040F670() // 0x45004
	v0014AAC4state = 19
	d.block2[14].f0040F670() // 0x45110
	v0014AAC4state = 20
	d.block2[15].f0040F670() // 0x45128
	v0014AAC4state = 21
	d.block2[15].f00417837()
	d.block2[15].f00417837()
	// LoadIcon()
	// d.f00406920()
	func() {
		// read config.ini
		// fileName := path.Join(filepath.Dir(os.Args[0]), "config.ini")
		// var version [10]uint8
		// dll.kernel32.GetPrivateProfileString("LOGIN", "Version", "0.0.0", version[:], 10, fileName[:])
		var minor, major, patch uint8
		minor = 4
		major = 1
		patch = 44

		// dll.kernel32.GetPrivateProfileString("LOGIN", "TestVersion", "0.0.0", version[:], 10, fileName[:])
		// var testmajor uint8
		// testmajor = 1
		if v0046327C == 1 {
			// ...
		}
		d.major = major
		d.minor = minor
		d.patch = patch
	}()
	println(v0014AAC4state)
	// memcpy(d.name[:], textcode(81)) // "奇迹(MU)"
	// f00433360memset(d.pathName, 0, 256)
}

// MESSAGE_MAP
// const AFX_MSGMAP* theClass::GetMessageMap() const {
// 	return GetThisMessageMap();
// }
// const AFX_MSGMAP* theClass::GetThisMessageMap() {
// 	static const AFX_MSGMAP_ENTRY _messageEntries[] = {
// 		// -----------
// 		{WM_SYSCOMMAND,    0, 0, 0, AfxSIg_vwl, (AFX_PMSG)(AFX_PMSGW)(static_cast<void (AFX_MSG_CALL_ CWnd::*)(UINT, LPARAM)>(&ThisClass::OnSysCommand))},
// 		{WM_PAINT,         0, 0, 0, AfxSig_vv,  (AFX_PMSG)(AFX_PMSGW)(static_cast<void (AFX_MSG_CALL_ CWnd::*)(void)>(&ThisClass::OnPaint))},
// 		{WM_QUERYDRAGICON, 0, 0, 0, AfxSig_hv,  (AFX_PMSG)(AFX_PMSGW)(static_cast<HCURSOR (AFX_MSG_CALL_ CWnd::*)(void)>(&ThisClass::OnQueryDragIcon))},
// 		{WM_COMMAND, (WORD)BN_CLICKED, (WORD)IDOK, (WORD)IDOK, AfxSigCmd_v, static_cast<AFX_PMSG>&CmuDlg::OnBnClickedOk},
// 		{WM_COMMAND, (WORD)BN_CLICKED, (WORD)IDCANCEL, (WORD)IDCANCEL, AfxSigCmd_v, static_cast<AFX_PMSG>&CmuDlg::OnBnClickedCancel},
// 		// -----------
// 		{0, 0, 0, 0, AfxSig_end, (AFX_PMSG)0}
// 	};
// 	static const AFX_MSGMAP messageMap = {&TheBaseClass::GetThisMessageMap, &_messageENtries[0]};
// 	return &messageMap;
// }

func (d *muDlg) f0040CA40OnInitDialog() {
	// hUpdate := dll.user32.FindWindowA("#32770", "MU Auto Update")
	// if hUpdate != nil {
	// 	dll.user32.SendMessage(hUpdate, 0x10, 0, 0)
	// }
	// d.f0040E390()
	// dll.user32.SendMessage(d.m20, 0x80, 1, d.lParam) // 0x00212582, 0x80, 1, 0x01EF242F
	// dll.user32.SendMessage(d.m20, 0x80, 1, d.lParam) // 0x00212582, 0x80, 0, 0x01EF242F
	// d.f0040E880(0xBE)
	// d.f0040E480(4)
	// d.f0040E3A0(0xFF00FF,1) // set window style
	// d.f00416427("MU Auto Update")
	// d.f0041A58D(0)
	// d.f00407B10()
	// d.f0040C770()
	func() {
		// d.f0040C440()
		func() {
			// d.f0040B2B0()
			func() {
				d.m45004.f00403F50open(".\\parent_partition.inf")
				// d.m45110.f00409F90('\n')
				// ...
				parition := "1 电信专区\n" //
				partiNum := 0
				partiName := ""
				fmt.Sscanf(parition, "%d %s", &partiNum, &partiName)
			}()
			// d.f0040B480()
			func() {
				f := d.m45004.f00403F50open(".\\partition.inf")
				if f == nil {
					return
				}
				// d.m45110.f00409F90('\n')
				for !d.m45004.f00403FB0EOF() {
					var partitionRecord [512]uint8 // "1\t小学生\t192.168.0.102\t44405\thttp://www.baidu.com\n"
					d.m45004.f00403FC0fread(f, partitionRecord[:])
					partiNum := 1
					partiName := "小学生"
					partiIP := "192.168.0.102"
					partiPort := 44405
					partiSite := "http://www.baidu.com"
					// f00433F69sscanf(parition, "%d %s %s %d %s\n", &partiNum, &partiName, &partiIP, &partiPort, &partiSite)
					fmt.Sscanf(string(partitionRecord[:]), "%d %s %s %d %s\n", &partiNum, &partiName, &partiIP, &partiPort, &partiSite)
					p := new(partition) // 0x02596278 0x02496340
					p.f00405E00()
					p.m04Num = partiNum
					p.m3AName = partiName
					p.m7AIP = partiIP
					p.mBCPort = partiPort
					p.m08Site = partiSite
					d.m45128.f0040ADF0(p)
					d.m44FE8partiSum++
				}
				d.m45004.f00403F80close()
			}()
		}()
		pSum := len(d.m45128.partitions)
		if pSum == 0 {
			return
		}
		// d.m44FEC.f00409F90(15)

		// 渲染按钮
		i := 0
		for i < pSum {
			p := d.m45128.partitions[i]
			if p == nil {
				continue
			}
			// ecxNum := p.m04Num
			// eaxIndex := p.m00Index
			// esi := ecxNum - 1
			// eaxIndex%4
			b := new(button) //f00414B5Cnew, 0x02596E78, 0x02595008
			b.f0040F770()
			// b.f0041E1B1()
			// b.f0041E242()
			// b.f0040F750()
			// b.f0040F6E0()
			// ...
			i++
		}
	}()
}
