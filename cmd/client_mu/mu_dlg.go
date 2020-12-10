package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xujintao/balgass/win"
)

// dialogEx base dialog
type dialogEx struct {
	m_hWnd win.HWND
}

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
	// return AfxCallWndProc(pWnd, hWnd, msg, wParam, lParam) // -> CWnd::WindowProc, f00418F37 -> CWnd::OnWndMsg, f0041C89F
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
		// f0041C6E8mfcAfxHookWindowCreate(CWnd* pWnd)
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

type partitionManager struct {
	// partitions
	partitions        []*partition
	m0CpartitionsHead *partition // head
	m10partitionsTail *partition // tail
	m14partitionsSafe *partition // safe
}

func (t *partitionManager) f0040ADF0(p *partition) {
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
	// 	// f00415460()
	// }
	// if v00A6DE08.f04 != 0 {
	// 	return 1
	// }
	var ret uint
	// ret = v0048F8C4.f00425146(f00417808)
	// if ret == 0 {
	// 	// f00415460()
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
	m88bindStatusCallback func()
	m140                  int
	verfile               [60]uint8  // 0x3F930
	host                  [256]uint8 // 0x3F96C
	port                  uint16     // 0x3FA6C
	user                  [50]uint8  // 0x3FA6E
	passwd                [50]uint8  // 0x3FAA0
	m3FAEC                int
	m3FAF0                int
	m3FAF8verCur          int
	major                 uint8 // 0x3FAFC
	minor                 uint8 // 0x3FAFD
	patch                 uint8 // 0x3FAFE
	m3FB00vers            [1000]struct {
		m00done int
		m04     [16]uint8
		m14     int
	}
	m44924vernum   int
	name           [256]uint8 // 0x44928, "奇迹("
	dir            [260]uint8 // 0x44A28
	block2         [100]block // 0x44B40 0x44B94 0x44BEC 0x44C60 0x44CD4
	m44E30btnStart button     // 0x44E30
	// m44F18
	// m44F84 label
	// m44FAC
	m44FE0                 int
	lParam                 uint             // 0x44FE4
	m44FE8partiSum         int              // 0x44FE8
	m44FEC                 struct{}         // 0x44FEC, 0x0018FC74
	m44FF8btns             []*button        // m44FF8 0x0018FC80, m44FFC 0x0018FC84
	m45004                 t0018FC8C        // 0x45004, 0x0018FC8C
	m45110                 struct{}         // 0x45110, 0x0018FD98
	m45128partitionManager partitionManager // 0x45128, 0x0018FDB0
}

func (d *muDlg) f0040AEC0construct(x uint) {
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
		var major, minor, patch uint8
		major = 1
		minor = 4
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
	// f00433360memset(d.host, 0, 256)
}

// MESSAGE_MAP
// const AFX_MSGMAP* theClass::GetMessageMap() const { // f28
// 	return GetThisMessageMap();
// }
// const AFX_MSGMAP* theClass::GetThisMessageMap() {
// 	static const AFX_MSGMAP_ENTRY _messageEntries[] = { // 0x0044E960
// 		// -----------
// 		{WM_PAINT, 0, 0, 0, 0x13, 0x0040A0F0},
// 		{WM_QUERYDRAGICON, 0, 0, 0, 0x28, 0x004062A0},
// 		{0x19, 0, 0, 0, 0x08, 0x0041B0FB}
// 		{WM_LBUTTONDOWN, 0, 0, 0, 0x35, 0x0040E560},
// 		{WM_CLOSE, 0, 0, 0, 0x13, 0x00407B00},
// 		{WM_COMMAND, BN_CLICKED, 0x7D6, 0x7D6, 0x39, f004064A0}, // 官方网站
// 		{WM_COMMAND, BN_CLICKED, 0x7D9, 0x7D9, 0x39, f004064E0}, // 游戏论坛
// 		{WM_COMMAND, BN_CLICKED, 0x7D7, 0x7D7, 0x39, f0040A8E0}, // 退出
// 		{WM_COMMAND, BN_CLICKED, 0x7D2, 0x7D2, 0x39, f0040A5A0start}, // 游戏开始
// 		{WM_COMMAND, BN_CLICKED, 0x7D0, 0x7D0, 0x39, f00406520set}, // 游戏设置
// 		{WM_COMMAND, BN_CLICKED, 0x7D1, 0x7D1, 0x39, f004062B0}, // 注册账号
// 		{WM_COMMAND, BN_CLICKED, 0x1965, 0x1975, 0x3B, f00408AA0active}, // 激活专区
// 		{WM_COMMAND, BN_CLICKED, 0x7DB, 0x7DB, 0x39, f004062F0}, // 查看用户协议
// 		{WM_COMMAND, BN_CLICKED, 0x7DC, 0x7DC, 0x39, f0040B270}, // 同意
// 		{WM_COMMAND, BN_CLICKED, 0, 0, 0x14, f0040A900onTimer}, // 定时器
// 		{WM_MOVE, 0, 0, 0, 0x17, f00407F30}, //{0x3, 0, 0, 0, 0x17, 0x00407F30},
// 		{0x7E8, 0, 0, 0, 0x0E, f0040A9A0net}, // User Message, net

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
					var partitionRecord [512]uint8 // "1\t双喜\t192.168.0.102\t44405\thttp://www.baidu.com\n"
					d.m45004.f00403FC0fread(f, partitionRecord[:])
					partiNum := 1
					partiName := "双喜"
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
					d.m45128partitionManager.f0040ADF0(p)
					d.m44FE8partiSum++
				}
				d.m45004.f00403F80close()
			}()
		}()
		pSum := len(d.m45128partitionManager.partitions)
		if pSum == 0 {
			return
		}
		// d.m44FEC.f00409F90(15)

		// 渲染按钮
		i := 0
		for i < pSum {
			p := d.m45128partitionManager.partitions[i]
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

func (d *muDlg) f00408AA0active() {
	// 拿到按钮索引
	index := 0

	// 拿到分区
	parti := d.m45128partitionManager.partitions[index]

	// disable start button
	/*if d.m44E30btnStart.f004164F0isButtonEnable() {
		d.m44E30btnStart.f0041650BenableButton(false)
	}*/

	v00463280IP = parti.m7AIP
	v004632B4port = uint16(parti.mBCPort)
	// d.f004084B0()
	func() {
		// d.f00408410(false) // disable all active buttons
		// d.f00406330(false) // disable ?
		v004633D0conn.f0040CB70()
		// f00434570remove("mu.tmp")
		// f00434A24createDir("Temp")
		// dll.user32.SetTimer(d.m_hWnd, 10, 500, 0) // 500ms
		// dll.user32.SetTimer(d.m_hWnd, 1, 100, 0) // 100ms
	}()
}

func (d *muDlg) f0040A900onTimer(nIDEvent int) {
	switch nIDEvent {
	case 1:
		// dll.user32.KillTimer(d.m_hWnd, 1)
		// d.f00408D80(v0046F448msg.Get(102)) // "连接服务器"
		v004633D0conn.f0040CC30socket(d.m_hWnd)
		v004633D0conn.f0040CD00dial(v00463280IP, v004632B4port, 0x7E8)
	default:
	}

	// d.f0041AB1E()
	func() {
		// v009CDE08 := v0048F8C8.f0042566A(f0041547C)
		// if v00A6DE08 == 0 {
		// 	// f00415460()
		// }
		// d.f118(v009CDE08.m5C, v009CDE08.m60, v009CDE08.m64) // d.f00418DDF(), CWnd::DefWindowProc(WM_TIMER, 10, 0)
	}()
}

func f0040CAF0xor(s []uint8, l int) {
	var keys = [3]uint8{0xFC, 0xCF, 0xAB}
	for i, c := range s {
		s[i] = c ^ keys[i%3]
	}
}

func (d *muDlg) f004066F0setVersionInfo(host []uint8, port uint16, user, passwd, verfile []uint8) {
	f0040CAF0xor(host, 100)
	f0040CAF0xor(user, 20)
	f0040CAF0xor(passwd, 20)
	f0040CAF0xor(verfile, 20)
	copy(d.host[:], host)
	d.port = port
	copy(d.user[:], user)
	copy(d.passwd[:], passwd)
	copy(d.verfile[:], verfile)
}

// var v00463368 [100]uint8
var v00463368buf bytes.Buffer

func f00405E50readOne() int {
	var c uint8
	var buf [1]uint8
	for {
		v004633CCfd.Read(buf[:]) // f0043446Cfgetc(v004633CCfd)
		c = buf[0]
		if c == 0xFF {
			return 2
		}
		if c == '/' {
			// ...
		}
		// if !f00434197(int(c)) {
		// 	break
		// }
	}
	switch c {
	case '"':
		// 0x0040601C:
		v004633CCfd.Read(buf[:]) // f0043446Cfgetc(v004633CCfd)
		c = buf[0]
		for c != 0xFF {
			if c == '"' {
				return 0
			}
			v00463368buf.WriteByte(c)
			v004633CCfd.Read(buf[:]) // f0043446Cfgetc(v004633CCfd)
			c = buf[0]
		}
	case '#':
	case '$':
	case '%':
	case '&':
	case '*':
	}
	return 0
}

func (d *muDlg) f00409AD0parseVersionFile() int {
	// d.f00408D80(v0046F448msg.Get(112)) // 分析更新信息
	d.m44924vernum = 0
	v004633CCfd, _ = os.Open(string(d.verfile[:])) // = f00432686fopen(d.verfile[:], "r")
	if v004633CCfd == nil {
		// d.f00408D80(v0046F448msg.Get(113)) // 更新信息读入失败
		return 0
	}
	for {
		v := f00405E50readOne()
		if v == 2 { // EOF
			break
		}
		if v == 0 {
			copy(d.m3FB00vers[d.m44924vernum].m04[:], v00463368buf.Bytes())
			subs := strings.Split(v00463368buf.String(), ".") // f00433E3Estrtok(v00463368buf.Bytes(), '.')
			major, _ := strconv.Atoi(subs[0])
			minor, _ := strconv.Atoi(subs[1])
			patch, _ := strconv.Atoi(subs[2])
			d.m3FB00vers[d.m44924vernum].m14 = major<<16 | minor<<8 | patch
			d.m44924vernum++
			if d.m44924vernum > 999 {
				// ...
			}
		}
	}
	v004633CCfd.Close() //f00432546fclose(v004633CCfd)
	return d.m44924vernum
}

func f0040AB70(d *muDlg) {
	/*
		var url [256]uint8
		dll.user32.wsprintfA(url[:], "http://%s/%s", d.host[:], d.verfile[:])
		if dll.urlmon.URLDownloadToFileA(0, url[:], d.verfile[:], 0, d.m88bindStatusCallback) != S_OK {
			d.f00408D80(v0046F448msg.Get(116)) // 列表信息接收失败！#2
			return
		}
	*/
	if d.f00409AD0parseVersionFile() == 0 {
		return
	}
	d.m3FAF8verCur = int(d.major)<<16 | int(d.minor)<<8 | int(d.patch)
	d.m3FB00vers[0].m00done = 0
	for num := d.m44924vernum; num > 0; num-- {
		if d.m3FB00vers[num-1].m14 <= d.m3FAF8verCur {
			break
		}
	}
	// draw
	os.Remove(string(d.verfile[:])) // f00434570remove(d.verfile)
	d.m3FAF0 = 1
	d.m3FAEC = 1
	// d.f0040A4A0()
	func() {
		// d.f00408410(true) // enable all active buttons
		// d.f00406330(true) // enable xx
		// d.f00408D80(v0046F448msg.Get(126)) // 选择“服务器”后，点击“游戏开始”进入游戏。
		d.m140 = 1
		d.m3FAF0 = 1
		// d.m44E30btnStart.f0041650BenableButton(true)
	}()
}

func (d *muDlg) f0040C3C0reqVersionFile() bool {
	// f00434570remove(string(d.verfile))
	// d.f00408D80(v0046F448msg.Get(109)) // 正在接收更新信息
	// d.f0041A502(1)
	// d.m44FAC.f0040F480(100)
	// d.m44E30btnStart.f0041650BenableButton(false)
	// f0042145B(f0040AB70, d, 0, 0, 0, 0)
	return false
}

func f00402C90handle(d *muDlg) {
	for {
		buf := v004633D0conn.f0040CEE0getPacket()
		if buf == nil {
			return
		}
		var code uint8
		switch buf[0] {
		case 0xC1:
			code = buf[2]
		case 0xC2:
			code = buf[3]
		}
		switch code {
		case 0x00: // connect result C1 04 00 01
			// f00402B10(buf)
			func(buf []uint8) {
				if buf[3] != 1 {
					return
				}
				var code uint8 = 5
				if d.m44FE0 == 1 { // d := f00417837AfxGetApp().m_hWnd
					code = 4
				}
				// f00402AB0(code, d.major, d.minor, d.patch)
				func(code, major, minor, patch uint8) {
					buf := [6]uint8{0xC1, 0x06, code, major, minor, patch}
					v004633D0conn.f0040CF30write(buf[:], len(buf))
				}(code, d.major, d.minor, d.patch)
			}(buf)
		case 0x02: // version match C1 04 02 01
			// f00402C20()
			func() {
				// d.f00408410(true) // enable all active buttons
				// d.f00406330(true) // enable ? buttons
				// d.f00408D80(v0046F448msg.Get(126)) // 选择“服务器”后，点击“游戏开始”进入游戏
				d.m140 = 1
				d.m3FAF0 = 1
				// d.m44E30btnStart.f0041650BenableButton(true)
			}()
		case 0x04: // auto update with http
			// f00402C30(buf)
			func(buf []uint8) {
				info := struct {
					version [3]uint8
					host    [100]uint8
					port    uint16
					user    [20]uint8
					passwd  [20]uint8
					verfile [20]uint8
				}{}
				d.f004066F0setVersionInfo(info.host[:], info.port, info.user[:], info.passwd[:], info.verfile[:])
				if d.f0040C3C0reqVersionFile() == false {
					// d.f00408D80(v0046F448msg.Get(138)) // 下载失败。请重新安装。
				}
			}(buf)
		case 0x05: // auto update with ftp
			// f00402B60(buf)
		}
	}
}

func (d *muDlg) f0040A9A0net(wParam, lParam uint) {
	switch lParam {
	case 1: // FD_READ
		v004633D0conn.f0040D090read()
		f00402C90handle(d)
	case 2: // FD_WRITE
		v004633D0conn.f0040D010write()
	case 8: // FD_ACCEPT
	case 16: // FD_CONNECT
		// dll.ws2_32.WSAAsyncSelect(v004633D0conn.fd, d.m_hWnd, 0x7E8, FD_READ|FD_WRITE|FD_CLOSE)
		if win.WSAGetLastError() == ^uint32(0) {
			v004633D0conn.f0040CCE0close()
		}
	case 32: // FD_CLOSE
	}
}

func (d *muDlg) f0040A5A0start() {
	// d.f00406380() // ?
	// d.f004078D0() // clean tmp directory
	v004633D0conn.f0040CCE0close()
	// d.f00406550()
	// 						   hWnd, Op,    FileName, Param,                             Dir, showCmd
	// dll.shell32.ShellExecute(0, "open", "main.exe", "connect /u192.168.0.102 /p44405", 0, SW_SHOW)
}

func (d *muDlg) f00406520set() {
	// f004025C0()
}
