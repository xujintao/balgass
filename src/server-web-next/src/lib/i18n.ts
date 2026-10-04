export const LOCALE_COOKIE = 'r2f2-locale';
export const locales = ['en', 'zh-CN', 'es'] as const;
export type Locale = (typeof locales)[number];

export function resolveLocale(value: string | undefined): Locale {
  return locales.find((locale) => locale === value) ?? 'en';
}

const en = {
  description: 'Sign up with email and sign in securely with a passkey.',
  language: 'Language',
  accountNavigation: 'Account navigation',
  accountMenu: 'Account menu',
  login: 'Log in',
  signup: 'Sign up',
  profile: 'Profile',
  security: 'Account security',
  gameAccounts: 'Game accounts',
  myGameAccounts: 'My accounts',
  noGameAccounts: 'No game accounts yet.',
  noGameCharacters: 'No characters yet',
  gameCharacterLevel: (level: number) => `Level ${level}`,
  createGameAccount: 'Create game account',
  gameAccountName: 'Account name',
  gameAccountPassword: 'Password',
  gameAccountPasswordConfirmation: 'Confirm password',
  creatingGameAccount: 'Creating…',
  logout: 'Log out',
  logoutFailed: 'Could not log out. Please try again.',
  homeEyebrow: 'Welcome to r2f2',
  homeTitle: 'Your r2f2 starts here.',
  homeText:
    'A place connecting players with great content. We are building the new r2f2 website. More is coming soon.',
  accountUnavailable: 'Your account is temporarily unavailable',
  serviceUnavailable:
    'The service is temporarily unavailable. Please try again later.',
  loginFirst: 'Please log in',
  restoringSession: 'Restoring your session…',
  sessionEnded:
    'Your session has ended. Log in again with a passkey or email code.',
  goToLogin: 'Go to login',
  emailVerified: 'Email verified',
  setUpPasskey: 'Set up your passkey',
  passkeyNextTime:
    'Next time, use your fingerprint, face, or device PIN to log in.',
  settingUp: 'Setting up…',
  createPasskey: 'Create passkey',
  later: 'Set up later',
  createAccount: 'Create your account',
  loggingIn: 'Logging in…',
  passkeyLogin: 'Log in with passkey',
  emailLogin: 'Use an email code',
  email: 'Email',
  emailCode: 'Email code',
  codeSent: (email: string) =>
    `If this address can receive the request, a code will be sent to ${email}. It expires in 10 minutes.`,
  pleaseWait: 'Please wait…',
  verifyContinue: 'Verify and continue',
  sendCode: 'Send code',
  resendIn: (seconds: number) => `Resend in ${seconds} seconds`,
  resend: 'Resend',
  changeEmail: 'Change email',
  captchaTitle: 'Complete the verification',
  cancel: 'Cancel',
  backToPasskey: 'Back to passkey login',
  haveAccount: 'Already have an account?',
  needAccount: 'Need an account?',
  captchaUnconfigured: 'Verification is unavailable. Please try again later.',
  captchaIncomplete: 'Verification was not completed. Please try again.',
  actionFailed: 'Something went wrong. Please try again.',
  loginEmail: 'Login email',
  nickname: 'Nickname',
  nicknameHint:
    '2–24 visible characters: letters from any language, digits 0–9, or underscores. Nicknames must be unique and can be changed once every 30 days.',
  nextNicknameChange: 'Next change available:',
  saveNickname: 'Save nickname',
  nicknameSaved: 'Nickname saved.',
  yourPasskeys: 'Your passkeys',
  passkeyDescription:
    'Add a passkey for devices you use often. You can always log in with an email code.',
  noPasskeys: 'No passkeys added yet.',
  addedOn: 'Added on',
  deletePasskey: (name: string) => `Delete ${name}`,
  delete: 'Delete',
  confirmDelete:
    'Delete this passkey? You can still log in with an email code.',
  passkeyDeleted: 'Passkey deleted.',
  passkeyAdded: 'Passkey added.',
  addPasskey: 'Add passkey',
  errors: {
    INVALID_INPUT: 'Please check your input and try again.',
    UNSUPPORTED_MEDIA_TYPE: 'Please send a JSON request.',
    PAYLOAD_TOO_LARGE: 'The request is too large.',
    NOT_FOUND: 'The requested resource was not found.',
    METHOD_NOT_ALLOWED: 'This request method is not supported.',
    EMAIL_ALREADY_REGISTERED:
      'This email is already registered. Please log in.',
    UNAUTHENTICATED: 'Your session has ended. Please log in again.',
    RATE_LIMITED: 'Too many attempts. Please try again later.',
    OTP_INVALID: 'The code is invalid or has expired.',
    CAPTCHA_FAILED: 'Please complete verification again.',
    PASSKEY_UNAVAILABLE: 'Passkeys are unavailable. Use an email code.',
    AUTH_UNAVAILABLE:
      'Authentication is temporarily unavailable. Please try again later.',
    AUTH_FAILED: 'Authentication failed. Check your details and try again.',
    PROFILE_UNAVAILABLE:
      'Your profile is temporarily unavailable. Please try again later.',
    NICKNAME_TAKEN: 'This nickname is already in use.',
    NICKNAME_COOLDOWN:
      'You can change your nickname again after the waiting period.',
    INVALID_NICKNAME: 'Please enter a valid nickname.',
    PROFILE_CHANGED: 'Your profile changed. Refresh the page and try again.',
    EMAIL_UNVERIFIED: 'Please verify your email first.',
    INVALID_ORIGIN: 'This request came from an invalid origin.',
    CONFIGURATION_REQUIRED:
      'The service is not configured. Please try again later.',
    CAPTCHA_REQUIRED: 'Please complete the verification.',
    SERVICE_UNAVAILABLE:
      'The service is temporarily unavailable. Please try again later.',
    GAME_UNAVAILABLE: 'The game account service is temporarily unavailable.',
    GAME_ACCOUNT_EXISTS: 'This game account name is already in use.',
    PASSKEY_UNSUPPORTED:
      'This device does not support passkeys. Use an email code.',
    PASSKEY_CANCELLED:
      'The passkey request was cancelled or timed out. Try again or use an email code.',
    PASSKEY_FAILED:
      'Could not use the passkey. Try again or use an email code.',
    WEBAUTHN: 'Passkey verification failed. Try again or use an email code.',
  },
};

export type Dictionary = typeof en;

const zh: Dictionary = {
  description: '邮箱注册，使用 Passkey 安全登录。',
  language: '语言',
  accountNavigation: '账号导航',
  accountMenu: '账号菜单',
  login: '登录',
  signup: '注册',
  profile: '个人资料',
  security: '账号与安全',
  gameAccounts: '游戏账号',
  myGameAccounts: '我的账号',
  noGameAccounts: '还没有游戏账号。',
  noGameCharacters: '暂无角色',
  gameCharacterLevel: (level) => `等级 ${level}`,
  createGameAccount: '创建游戏账号',
  gameAccountName: '账号名',
  gameAccountPassword: '密码',
  gameAccountPasswordConfirmation: '确认密码',
  creatingGameAccount: '创建中…',
  logout: '退出登录',
  logoutFailed: '退出失败，请重试。',
  homeEyebrow: '欢迎来到 r2f2',
  homeTitle: '你的 r2f2，从这里开始。',
  homeText:
    '一个连接玩家与精彩内容的空间。我们正在打造新的 r2f2 网站，更多内容即将到来。',
  accountUnavailable: '暂时无法读取账号',
  serviceUnavailable: '服务暂不可用，请稍后重试。',
  loginFirst: '请先登录',
  restoringSession: '正在恢复会话…',
  sessionEnded: '会话已结束，请使用 Passkey 或邮件验证码重新登录。',
  goToLogin: '前往登录',
  emailVerified: '邮箱已验证',
  setUpPasskey: '设置你的 Passkey',
  passkeyNextTime: '下次使用指纹、面容或设备 PIN 即可登录。',
  settingUp: '正在设置…',
  createPasskey: '创建 Passkey',
  later: '稍后设置',
  createAccount: '创建你的账号',
  loggingIn: '正在登录…',
  passkeyLogin: '使用 Passkey 登录',
  emailLogin: '使用邮件验证码',
  email: '邮箱',
  emailCode: '邮件验证码',
  codeSent: (email) =>
    `若邮箱可以接收此请求，验证码将发送到 ${email}，10 分钟内有效。`,
  pleaseWait: '请稍候…',
  verifyContinue: '验证并继续',
  sendCode: '发送验证码',
  resendIn: (seconds) => `${seconds} 秒后可重发`,
  resend: '重新发送',
  changeEmail: '更换邮箱',
  captchaTitle: '请完成人机验证',
  cancel: '取消',
  backToPasskey: '返回 Passkey 登录',
  haveAccount: '已有账号？',
  needAccount: '还没有账号？',
  captchaUnconfigured: '人机验证未配置，请稍后重试。',
  captchaIncomplete: '人机验证未完成，请重试。',
  actionFailed: '操作失败，请重试。',
  loginEmail: '登录邮箱',
  nickname: '昵称',
  nicknameHint:
    '2–24 个可见字符，可使用各种语言的字母、数字 0–9 或下划线。昵称唯一，每 30 天可修改一次。',
  nextNicknameChange: '下次可修改时间：',
  saveNickname: '保存昵称',
  nicknameSaved: '昵称已保存。',
  yourPasskeys: '你的 Passkey',
  passkeyDescription: '为常用设备添加 Passkey。也可以随时使用邮件验证码登录。',
  noPasskeys: '尚未添加 Passkey。',
  addedOn: '添加于',
  deletePasskey: (name) => `删除 ${name}`,
  delete: '删除',
  confirmDelete: '删除此 Passkey？你仍可使用邮件验证码登录。',
  passkeyDeleted: 'Passkey 已删除。',
  passkeyAdded: 'Passkey 已添加。',
  addPasskey: '添加 Passkey',
  errors: {
    INVALID_INPUT: '输入格式无效，请检查后重试。',
    UNSUPPORTED_MEDIA_TYPE: '请使用 JSON 请求。',
    PAYLOAD_TOO_LARGE: '请求过大。',
    NOT_FOUND: '请求的内容不存在。',
    METHOD_NOT_ALLOWED: '请求方法不支持。',
    EMAIL_ALREADY_REGISTERED: '该邮箱已注册，请登录。',
    UNAUTHENTICATED: '会话已失效，请重新登录。',
    RATE_LIMITED: '操作过于频繁，请稍后重试。',
    OTP_INVALID: '验证码无效或已过期。',
    CAPTCHA_FAILED: '请重新完成人机验证。',
    PASSKEY_UNAVAILABLE: 'Passkey 暂不可用，请使用邮件验证码。',
    AUTH_UNAVAILABLE: '认证服务暂不可用，请稍后重试。',
    AUTH_FAILED: '认证失败，请检查输入后重试。',
    PROFILE_UNAVAILABLE: '用户资料暂不可用，请稍后重试。',
    NICKNAME_TAKEN: '昵称已被使用。',
    NICKNAME_COOLDOWN: '昵称尚未到可修改时间。',
    INVALID_NICKNAME: '昵称格式无效。',
    PROFILE_CHANGED: '用户资料已改变，请刷新页面后重试。',
    EMAIL_UNVERIFIED: '请先验证邮箱。',
    INVALID_ORIGIN: '请求来源无效。',
    CONFIGURATION_REQUIRED: '服务尚未配置，请稍后重试。',
    CAPTCHA_REQUIRED: '请完成人机验证。',
    SERVICE_UNAVAILABLE: '服务暂不可用，请稍后重试。',
    GAME_UNAVAILABLE: '游戏账号服务暂不可用。',
    GAME_ACCOUNT_EXISTS: '游戏账号名已被使用。',
    PASSKEY_UNSUPPORTED: '此设备不支持 Passkey，请使用邮件验证码。',
    PASSKEY_CANCELLED: 'Passkey 操作已取消或超时，可以重试或使用邮件验证码。',
    PASSKEY_FAILED: '无法使用 Passkey，请重试或使用邮件验证码。',
    WEBAUTHN: 'Passkey 验证失败，请重试或使用邮件验证码。',
  },
};

const es: Dictionary = {
  description:
    'Regístrate con tu correo e inicia sesión de forma segura con una llave de acceso.',
  language: 'Idioma',
  accountNavigation: 'Navegación de la cuenta',
  accountMenu: 'Menú de la cuenta',
  login: 'Iniciar sesión',
  signup: 'Registrarse',
  profile: 'Perfil',
  security: 'Seguridad de la cuenta',
  gameAccounts: 'Cuentas de juego',
  myGameAccounts: 'Mis cuentas',
  noGameAccounts: 'Aún no tienes cuentas de juego.',
  noGameCharacters: 'Aún no hay personajes',
  gameCharacterLevel: (level) => `Nivel ${level}`,
  createGameAccount: 'Crear cuenta de juego',
  gameAccountName: 'Nombre de cuenta',
  gameAccountPassword: 'Contraseña',
  gameAccountPasswordConfirmation: 'Confirmar contraseña',
  creatingGameAccount: 'Creando…',
  logout: 'Cerrar sesión',
  logoutFailed: 'No se pudo cerrar sesión. Inténtalo de nuevo.',
  homeEyebrow: 'Bienvenido a r2f2',
  homeTitle: 'Tu r2f2 comienza aquí.',
  homeText:
    'Un espacio que conecta a jugadores con contenido interesante. Estamos creando el nuevo sitio de r2f2. Pronto habrá más novedades.',
  accountUnavailable: 'Tu cuenta no está disponible temporalmente',
  serviceUnavailable:
    'El servicio no está disponible temporalmente. Inténtalo más tarde.',
  loginFirst: 'Inicia sesión',
  restoringSession: 'Restaurando tu sesión…',
  sessionEnded:
    'Tu sesión terminó. Vuelve a iniciar sesión con una llave de acceso o un código por correo.',
  goToLogin: 'Ir a iniciar sesión',
  emailVerified: 'Correo verificado',
  setUpPasskey: 'Configura tu llave de acceso',
  passkeyNextTime:
    'La próxima vez, usa tu huella, rostro o PIN del dispositivo para iniciar sesión.',
  settingUp: 'Configurando…',
  createPasskey: 'Crear llave de acceso',
  later: 'Configurar después',
  createAccount: 'Crea tu cuenta',
  loggingIn: 'Iniciando sesión…',
  passkeyLogin: 'Entrar con llave de acceso',
  emailLogin: 'Usar un código por correo',
  email: 'Correo electrónico',
  emailCode: 'Código por correo',
  codeSent: (email) =>
    `Si esta dirección puede recibir la solicitud, enviaremos un código a ${email}. Caduca en 10 minutos.`,
  pleaseWait: 'Espera…',
  verifyContinue: 'Verificar y continuar',
  sendCode: 'Enviar código',
  resendIn: (seconds) => `Reenviar en ${seconds} segundos`,
  resend: 'Reenviar',
  changeEmail: 'Cambiar correo',
  captchaTitle: 'Completa la verificación',
  cancel: 'Cancelar',
  backToPasskey: 'Volver a la llave de acceso',
  haveAccount: '¿Ya tienes una cuenta?',
  needAccount: '¿Necesitas una cuenta?',
  captchaUnconfigured:
    'La verificación no está disponible. Inténtalo más tarde.',
  captchaIncomplete: 'La verificación no se completó. Inténtalo de nuevo.',
  actionFailed: 'Algo salió mal. Inténtalo de nuevo.',
  loginEmail: 'Correo de acceso',
  nickname: 'Apodo',
  nicknameHint:
    'De 2 a 24 caracteres visibles: letras de cualquier idioma, dígitos del 0 al 9 o guiones bajos. Los apodos deben ser únicos y se pueden cambiar cada 30 días.',
  nextNicknameChange: 'Próximo cambio disponible:',
  saveNickname: 'Guardar apodo',
  nicknameSaved: 'Apodo guardado.',
  yourPasskeys: 'Tus llaves de acceso',
  passkeyDescription:
    'Añade una llave de acceso para tus dispositivos habituales. También puedes iniciar sesión con un código por correo.',
  noPasskeys: 'Aún no has añadido llaves de acceso.',
  addedOn: 'Añadida el',
  deletePasskey: (name) => `Eliminar ${name}`,
  delete: 'Eliminar',
  confirmDelete:
    '¿Eliminar esta llave de acceso? Podrás seguir iniciando sesión con un código por correo.',
  passkeyDeleted: 'Llave de acceso eliminada.',
  passkeyAdded: 'Llave de acceso añadida.',
  addPasskey: 'Añadir llave de acceso',
  errors: {
    INVALID_INPUT: 'Revisa los datos e inténtalo de nuevo.',
    UNSUPPORTED_MEDIA_TYPE: 'Envía una solicitud JSON.',
    PAYLOAD_TOO_LARGE: 'La solicitud es demasiado grande.',
    NOT_FOUND: 'No se encontró el recurso solicitado.',
    METHOD_NOT_ALLOWED: 'Este método de solicitud no está permitido.',
    EMAIL_ALREADY_REGISTERED: 'Este correo ya está registrado. Inicia sesión.',
    UNAUTHENTICATED: 'Tu sesión terminó. Vuelve a iniciar sesión.',
    RATE_LIMITED: 'Demasiados intentos. Inténtalo más tarde.',
    OTP_INVALID: 'El código no es válido o ha caducado.',
    CAPTCHA_FAILED: 'Completa la verificación de nuevo.',
    PASSKEY_UNAVAILABLE:
      'Las llaves de acceso no están disponibles. Usa un código por correo.',
    AUTH_UNAVAILABLE:
      'El servicio de autenticación no está disponible temporalmente. Inténtalo más tarde.',
    AUTH_FAILED:
      'Error de autenticación. Revisa los datos e inténtalo de nuevo.',
    PROFILE_UNAVAILABLE:
      'Tu perfil no está disponible temporalmente. Inténtalo más tarde.',
    NICKNAME_TAKEN: 'Este apodo ya está en uso.',
    NICKNAME_COOLDOWN:
      'Podrás cambiar tu apodo al terminar el período de espera.',
    INVALID_NICKNAME: 'Introduce un apodo válido.',
    PROFILE_CHANGED:
      'Tu perfil cambió. Actualiza la página e inténtalo de nuevo.',
    EMAIL_UNVERIFIED: 'Verifica primero tu correo electrónico.',
    INVALID_ORIGIN: 'La solicitud proviene de un origen no válido.',
    CONFIGURATION_REQUIRED:
      'El servicio no está configurado. Inténtalo más tarde.',
    CAPTCHA_REQUIRED: 'Completa la verificación.',
    SERVICE_UNAVAILABLE:
      'El servicio no está disponible temporalmente. Inténtalo más tarde.',
    GAME_UNAVAILABLE: 'El servicio de cuentas de juego no está disponible temporalmente.',
    GAME_ACCOUNT_EXISTS: 'Este nombre de cuenta de juego ya está en uso.',
    PASSKEY_UNSUPPORTED:
      'Este dispositivo no admite llaves de acceso. Usa un código por correo.',
    PASSKEY_CANCELLED:
      'La solicitud se canceló o caducó. Inténtalo de nuevo o usa un código por correo.',
    PASSKEY_FAILED:
      'No se pudo usar la llave de acceso. Inténtalo de nuevo o usa un código por correo.',
    WEBAUTHN:
      'Falló la verificación de la llave de acceso. Inténtalo de nuevo o usa un código por correo.',
  },
};

export const dictionaries: Record<Locale, Dictionary> = { en, 'zh-CN': zh, es };
export function formatDate(locale: Locale, value: string): string {
  return new Intl.DateTimeFormat(locale).format(new Date(value));
}
export function formatDateTime(locale: Locale, value: string): string {
  return new Intl.DateTimeFormat(locale, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(value));
}
export function errorMessage(dictionary: Dictionary, error: unknown): string {
  if (typeof error !== 'object' || error === null || !('code' in error))
    return dictionary.actionFailed;
  const code = (error as { code?: unknown }).code;
  if (typeof code !== 'string') return dictionary.actionFailed;
  if (code.startsWith('WEBAUTHN_')) return dictionary.errors.WEBAUTHN;
  return (
    dictionary.errors[code as keyof Dictionary['errors']] ??
    dictionary.actionFailed
  );
}
