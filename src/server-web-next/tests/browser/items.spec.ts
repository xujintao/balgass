import { expect, test } from '@playwright/test';

test('道具目录分类、详情和未登录订单入口', async ({ page, context, baseURL }) => {
  await context.addCookies([{ name: 'r2f2-locale', value: 'zh-CN', url: baseURL! }]);
  await page.goto('/items');
  await expect(page).toHaveTitle('道具商城 · r2f2');
  await expect(page.getByRole('heading', { name: '道具商城' })).toBeVisible();
  await expect(page.getByText(/^(测试剑|短剑)$/).first()).toBeVisible();
  const groups = page.getByRole('navigation', { name: '道具商城' }).locator(':scope > div');
  await expect(groups).toHaveCount(2);
  await expect(groups.nth(0).getByRole('link')).toHaveCount(12);
  await expect(groups.nth(1).getByRole('link')).toHaveCount(9);
  await expect(groups.nth(0).getByRole('link').first()).toHaveAttribute('href', '/items?kind=sword');
  await expect(groups.nth(1).getByRole('link').first()).toHaveAttribute('href', '/items?kind=helmet');
  await page.locator('details').first().locator('summary').click();
  await expect(page.getByRole('table')).toBeVisible();
  await expect(page.getByRole('button', { name: '提交订单' })).toBeVisible();
  await page.getByRole('link', { name: '戒指' }).click();
  await page.locator('details').first().locator('summary').click();
  await expect(page.getByRole('button', { name: '提交订单' })).toBeVisible();
  await page.getByRole('link', { name: '套装' }).click();
  await expect(page.locator('article').first()).toBeVisible();
  await page.goto('/orders');
  await expect(page.getByRole('heading', { name: '请先登录' })).toBeVisible();
});
