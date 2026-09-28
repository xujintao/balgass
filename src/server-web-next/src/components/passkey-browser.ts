'use client';
import {
  startAuthentication,
  startRegistration,
  browserSupportsWebAuthn,
} from '@simplewebauthn/browser';
import type {
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
} from '@simplewebauthn/browser';
import { api, ClientError } from './client-api';
export async function performPasskey(
  mode: 'login' | 'register',
  captchaToken?: string,
) {
  if (!browserSupportsWebAuthn())
    throw new ClientError(
      'PASSKEY_UNSUPPORTED',
      '此设备不支持 Passkey，请使用邮件验证码。',
    );
  try {
    const result = await api<{
      challenge_id: string;
      options: PublicKeyCredentialCreationOptionsJSON &
        PublicKeyCredentialRequestOptionsJSON;
    }>(
      `auth/passkeys/${mode}/options`,
      'POST',
      mode === 'login' ? { captchaToken: captchaToken || undefined } : {},
    );
    const credential =
      mode === 'login'
        ? await startAuthentication({ optionsJSON: result.options })
        : await startRegistration({ optionsJSON: result.options });
    return await api(`auth/passkeys/${mode}/verify`, 'POST', {
      challengeId: result.challenge_id,
      credential,
    });
  } catch (error) {
    if (error instanceof ClientError) throw error;
    if (
      error instanceof Error &&
      ['NotAllowedError', 'AbortError', 'ERROR_CEREMONY_ABORTED'].includes(
        error.name,
      )
    )
      throw new ClientError(
        'PASSKEY_CANCELLED',
        'Passkey 操作已取消或超时，可以重试或使用邮件验证码。',
      );
    throw new ClientError(
      'PASSKEY_FAILED',
      '无法使用 Passkey，请重试或使用邮件验证码。',
    );
  }
}
