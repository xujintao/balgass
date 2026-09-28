import { context } from './session';
import { me } from './profile';
import { ApiError } from './errors';
export async function pageProfile() {
  try {
    return { kind: 'ok' as const, profile: await me(await context()) };
  } catch (error) {
    if (error instanceof ApiError && error.status === 401)
      return { kind: 'unauthenticated' as const };
    return {
      kind: 'error' as const,
      message:
        error instanceof ApiError
          ? error.message
          : '服务暂不可用，请稍后重试。',
    };
  }
}
