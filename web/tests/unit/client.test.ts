import { describe, expect, it, vi } from 'vitest';
import { ApiError, createClient } from '$lib/api/client';
import { BECKY_ID, TONY_ID } from '$lib/api/types';

function jsonResponse(status: number, body: unknown): Response {
	return new Response(JSON.stringify(body), {
		status,
		headers: { 'Content-Type': 'application/json' }
	});
}

describe('createClient', () => {
	it('lists archers from GET /api/v1/archers', async () => {
		const fetchImpl = vi.fn(async (input: RequestInfo | URL) => {
			expect(String(input)).toBe('/api/v1/archers');
			return jsonResponse(200, [
				{ id: BECKY_ID, displayName: 'Becky' },
				{ id: TONY_ID, displayName: 'Tony' }
			]);
		});
		const client = createClient(fetchImpl as typeof fetch);
		const archers = await client.listArchers();
		expect(archers[0].displayName).toBe('Becky');
	});

	it('posts an arrow scoreCode and returns the API summary', async () => {
		const fetchImpl = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
			expect(String(input)).toBe('/api/v1/sessions/sess-1/arrows');
			expect(init?.method).toBe('POST');
			expect(JSON.parse(String(init?.body))).toEqual({ scoreCode: 'X' });
			return jsonResponse(201, {
				arrow: {
					id: 'a1',
					arrowNumber: 1,
					scoreCode: 'X',
					scoreValue: 10,
					recordedAt: '2026-10-06T12:00:00Z'
				},
				endNumber: 1,
				summary: { total: 10, hits: 1, golds: 1, xCount: 1, arrowCount: 1 }
			});
		});
		const client = createClient(fetchImpl as typeof fetch);
		const got = await client.recordArrow('sess-1', 'X');
		expect(got.summary.total).toBe(10);
		expect(got.endNumber).toBe(1);
	});

	it('maps error JSON to ApiError', async () => {
		const fetchImpl = vi.fn(async () =>
			jsonResponse(409, { code: 'session_completed', message: 'Session is already completed' })
		);
		const client = createClient(fetchImpl as typeof fetch);
		await expect(client.completeSession('sess-1')).rejects.toMatchObject({
			name: 'ApiError',
			status: 409,
			code: 'session_completed'
		} satisfies Partial<ApiError>);
	});
});
