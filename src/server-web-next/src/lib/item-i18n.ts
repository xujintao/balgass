import type { Locale } from './i18n';
import type { ItemKind } from './item-catalog';
import type { ExcellentOption } from './item-options';

const en = {
  items: 'Items', orders: 'My orders', unavailable: 'The item catalog is temporarily unavailable.',
  noItems: 'No items in this category.', noOrders: 'No orders yet.', setPieces: 'Set pieces',
  level: 'Level', additional: 'Additional', excellent: 'Excellent options',
  skill: 'Skill', twoHand: 'Two handed', attackSpeed: 'Attack speed', wingKind: 'Wing kind',
  attack: 'Attack', magicPower: 'Magic power', defense: 'Defense', defenseRate: 'Defense rate',
  requirements: 'Strength / dexterity required', yes: 'Yes', no: 'No',
  submit: 'Submit order', submitting: 'Submitting…',
  submitted: 'Order submitted. It is pending.',
  pendingNotice: 'Submitting creates a pending order. No payment or delivery happens here.',
  order: 'Order', status: 'Status', pending: 'Pending',
  categories: {
    sword: 'Swords', claw: 'Claws', axe: 'Axes', mace: 'Maces', scepter: 'Scepters',
    spear: 'Spears', bow: 'Bows', crossbow: 'Crossbows', staff: 'Staves',
    stick: 'Sticks', book: 'Books', shield: 'Shields', helmet: 'Helmets',
    armor: 'Armor', pants: 'Pants', gloves: 'Gloves', boots: 'Boots',
    wing: 'Wings', pendant: 'Pendants', ring: 'Rings', set: 'Sets',
  } satisfies Record<ItemKind, string>,
  excellentLabels: {
    excellent_attack_rate: 'Excellent damage chance',
    excellent_attack_level: 'Attack per level',
    excellent_attack_percent: 'Attack increase',
    excellent_attack_speed: 'Attack speed increase',
    excellent_attack_hp: 'Life from monsters',
    excellent_attack_mp: 'Mana from monsters',
  } satisfies Record<ExcellentOption, string>,
  errors: {
    ITEM_NOT_FOUND: 'The item was not found.',
    ITEM_CATALOG_UNAVAILABLE: 'The item catalog is temporarily unavailable.',
    ORDERS_UNAVAILABLE: 'Orders are temporarily unavailable.',
  },
};
type ItemDictionary = typeof en;
const zh: ItemDictionary = {
  items: '道具商城', orders: '我的订单', unavailable: '道具目录暂不可用。',
  noItems: '该分类暂无道具。', noOrders: '暂无订单。', setPieces: '套装部件',
  level: '等级', additional: '追加', excellent: '卓越属性',
  skill: '技能', twoHand: '双手', attackSpeed: '攻击速度', wingKind: '翅膀类型',
  attack: '攻击力', magicPower: '魔法攻击力', defense: '防御力', defenseRate: '防御成功率',
  requirements: '力量／敏捷需求', yes: '是', no: '否',
  submit: '提交订单', submitting: '提交中…',
  submitted: '订单已提交，当前为待处理状态。',
  pendingNotice: '提交后会生成待处理订单；此处不付款，也不会自动发货。',
  order: '订单', status: '状态', pending: '待处理',
  categories: {
    sword: '剑', claw: '爪', axe: '斧', mace: '锤', scepter: '权杖',
    spear: '矛', bow: '弓', crossbow: '弩', staff: '法杖',
    stick: '短杖', book: '魔法书', shield: '盾牌', helmet: '头盔',
    armor: '铠甲', pants: '护腿', gloves: '手套', boots: '靴子',
    wing: '翅膀', pendant: '项链', ring: '戒指', set: '套装',
  },
  excellentLabels: {
    excellent_attack_rate: '卓越伤害几率',
    excellent_attack_level: '按等级增加攻击力',
    excellent_attack_percent: '增加攻击力',
    excellent_attack_speed: '增加攻击速度',
    excellent_attack_hp: '打怪获得生命值',
    excellent_attack_mp: '打怪获得魔法值',
  },
  errors: {
    ITEM_NOT_FOUND: '道具不存在。',
    ITEM_CATALOG_UNAVAILABLE: '道具目录暂不可用。',
    ORDERS_UNAVAILABLE: '订单暂不可用。',
  },
};
const es: ItemDictionary = {
  items: 'Objetos', orders: 'Mis pedidos', unavailable: 'El catálogo de objetos no está disponible temporalmente.',
  noItems: 'No hay objetos en esta categoría.', noOrders: 'Aún no hay pedidos.', setPieces: 'Piezas del conjunto',
  level: 'Nivel', additional: 'Adicional', excellent: 'Opciones excelentes',
  skill: 'Habilidad', twoHand: 'Dos manos', attackSpeed: 'Velocidad de ataque', wingKind: 'Tipo de alas',
  attack: 'Ataque', magicPower: 'Poder mágico', defense: 'Defensa', defenseRate: 'Tasa de defensa',
  requirements: 'Fuerza / destreza requerida', yes: 'Sí', no: 'No',
  submit: 'Enviar pedido', submitting: 'Enviando…',
  submitted: 'Pedido enviado. Está pendiente.',
  pendingNotice: 'Esto crea un pedido pendiente. Aquí no se paga ni se entrega.',
  order: 'Pedido', status: 'Estado', pending: 'Pendiente',
  categories: {
    sword: 'Espadas', claw: 'Garras', axe: 'Hachas', mace: 'Mazas', scepter: 'Cetros',
    spear: 'Lanzas', bow: 'Arcos', crossbow: 'Ballestas', staff: 'Báculos',
    stick: 'Varas', book: 'Libros', shield: 'Escudos', helmet: 'Cascos',
    armor: 'Armaduras', pants: 'Pantalones', gloves: 'Guantes', boots: 'Botas',
    wing: 'Alas', pendant: 'Colgantes', ring: 'Anillos', set: 'Conjuntos',
  },
  excellentLabels: {
    excellent_attack_rate: 'Probabilidad de daño excelente',
    excellent_attack_level: 'Ataque por nivel',
    excellent_attack_percent: 'Aumento de ataque',
    excellent_attack_speed: 'Aumento de velocidad de ataque',
    excellent_attack_hp: 'Vida obtenida de monstruos',
    excellent_attack_mp: 'Maná obtenido de monstruos',
  },
  errors: {
    ITEM_NOT_FOUND: 'No se encontró el objeto.',
    ITEM_CATALOG_UNAVAILABLE: 'El catálogo de objetos no está disponible temporalmente.',
    ORDERS_UNAVAILABLE: 'Los pedidos no están disponibles temporalmente.',
  },
};
export const itemDictionaries: Record<Locale, ItemDictionary> = { en, 'zh-CN': zh, es };
