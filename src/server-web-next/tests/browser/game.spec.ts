import { expect, test } from '@playwright/test';

test('实时地图订阅、切换和绘制', async ({ page, context }) => {
  await context.addCookies([
    { name: 'r2f2-locale', value: 'zh-CN', url: 'http://localhost:3000' },
  ]);
  await page.addInitScript(`
    window.__gameMessages = [];
    class MockWebSocket {
      static OPEN = 1;
      readyState = 0;
      constructor(url) {
        this.url = url;
        queueMicrotask(() => {
          this.readyState = 1;
          this.onopen?.();
        });
      }
      send(data) {
        const message = JSON.parse(data);
        window.__gameMessages.push(message);
        const name = message.in.name;
        for (const out of [
          { name: 'map-name', data: name },
          { name: 'map-data', data: [{ x: 10, y: 10 }] },
          { name: 'object', data: {
            players: [{ x: 20, y: 20 }],
            monsters: [{ x: 30, y: 30 }],
            npcs: [{ x: 40, y: 40 }],
            stands: [],
          } },
        ]) this.onmessage?.({ data: JSON.stringify({ action: 'SubscribeMapReply', out }) });
      }
      close() {
        this.readyState = 3;
        this.onclose?.();
      }
    }
    window.WebSocket = MockWebSocket;
  `);

  await page.goto('/game');
  await expect(page).toHaveTitle('实时地图 · r2f2');
  await expect(page.getByRole('heading', { name: '实时地图' })).toBeVisible();
  await expect(page.getByRole('status')).toHaveText('已连接');
  await expect(page.getByText('玩家: 1')).toBeVisible();
  await expect(page.getByRole('combobox', { name: '地图' }).locator('option')).toHaveCount(4);
  expect(await page.evaluate(() => (window as unknown as { __gameMessages: { in: { name: string } }[] }).__gameMessages[0].in.name)).toBe('Lorencia');

  const nonBackgroundPixels = await page.locator('canvas').evaluate((canvas) => {
    const context = (canvas as HTMLCanvasElement).getContext('2d')!;
    const data = context.getImageData(0, 0, 512, 512).data;
    let count = 0;
    for (let i = 0; i < data.length; i += 4) {
      if (data[i] !== 248 || data[i + 1] !== 250 || data[i + 2] !== 247) count++;
    }
    return count;
  });
  expect(nonBackgroundPixels).toBeGreaterThan(20);

  await page.getByRole('combobox', { name: '地图' }).selectOption('Dungeon');
  await expect(page.locator('canvas')).toHaveAttribute('aria-label', 'Dungeon');
  expect(await page.evaluate(() => (window as unknown as { __gameMessages: { in: { name: string } }[] }).__gameMessages.at(-1)?.in.name)).toBe('Dungeon');
});
