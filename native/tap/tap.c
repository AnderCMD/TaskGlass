/*
 * TaskGlass TAP: a tiny XAML Diagnostics "Tool Assembly Provider" that is
 * loaded into explorer.exe via InitializeXamlDiagnosticsEx.
 *
 * Windows 11's taskbar paints its own semi-transparent XAML fill on top of
 * the DWM accent. This DLL finds those fill elements (BackgroundFill /
 * BackgroundStroke) in the taskbar's visual tree and sets their Opacity to 0
 * while TaskGlass asks for it, and back to 1 when it stops.
 *
 * How it gets loaded (same scheme as TranslucentTB): TaskGlass installs a
 * WH_CALLWNDPROC hook (TaskGlassHookProc) on the taskbar's UI thread, which
 * makes Windows load this DLL into explorer.exe on that thread. DllMain then
 * starts a thread that calls InitializeXamlDiagnosticsEx on Explorer itself
 * (retrying until XAML is ready). XAML then loads this same DLL again as the
 * TAP and calls SetSite, which creates an Inst below.
 *
 * Control channel: a named shared memory block "Local\TaskGlass.Xaml.State.v1"
 * of two LONGs created by the TaskGlass process:
 *   [0] mode   TaskGlass -> TAP: 0 = leave the taskbar alone, 1 = hide the fill
 *   [1] status TAP -> TaskGlass: 0 = pending, 1 = attached, <0 = failure HRESULT
 * plus a manual-reset event "Local\TaskGlass.Xaml.Ready" set once status is
 * final. If the block does not exist (TaskGlass exited or crashed) the mode is
 * 0, so the stock look always comes back by itself.
 *
 * Threading: XAML objects may only be touched on the UI thread. The visual
 * tree callback already runs there; later mode changes are marshalled with
 * CoreDispatcher::RunAsync.
 *
 * Optional debugging: if %LOCALAPPDATA%\TaskGlass\tap.dump exists when the DLL
 * is loaded, every visual tree change is written to tap.log.
 */
#define WIN32_LEAN_AND_MEAN
#define COBJMACROS
#define INITGUID
#include <windows.h>
#include <ocidl.h>
#include <inspectable.h>
#include <xamlom.h>
#include <windows.ui.core.h>
#include <winstring.h>
#include <roapi.h>
#include <stdio.h>
#include <stdarg.h>
#include <wchar.h>

typedef __x_ABI_CWindows_CUI_CCore_CICoreDispatcher ICoreDispatcher_;
typedef __x_ABI_CWindows_CUI_CCore_CIDispatchedHandler IDispatchedHandler_;
typedef __x_ABI_CWindows_CUI_CCore_CIDispatchedHandlerVtbl IDispatchedHandlerVtbl_;

/* {B0A3F1E2-6C1D-4F5A-9E52-7D2A3C4B5E61} */
static const CLSID CLSID_TaskGlassTap = {0xb0a3f1e2, 0x6c1d, 0x4f5a, {0x9e, 0x52, 0x7d, 0x2a, 0x3c, 0x4b, 0x5e, 0x61}};
/* Windows.UI.Xaml.IUIElement */
static const IID IID_IUIElement_ = {0x676d0be9, 0xb65c, 0x41c6, {0xba, 0x40, 0x58, 0xcf, 0x87, 0xf2, 0x01, 0xc1}};
/* Windows.UI.Xaml.IDependencyObject, IFrameworkElement, Media.IVisualTreeHelperStatics.
 * IIDs and the slots below were read from Windows.UI.Xaml.winmd on a Windows 11 box. */
static const IID IID_IDependencyObject_ = {0x5c526665, 0xf60e, 0x4912, {0xaf, 0x59, 0x5f, 0xe0, 0x68, 0x0f, 0x08, 0x9d}};
static const IID IID_IFrameworkElement_ = {0xa391d09b, 0x4a99, 0x4b7c, {0x9d, 0x8d, 0x6f, 0xa5, 0xd0, 0x1f, 0x6f, 0xbf}};
static const IID IID_IVisualTreeHelperStatics_ = {0xe75758c4, 0xd25d, 0x4b1d, {0x97, 0x1f, 0x59, 0x6f, 0x17, 0xf1, 0x2b, 0xaa}};
#define SLOT_VTH_GETCHILD 10
#define SLOT_VTH_GETCOUNT 11
#define SLOT_VTH_GETPARENT 12
#define SLOT_FE_GETNAME 33
typedef HRESULT(STDMETHODCALLTYPE *GetChildFn)(void *self, void *parent, INT32 idx, void **child);
typedef HRESULT(STDMETHODCALLTYPE *GetCountFn)(void *self, void *ref, INT32 *count);
typedef HRESULT(STDMETHODCALLTYPE *GetParentFn)(void *self, void *ref, void **parent);
typedef HRESULT(STDMETHODCALLTYPE *GetNameFn)(void *self, HSTRING *name);
#define VT(o) (*(void ***)(o))

/* IAgileObject marker */
static const IID IID_IAgileObject_ = {0x94ea2b94, 0xe9cc, 0x49e0, {0xc0, 0xff, 0xee, 0x64, 0xca, 0x8f, 0x5b, 0x90}};

/* IUIElement vtable slot of put_Opacity: IUnknown(3) + IInspectable(3) +
 * get_DesiredSize, get_AllowDrop, put_AllowDrop, get_Opacity, put_Opacity.
 * Verified against C:\Windows\System32\WinMetadata\Windows.UI.Xaml.winmd. */
#define SLOT_PUT_OPACITY 10
typedef HRESULT(STDMETHODCALLTYPE *PutOpacityFn)(void *self, double v);

#define MAX_TRACKED 32
typedef struct {
    InstanceHandle handle;
    void *ui; /* IUIElement*, AddRef'd */
} Tracked;

typedef struct Inst Inst;
struct Inst {
    const IObjectWithSiteVtbl *tap;            /* IObjectWithSite */
    const IVisualTreeServiceCallbackVtbl *cb;  /* IVisualTreeServiceCallback */
    const IDispatchedHandlerVtbl_ *handler;    /* IDispatchedHandler */
    int id;
    IXamlDiagnostics *diag;
    IVisualTreeService *vts;
    ICoreDispatcher_ *disp;
    CRITICAL_SECTION cs;
    Tracked tracked[MAX_TRACKED];
    InstanceHandle scan_anchor; /* newest live element, set by the callback */
    volatile LONG scan_queued;
    DWORD last_scan;
    Inst *next;
};

#define INST_FROM(ptr, field) ((Inst *)((char *)(ptr)-offsetof(Inst, field)))

static CRITICAL_SECTION g_list_cs;
static Inst *g_insts;
static int g_next_id = 1;
static volatile LONG g_applied_mode;
static volatile LONG g_thread_started;
static BOOL g_dump;

/* ---------------------------------------------------------------- logging */
static void logf_(const wchar_t *fmt, ...) {
    wchar_t dir[MAX_PATH], path[MAX_PATH], line[1024], out[1100];
    DWORD n = GetEnvironmentVariableW(L"LOCALAPPDATA", dir, MAX_PATH);
    if (n == 0 || n >= MAX_PATH) return;
    _snwprintf(path, MAX_PATH, L"%ls\\TaskGlass\\tap.log", dir);
    va_list ap;
    va_start(ap, fmt);
    _vsnwprintf(line, 1023, fmt, ap);
    va_end(ap);
    line[1023] = 0;
    SYSTEMTIME t;
    GetLocalTime(&t);
    _snwprintf(out, 1099, L"%02d:%02d:%02d.%03d %ls\r\n", t.wHour, t.wMinute, t.wSecond, t.wMilliseconds, line);
    char u8[4400];
    int len = WideCharToMultiByte(CP_UTF8, 0, out, -1, u8, sizeof u8, NULL, NULL);
    if (len <= 1) return;
    HANDLE f = CreateFileW(path, FILE_APPEND_DATA | FILE_READ_ATTRIBUTES, FILE_SHARE_READ | FILE_SHARE_WRITE, NULL,
                           OPEN_ALWAYS, FILE_ATTRIBUTE_NORMAL, NULL);
    if (f == INVALID_HANDLE_VALUE) return;
    LARGE_INTEGER size;
    if (GetFileSizeEx(f, &size) && size.QuadPart > (1 << 20)) { /* 1 MiB cap */
        CloseHandle(f);
        return;
    }
    DWORD w;
    WriteFile(f, u8, len - 1, &w, NULL);
    CloseHandle(f);
}

/* ------------------------------------------------------------ shared state */
static LONG read_mode(void) {
    HANDLE m = OpenFileMappingW(FILE_MAP_READ, FALSE, L"Local\\TaskGlass.Xaml.State.v1");
    if (!m) return 0;
    LONG v = 0;
    LONG *p = (LONG *)MapViewOfFile(m, FILE_MAP_READ, 0, 0, sizeof(LONG));
    if (p) {
        v = *p;
        UnmapViewOfFile(p);
    }
    CloseHandle(m);
    return v;
}

static void publish_status(LONG status) {
    HANDLE m = OpenFileMappingW(FILE_MAP_WRITE, FALSE, L"Local\\TaskGlass.Xaml.State.v1");
    if (m) {
        LONG *p = (LONG *)MapViewOfFile(m, FILE_MAP_WRITE, 0, 0, 2 * sizeof(LONG));
        if (p) {
            p[1] = status;
            UnmapViewOfFile(p);
        }
        CloseHandle(m);
    }
    HANDLE e = OpenEventW(EVENT_MODIFY_STATE, FALSE, L"Local\\TaskGlass.Xaml.Ready");
    if (e) {
        SetEvent(e);
        CloseHandle(e);
    }
}

static void set_opacity(Inst *in, void *ui, double v) {
    void **vtbl = *(void ***)ui;
    HRESULT hr = ((PutOpacityFn)vtbl[SLOT_PUT_OPACITY])(ui, v);
    if (FAILED(hr)) logf_(L"[%d] put_Opacity failed: 0x%08lx", in->id, (unsigned long)hr);
}

/* Runs on the instance's UI thread. */
static void apply_all(Inst *in, LONG mode) {
    EnterCriticalSection(&in->cs);
    for (int i = 0; i < MAX_TRACKED; i++)
        if (in->tracked[i].ui) set_opacity(in, in->tracked[i].ui, mode ? 0.0 : 1.0);
    LeaveCriticalSection(&in->cs);
}

/* ------------------------------------------- IDispatchedHandler (agile) */
static HRESULT STDMETHODCALLTYPE H_QI(IDispatchedHandler_ *self, REFIID riid, void **out) {
    if (IsEqualIID(riid, &IID_IUnknown) || IsEqualIID(riid, &IID___x_ABI_CWindows_CUI_CCore_CIDispatchedHandler) ||
        IsEqualIID(riid, &IID_IAgileObject_)) {
        *out = self;
        return S_OK;
    }
    *out = NULL;
    return E_NOINTERFACE;
}
static ULONG STDMETHODCALLTYPE H_AddRef(IDispatchedHandler_ *self) { (void)self; return 2; }
static ULONG STDMETHODCALLTYPE H_Release(IDispatchedHandler_ *self) { (void)self; return 1; }
static void scan_from(Inst *in, InstanceHandle anchor);
static HRESULT STDMETHODCALLTYPE H_Invoke(IDispatchedHandler_ *self) {
    Inst *in = INST_FROM(self, handler);
    if (InterlockedExchange(&in->scan_queued, 0) && in->scan_anchor) scan_from(in, in->scan_anchor);
    apply_all(in, g_applied_mode);
    return S_OK;
}
static const IDispatchedHandlerVtbl_ g_handler_vtbl = {H_QI, H_AddRef, H_Release, H_Invoke};

static void dispatch_apply(Inst *in) {
    if (!in->disp) return;
    __x_ABI_CWindows_CFoundation_CIAsyncAction *act = NULL;
    HRESULT hr = __x_ABI_CWindows_CUI_CCore_CICoreDispatcher_RunAsync(
        in->disp, (__x_ABI_CWindows_CUI_CCore_CCoreDispatcherPriority)0, (IDispatchedHandler_ *)&in->handler, &act);
    if (FAILED(hr)) logf_(L"[%d] RunAsync failed: 0x%08lx", in->id, (unsigned long)hr);
    else if (act) IUnknown_Release((IUnknown *)act);
}

static BOOL any_tracked(void) {
    BOOL any = FALSE;
    EnterCriticalSection(&g_list_cs);
    for (Inst *in = g_insts; in && !any; in = in->next)
        for (int i = 0; i < MAX_TRACKED; i++)
            if (in->tracked[i].ui) { any = TRUE; break; }
    LeaveCriticalSection(&g_list_cs);
    return any;
}

/* status 2 = "controlling the taskbar fill"; republished every poll so a new
 * TaskGlass process (new shared block) learns about fills found by an earlier one. */
static void publish_filled(void) {
    HANDLE m = OpenFileMappingW(FILE_MAP_WRITE, FALSE, L"Local\\TaskGlass.Xaml.State.v1");
    if (!m) return;
    LONG *p = (LONG *)MapViewOfFile(m, FILE_MAP_WRITE, 0, 0, 2 * sizeof(LONG));
    if (p) {
        if (p[1] != 2) p[1] = 2;
        UnmapViewOfFile(p);
    }
    CloseHandle(m);
}

static DWORD WINAPI poll_thread(LPVOID unused) {
    (void)unused;
    for (;;) {
        if (any_tracked()) publish_filled();
        LONG m = read_mode();
        if (m != g_applied_mode) {
            g_applied_mode = m;
            logf_(L"mode -> %ld", (long)m);
            EnterCriticalSection(&g_list_cs);
            for (Inst *in = g_insts; in; in = in->next) dispatch_apply(in);
            LeaveCriticalSection(&g_list_cs);
        }
        Sleep(150);
    }
    return 0;
}


/* ------------------------------------------------ tracking and scanning */
static BOOL hs_eq(HSTRING h, const wchar_t *w) {
    UINT32 n = 0;
    const wchar_t *b = h ? WindowsGetStringRawBuffer(h, &n) : NULL;
    return b && wcslen(w) == n && wmemcmp(b, w, n) == 0;
}

/* Keeps a reference to a taskbar fill. Takes ownership of `ui`. UI thread. */
static BOOL track(Inst *in, InstanceHandle handle, void *ui, const wchar_t *what) {
    BOOL stored = FALSE, dup = FALSE;
    EnterCriticalSection(&in->cs);
    int slot = -1;
    for (int i = 0; i < MAX_TRACKED; i++) {
        if (in->tracked[i].ui && in->tracked[i].handle == handle) dup = TRUE;
        else if (!in->tracked[i].ui && slot < 0) slot = i;
    }
    if (!dup && slot >= 0) {
        in->tracked[slot].handle = handle;
        in->tracked[slot].ui = ui;
        stored = TRUE;
    }
    LeaveCriticalSection(&in->cs);
    if (!stored) {
        IUnknown_Release((IUnknown *)ui);
        if (!dup) logf_(L"[%d] tracked table full", in->id);
        return FALSE;
    }
    logf_(L"[%d] tracking %ls handle=%llu", in->id, what, (unsigned long long)handle);
    if (g_applied_mode) set_opacity(in, ui, 0.0);
    return TRUE;
}

typedef struct {
    Inst *in;
    void *statics;
    int nodes;
    int found;
} Scan;

static void scan_node(Scan *sc, void *dep, int depth, BOOL in_frame) {
    if (depth > 80 || sc->nodes > 20000) return;
    sc->nodes++;

    HSTRING name = NULL;
    void *fe = NULL;
    if (SUCCEEDED(IUnknown_QueryInterface((IUnknown *)dep, &IID_IFrameworkElement_, &fe)) && fe) {
        ((GetNameFn)VT(fe)[SLOT_FE_GETNAME])(fe, &name);
        IUnknown_Release((IUnknown *)fe);
    }
    if (g_dump) {
        HSTRING c = NULL;
        const wchar_t *cn = L"?", *nn = name ? WindowsGetStringRawBuffer(name, NULL) : L"";
        if (SUCCEEDED(IInspectable_GetRuntimeClassName((IInspectable *)dep, &c)) && c) cn = WindowsGetStringRawBuffer(c, NULL);
        logf_(L"[%d]   scan d=%d %ls \"%ls\"", sc->in->id, depth, cn, nn ? nn : L"");
        if (c) WindowsDeleteString(c);
    }
    BOOL frame = in_frame || hs_eq(name, L"TaskbarFrame");
    BOOL fill = hs_eq(name, L"BackgroundFill"), stroke = hs_eq(name, L"BackgroundStroke");
    if (name) WindowsDeleteString(name);

    if (frame && (fill || stroke)) {
        HSTRING cls = NULL;
        BOOL is_rect = SUCCEEDED(IInspectable_GetRuntimeClassName((IInspectable *)dep, &cls)) &&
                       hs_eq(cls, L"Windows.UI.Xaml.Shapes.Rectangle");
        if (cls) WindowsDeleteString(cls);
        InstanceHandle h = 0;
        void *ui = NULL;
        if (is_rect && SUCCEEDED(IXamlDiagnostics_GetHandleFromIInspectable(sc->in->diag, (IInspectable *)dep, &h)) &&
            SUCCEEDED(IUnknown_QueryInterface((IUnknown *)dep, &IID_IUIElement_, &ui)) && ui) {
            if (track(sc->in, h, ui, fill ? L"BackgroundFill" : L"BackgroundStroke")) sc->found++;
        }
    }

    INT32 count = 0;
    if (FAILED(((GetCountFn)VT(sc->statics)[SLOT_VTH_GETCOUNT])(sc->statics, dep, &count))) return;
    for (INT32 i = 0; i < count; i++) {
        void *child = NULL;
        if (SUCCEEDED(((GetChildFn)VT(sc->statics)[SLOT_VTH_GETCHILD])(sc->statics, dep, i, &child)) && child) {
            scan_node(sc, child, depth + 1, frame);
            IUnknown_Release((IUnknown *)child);
        }
    }
}

/* From any live element: climb to the root of its XAML tree, then search the
 * whole tree for the taskbar fills. Runs on the UI thread. */
static void scan_from(Inst *in, InstanceHandle anchor) {
    IInspectable *obj = NULL;
    if (FAILED(IXamlDiagnostics_GetIInspectableFromHandle(in->diag, anchor, &obj)) || !obj) return;
    void *cur = NULL;
    HRESULT hr = IInspectable_QueryInterface(obj, &IID_IDependencyObject_, &cur);
    IInspectable_Release(obj);
    if (FAILED(hr) || !cur) return;

    HSTRING cn = NULL;
    WindowsCreateString(L"Windows.UI.Xaml.Media.VisualTreeHelper", 38, &cn);
    void *statics = NULL;
    hr = RoGetActivationFactory(cn, &IID_IVisualTreeHelperStatics_, &statics);
    WindowsDeleteString(cn);
    if (FAILED(hr) || !statics) {
        logf_(L"[%d] VisualTreeHelper unavailable: 0x%08lx", in->id, (unsigned long)hr);
        IUnknown_Release((IUnknown *)cur);
        return;
    }

    for (int up = 0; up < 200; up++) {
        void *parent = NULL;
        if (FAILED(((GetParentFn)VT(statics)[SLOT_VTH_GETPARENT])(statics, cur, &parent)) || !parent) break;
        IUnknown_Release((IUnknown *)cur);
        cur = parent;
    }
    Scan sc = {in, statics, 0, 0};
    scan_node(&sc, cur, 0, FALSE);
    logf_(L"[%d] scan: %d nodes, %d fill element(s) found", in->id, sc.nodes, sc.found);
    if (sc.found) publish_status(2);
    IUnknown_Release((IUnknown *)cur);
    IUnknown_Release((IUnknown *)statics);
}

/* ---------------------------------------- IVisualTreeServiceCallback */
static BOOL wants(const wchar_t *name) {
    return name && (wcscmp(name, L"BackgroundFill") == 0 || wcscmp(name, L"BackgroundStroke") == 0);
}

static HRESULT STDMETHODCALLTYPE C_QI(IVisualTreeServiceCallback *self, REFIID riid, void **out) {
    if (IsEqualIID(riid, &IID_IUnknown) || IsEqualIID(riid, &IID_IVisualTreeServiceCallback)) {
        *out = self;
        return S_OK;
    }
    *out = NULL;
    return E_NOINTERFACE;
}
static ULONG STDMETHODCALLTYPE C_AddRef(IVisualTreeServiceCallback *self) { (void)self; return 2; }
static ULONG STDMETHODCALLTYPE C_Release(IVisualTreeServiceCallback *self) { (void)self; return 1; }

static HRESULT on_change(IVisualTreeServiceCallback *self, ParentChildRelation rel, VisualElement el,
                         VisualMutationType type) {
    Inst *in = INST_FROM(self, cb);
    if (g_dump)
        logf_(L"[%d] t%lu %ls parent=%llu handle=%llu type=%ls name=%ls children=%u", in->id, GetCurrentThreadId(), type == Add ? L"ADD" : L"DEL",
              (unsigned long long)rel.Parent, (unsigned long long)el.Handle, el.Type ? el.Type : L"?",
              el.Name ? el.Name : L"", el.NumChildren);

    if (type == Add && wants(el.Name) && el.Type && wcscmp(el.Type, L"Windows.UI.Xaml.Shapes.Rectangle") == 0 && in->diag) {
        /* a fill appearing live: only keep it if it sits inside a TaskbarFrame,
         * which the scan verifies; just queue a scan from here. */
        in->scan_anchor = el.Handle;
        if (!InterlockedExchange(&in->scan_queued, 1)) dispatch_apply(in);
    } else if (type == Add && rel.Parent != 0 && in->diag && GetTickCount() - in->last_scan > 1500) {
        /* any other live element works as an anchor to find the taskbar tree */
        in->last_scan = GetTickCount();
        in->scan_anchor = el.Handle;
        if (!InterlockedExchange(&in->scan_queued, 1)) dispatch_apply(in);
    } else if (type == Remove) {
        EnterCriticalSection(&in->cs);
        for (int i = 0; i < MAX_TRACKED; i++)
            if (in->tracked[i].ui && in->tracked[i].handle == el.Handle) {
                IUnknown_Release((IUnknown *)in->tracked[i].ui);
                in->tracked[i].ui = NULL;
            }
        LeaveCriticalSection(&in->cs);
    }
    return S_OK;
}
/* The service hands ownership of the BSTRs in `el` to the callback. */
static HRESULT STDMETHODCALLTYPE C_OnChange(IVisualTreeServiceCallback *self, ParentChildRelation rel, VisualElement el,
                                            VisualMutationType type) {
    HRESULT hr = on_change(self, rel, el, type);
    SysFreeString(el.SrcInfo.FileName);
    SysFreeString(el.SrcInfo.Hash);
    SysFreeString(el.Type);
    SysFreeString(el.Name);
    return hr;
}
static const IVisualTreeServiceCallbackVtbl g_cb_vtbl = {C_QI, C_AddRef, C_Release, C_OnChange};

/* ------------------------------------------------------ IObjectWithSite */
static HRESULT STDMETHODCALLTYPE T_QI(IObjectWithSite *self, REFIID riid, void **out) {
    if (IsEqualIID(riid, &IID_IUnknown) || IsEqualIID(riid, &IID_IObjectWithSite)) {
        *out = self;
        return S_OK;
    }
    *out = NULL;
    return E_NOINTERFACE;
}
static ULONG STDMETHODCALLTYPE T_AddRef(IObjectWithSite *self) { (void)self; return 2; }
static ULONG STDMETHODCALLTYPE T_Release(IObjectWithSite *self) { (void)self; return 1; }
static HRESULT STDMETHODCALLTYPE T_GetSite(IObjectWithSite *self, REFIID riid, void **out) {
    (void)self; (void)riid;
    *out = NULL;
    return E_FAIL;
}
static DWORD WINAPI advise_thread(LPVOID arg) {
    Inst *in = (Inst *)arg;
    HRESULT hr = IVisualTreeService_AdviseVisualTreeChange(in->vts, (IVisualTreeServiceCallback *)&in->cb);
    logf_(L"[%d] TAP attached, advise=0x%08lx dump=%d", in->id, (unsigned long)hr, (int)g_dump);
    publish_status(SUCCEEDED(hr) ? 1 : (LONG)hr);
    return 0;
}

static HRESULT STDMETHODCALLTYPE T_SetSite(IObjectWithSite *self, IUnknown *site) {
    Inst *in = INST_FROM(self, tap);
    if (!site) return S_OK; /* explorer never tears the TAP down while running */

    HRESULT hr = IUnknown_QueryInterface(site, &IID_IXamlDiagnostics, (void **)&in->diag);
    if (FAILED(hr)) { logf_(L"[%d] site is not IXamlDiagnostics: 0x%08lx", in->id, (unsigned long)hr); return hr; }
    hr = IUnknown_QueryInterface(site, &IID_IVisualTreeService, (void **)&in->vts);
    if (FAILED(hr)) { logf_(L"[%d] site is not IVisualTreeService: 0x%08lx", in->id, (unsigned long)hr); return hr; }

    IInspectable *d = NULL;
    hr = IXamlDiagnostics_GetDispatcher(in->diag, &d);
    if (SUCCEEDED(hr) && d) {
        hr = IInspectable_QueryInterface(d, &IID___x_ABI_CWindows_CUI_CCore_CICoreDispatcher, (void **)&in->disp);
        IInspectable_Release(d);
    }
    if (FAILED(hr))
        logf_(L"[%d] no CoreDispatcher (mode changes after startup will not apply): 0x%08lx", in->id, (unsigned long)hr);

    EnterCriticalSection(&g_list_cs);
    in->next = g_insts;
    g_insts = in;
    LeaveCriticalSection(&g_list_cs);

    g_applied_mode = read_mode();
    if (InterlockedCompareExchange(&g_thread_started, 1, 0) == 0) {
        HANDLE th = CreateThread(NULL, 0, poll_thread, NULL, 0, NULL);
        if (th) CloseHandle(th);
    }

    /* Advising from SetSite's thread can hang; the service hops to the UI
     * thread by itself for the callbacks (same approach as TranslucentTB). */
    HANDLE th = CreateThread(NULL, 0, advise_thread, in, 0, NULL);
    if (th) CloseHandle(th);
    else logf_(L"[%d] CreateThread failed: %lu", in->id, GetLastError());
    return S_OK;
}
static const IObjectWithSiteVtbl g_tap_vtbl = {T_QI, T_AddRef, T_Release, T_SetSite, T_GetSite};

/* ------------------------------------------------------- class factory */
static HRESULT STDMETHODCALLTYPE F_QI(IClassFactory *self, REFIID riid, void **out) {
    if (IsEqualIID(riid, &IID_IUnknown) || IsEqualIID(riid, &IID_IClassFactory)) {
        *out = self;
        return S_OK;
    }
    *out = NULL;
    return E_NOINTERFACE;
}
static ULONG STDMETHODCALLTYPE F_AddRef(IClassFactory *self) { (void)self; return 2; }
static ULONG STDMETHODCALLTYPE F_Release(IClassFactory *self) { (void)self; return 1; }
static HRESULT STDMETHODCALLTYPE F_Create(IClassFactory *self, IUnknown *outer, REFIID riid, void **out) {
    (void)self;
    if (outer) return CLASS_E_NOAGGREGATION;
    Inst *in = (Inst *)calloc(1, sizeof *in); /* lives as long as explorer: never freed */
    if (!in) return E_OUTOFMEMORY;
    in->tap = &g_tap_vtbl;
    in->cb = &g_cb_vtbl;
    in->handler = &g_handler_vtbl;
    InitializeCriticalSection(&in->cs);
    EnterCriticalSection(&g_list_cs);
    in->id = g_next_id++;
    LeaveCriticalSection(&g_list_cs);
    return T_QI((IObjectWithSite *)&in->tap, riid, out);
}
static HRESULT STDMETHODCALLTYPE F_Lock(IClassFactory *self, BOOL lock) { (void)self; (void)lock; return S_OK; }
static const IClassFactoryVtbl g_factory_vtbl = {F_QI, F_AddRef, F_Release, F_Create, F_Lock};
static IClassFactory g_factory = {(IClassFactoryVtbl *)&g_factory_vtbl};

STDAPI DllGetClassObject(REFCLSID clsid, REFIID riid, LPVOID *out) {
    if (!IsEqualCLSID(clsid, &CLSID_TaskGlassTap)) return CLASS_E_CLASSNOTAVAILABLE;
    return F_QI(&g_factory, riid, out);
}

/* Explorer must never unload us while XAML still holds our callbacks. */
STDAPI DllCanUnloadNow(void) { return S_FALSE; }

/* --------------------------------------------- loading into explorer */
static HINSTANCE g_inst;

/* Hook procedure: its only job is to make Windows load this DLL into the
 * thread that owns the taskbar window. */
LRESULT CALLBACK TaskGlassHookProc(int code, WPARAM w, LPARAM l) { return CallNextHookEx(NULL, code, w, l); }

typedef HRESULT(WINAPI *InitXamlDiagFn)(PCWSTR, DWORD, PCWSTR, PCWSTR, CLSID, PCWSTR);
typedef struct {
    InitXamlDiagFn fn;
    const wchar_t *path;
    wchar_t conn[32];
    HRESULT hr;
} InitCall;

/* XAML Diagnostics can only be initialised once per thread, so every attempt
 * gets a fresh thread. */
static DWORD WINAPI init_call_thread(LPVOID arg) {
    InitCall *c = (InitCall *)arg;
    c->hr = c->fn(c->conn, GetCurrentProcessId(), NULL, c->path, CLSID_TaskGlassTap, NULL);
    return 0;
}

static DWORD WINAPI install_thread(LPVOID unused) {
    (void)unused;
    wchar_t path[MAX_PATH];
    DWORD n = GetModuleFileNameW(g_inst, path, MAX_PATH);
    if (n == 0 || n >= MAX_PATH) { publish_status(-1); return 1; }

    HMODULE wux = LoadLibraryExW(L"Windows.UI.Xaml.dll", NULL, LOAD_LIBRARY_SEARCH_SYSTEM32);
    InitXamlDiagFn fn = wux ? (InitXamlDiagFn)GetProcAddress(wux, "InitializeXamlDiagnosticsEx") : NULL;
    if (!fn) {
        logf_(L"InitializeXamlDiagnosticsEx unavailable");
        publish_status(-2);
        return 1;
    }

    HRESULT hr = E_FAIL;
    for (int attempt = 1; attempt <= 60; attempt++) { /* 60 x 500 ms = 30 s */
        InitCall c = {fn, path, {0}, E_FAIL};
        _snwprintf(c.conn, 31, L"VisualDiagConnection%d", attempt);
        HANDLE t = CreateThread(NULL, 0, init_call_thread, &c, 0, NULL);
        if (!t) { hr = HRESULT_FROM_WIN32(GetLastError()); break; }
        WaitForSingleObject(t, INFINITE);
        CloseHandle(t);
        hr = c.hr;
        logf_(L"InitializeXamlDiagnosticsEx(%ls) = 0x%08lx", c.conn, (unsigned long)hr);
        if (SUCCEEDED(hr)) break;
        Sleep(500);
    }
    if (SUCCEEDED(hr)) {
        HMODULE pin;
        GetModuleHandleExW(GET_MODULE_HANDLE_EX_FLAG_PIN | GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS, (LPCWSTR)&install_thread, &pin);
        /* status becomes 1 (or an error) once SetSite/Advise finish */
    } else {
        publish_status((LONG)hr);
    }
    return 0;
}

static BOOL host_is_explorer(void) {
    wchar_t exe[MAX_PATH];
    if (!GetModuleFileNameW(NULL, exe, MAX_PATH)) return FALSE;
    const wchar_t *base = wcsrchr(exe, L'\\');
    return base && _wcsicmp(base + 1, L"explorer.exe") == 0;
}

BOOL WINAPI DllMain(HINSTANCE inst, DWORD reason, LPVOID reserved) {
    (void)reserved;
    if (reason == DLL_PROCESS_ATTACH) {
        g_inst = inst;
        DisableThreadLibraryCalls(inst);
        InitializeCriticalSection(&g_list_cs);
        wchar_t dir[MAX_PATH], path[MAX_PATH];
        DWORD n = GetEnvironmentVariableW(L"LOCALAPPDATA", dir, MAX_PATH);
        if (n > 0 && n < MAX_PATH) {
            _snwprintf(path, MAX_PATH, L"%ls\\TaskGlass\\tap.dump", dir);
            g_dump = GetFileAttributesW(path) != INVALID_FILE_ATTRIBUTES;
        }
        /* The TaskGlass process also loads this DLL to get the hook procedure's
         * address; only start the installer inside Explorer itself. */
        if (host_is_explorer()) {
            HANDLE t = CreateThread(NULL, 0, install_thread, NULL, 0, NULL);
            if (t) CloseHandle(t);
        }
    }
    return TRUE;
}
