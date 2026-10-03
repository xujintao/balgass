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
    throw new ClientError('PASSKEY_UNSUPPORTED', 'PASSKEY_UNSUPPORTED');
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
      throw new ClientError('PASSKEY_CANCELLED', 'PASSKEY_CANCELLED');
    throw new ClientError('PASSKEY_FAILED', 'PASSKEY_FAILED');
  }
}
