export const excellentOptions = [
  'excellent_attack_rate', 'excellent_attack_level', 'excellent_attack_percent',
  'excellent_attack_speed', 'excellent_attack_hp', 'excellent_attack_mp',
] as const;
export type ExcellentOption = (typeof excellentOptions)[number];
