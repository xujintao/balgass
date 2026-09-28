import { describe, expect, it } from 'vitest';
import {
  nickname,
  otpSend,
  otpVerify,
  authenticationVerify,
} from '../../src/lib/validation';
describe('输入规则', () => {
  it('去除空白并保留合法中文与大小写', () =>
    expect(nickname.parse('  玩家_Ab12\u3000')).toBe('玩家_Ab12'));
  it.each(['a', 'a'.repeat(25), 'hello world', '小猫🐈', 'a-b', 'éé', 'ＡＢ'])(
    '拒绝非法昵称 %s',
    (value) => expect(nickname.safeParse(value).success).toBe(false),
  );
  it('按 Unicode 字符计数', () => expect(nickname.parse('𠀀𠀀')).toBe('𠀀𠀀'));
  it('邮箱大小写归一化，发送必须声明用途', () => {
    expect(
      otpSend.parse({ email: 'Me@Example.com', intent: 'login' }).email,
    ).toBe('me@example.com');
    expect(otpSend.safeParse({ email: 'me@example.com' }).success).toBe(false);
  });
  it('验证码必须为 6 位数字', () => {
    expect(
      otpVerify.safeParse({ email: 'me@example.com', code: '012345' }).success,
    ).toBe(true);
    expect(
      otpVerify.safeParse({ email: 'me@example.com', code: '12345a' }).success,
    ).toBe(false);
  });
  it('不接受畸形 WebAuthn 提交', () =>
    expect(
      authenticationVerify.safeParse({ challengeId: 'bad', credential: {} })
        .success,
    ).toBe(false));
});
