import { z } from 'zod';
const nicknameSegments = new Intl.Segmenter('und', { granularity: 'grapheme' });
const nicknameCharacters = /^(?:\p{L}\p{M}*|[0-9_])+$/u;
export const nickname = z
  .string()
  .transform((v) => v.trim().normalize('NFC'))
  .refine((v) => {
    if (!nicknameCharacters.test(v)) return false;
    const length = [...nicknameSegments.segment(v)].length;
    return length >= 2 && length <= 24;
  }, '昵称需为 2–24 个可见字符，只能使用各语言字母、数字或下划线');
export const email = z
  .email()
  .max(254)
  .transform((v) => v.toLowerCase());
export const otpSend = z
  .object({
    email,
    intent: z.enum(['signup', 'login']),
    captchaToken: z.string().min(1).max(4096).optional(),
  })
  .strict();
export const otpVerify = z
  .object({ email, code: z.string().regex(/^\d{6}$/) })
  .strict();
export const refreshInput = z
  .object({ refreshToken: z.string().min(1).max(8192) })
  .strict();
export const nicknameInput = z.object({ nickname }).strict();
export const passkeyStart = z
  .object({ captchaToken: z.string().min(1).max(4096).optional() })
  .strict();
const encoded = z
  .string()
  .min(1)
  .max(20000)
  .regex(/^[A-Za-z0-9_-]+$/);
const credentialBase = z.object({
  id: encoded,
  rawId: encoded,
  type: z.literal('public-key'),
  clientExtensionResults: z.record(z.string(), z.unknown()).default({}),
  authenticatorAttachment: z.enum(['platform', 'cross-platform']).optional(),
});
export const authenticationVerify = z
  .object({
    challengeId: z.uuid(),
    credential: credentialBase
      .extend({
        response: z
          .object({
            clientDataJSON: encoded,
            authenticatorData: encoded,
            signature: encoded,
            userHandle: encoded.optional(),
          })
          .passthrough(),
      })
      .passthrough(),
  })
  .strict();
export const registrationVerify = z
  .object({
    challengeId: z.uuid(),
    credential: credentialBase
      .extend({
        response: z
          .object({
            clientDataJSON: encoded,
            attestationObject: encoded,
            transports: z
              .array(
                z.enum([
                  'ble',
                  'hybrid',
                  'internal',
                  'nfc',
                  'usb',
                  'smart-card',
                ]),
              )
              .optional(),
          })
          .passthrough(),
      })
      .passthrough(),
  })
  .strict();
