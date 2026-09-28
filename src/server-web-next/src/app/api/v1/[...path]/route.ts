import { handle } from '@/lib/api';
export const dynamic = 'force-dynamic';
async function route(
  request: Request,
  ctx: { params: Promise<{ path: string[] }> },
) {
  return handle(request, (await ctx.params).path.join('/'));
}
export { route as GET, route as POST, route as PATCH, route as DELETE };
