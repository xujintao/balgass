import type { Locale } from './i18n';

export const adminText: Record<Locale, {
  title: string; email: string; nickname: string; role: string;
  admin: string; player: string; previous: string; next: string;
  page: (current: number, total: number) => string;
  saved: string; unavailable: string;
  invalidPage: string; lastAdmin: string;
}> = {
  en: {
    title: 'Users', email: 'Email', nickname: 'Nickname', role: 'Role',
    admin: 'Admin', player: 'Player', previous: 'Previous', next: 'Next',
    page: (current, total) => `Page ${current} of ${total}`,
    saved: 'Role updated.',
    unavailable: 'Users are temporarily unavailable.', invalidPage: 'Invalid page number.',
    lastAdmin: 'The last admin cannot be removed.',
  },
  'zh-CN': {
    title: '用户管理', email: '邮箱', nickname: '昵称', role: '角色',
    admin: '管理员', player: '普通用户', previous: '上一页', next: '下一页',
    page: (current, total) => `第 ${current} / ${total} 页`,
    saved: '角色已更新。',
    unavailable: '用户列表暂不可用。', invalidPage: '页码无效。',
    lastAdmin: '不能撤销最后一位管理员。',
  },
  es: {
    title: 'Usuarios', email: 'Correo', nickname: 'Apodo', role: 'Rol',
    admin: 'Administrador', player: 'Jugador', previous: 'Anterior', next: 'Siguiente',
    page: (current, total) => `Página ${current} de ${total}`,
    saved: 'Rol actualizado.',
    unavailable: 'La lista de usuarios no está disponible.', invalidPage: 'Número de página inválido.',
    lastAdmin: 'No se puede quitar al último administrador.',
  },
};
