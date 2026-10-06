import { describe, expect, it } from 'vitest';
import type { PersonalBest, Session, SessionListItem } from '$lib/api/types';
import { livePadView, personalBestCard, sessionListRows, lastFive, endRow } from '$lib/session/mappers';

const emptySummary = { total: 0, hits: 0, golds: 0, xCount: 0, arrowCount: 0 };

function session(partial: Partial<Session> & Pick<Session, 'id'>): Session {
	return {
		archerId: 'a',
		status: 'in_progress',
		startedAt: '2026-10-06T12:00:00Z',
		completedAt: null,
		summary: emptySummary,
		arrowsPerEnd: 6,
		ends: [],
		...partial
	};
}

describe('livePadView', () => {
	it('shows End 1 · 0/6 with empty slots when there are no arrows', () => {
		const view = livePadView(session({ id: 's1' }));
		expect(view.progress).toBe('End 1 · 0/6');
		expect(view.slots).toEqual([null, null, null, null, null, null]);
		expect(view.total).toBe(0);
	});

	it('fills current-end slots from API codes and uses the API summary total', () => {
		const view = livePadView(
			session({
				id: 's1',
				summary: { total: 147, hits: 14, golds: 6, xCount: 2, arrowCount: 15 },
				ends: [
					{
						endNumber: 2,
						arrows: [
							{
								id: '1',
								arrowNumber: 1,
								scoreCode: 'X',
								scoreValue: 10,
								recordedAt: '2026-10-06T12:00:00Z'
							},
							{
								id: '2',
								arrowNumber: 2,
								scoreCode: '10',
								scoreValue: 10,
								recordedAt: '2026-10-06T12:01:00Z'
							},
							{
								id: '3',
								arrowNumber: 3,
								scoreCode: '9',
								scoreValue: 9,
								recordedAt: '2026-10-06T12:02:00Z'
							}
						]
					}
				]
			})
		);
		expect(view.total).toBe(147);
		expect(view.hits).toBe(14);
		expect(view.progress).toBe('End 2 · 3/6');
		expect(view.slots).toEqual(['X', '10', '9', null, null, null]);
	});

	it('opens the next end after six arrows', () => {
		const arrows = Array.from({ length: 6 }, (_, i) => ({
			id: String(i),
			arrowNumber: i + 1,
			scoreCode: '9' as const,
			scoreValue: 9,
			recordedAt: '2026-10-06T12:00:00Z'
		}));
		const view = livePadView(session({ id: 's1', ends: [{ endNumber: 1, arrows }] }));
		expect(view.progress).toBe('End 2 · 0/6');
		expect(view.slots).toEqual([null, null, null, null, null, null]);
	});
});

describe('sessionListRows', () => {
	it('labels yesterday, totals from summary, and stars the PB session for that length', () => {
		const now = new Date('2026-10-06T15:00:00Z');
		const sessions: SessionListItem[] = [
			{
				id: 'pb-60',
				archerId: 'a',
				status: 'completed',
				startedAt: '2026-10-05T10:00:00Z',
				completedAt: '2026-10-05T11:00:00Z',
				summary: { total: 498, hits: 58, golds: 22, xCount: 4, arrowCount: 60 }
			},
			{
				id: 'other',
				archerId: 'a',
				status: 'completed',
				startedAt: '2026-09-28T10:00:00Z',
				completedAt: '2026-09-28T11:00:00Z',
				summary: { total: 220, hits: 28, golds: 4, xCount: 0, arrowCount: 30 }
			}
		];
		const bests: PersonalBest[] = [
			{
				archerId: 'a',
				sessionId: 'pb-60',
				total: 498,
				arrowCount: 60,
				achievedAt: '2026-10-05T11:00:00Z'
			}
		];
		const rows = sessionListRows(sessions, bests, now);
		expect(rows[0]).toMatchObject({
			id: 'pb-60',
			dateLabel: 'Yesterday',
			arrowCount: 60,
			total: 498,
			isPersonalBest: true
		});
		expect(rows[1].isPersonalBest).toBe(false);
		expect(rows[1].total).toBe(220);
		expect(rows[1].dateLabel).toMatch(/28 Sep/);
	});
});

describe('personalBestCard', () => {
	it('shows empty copy when there are no completed bests', () => {
		expect(personalBestCard([])).toEqual({ headline: 'No personal best yet', hint: null });
	});

	it('uses the most recently achieved bucket and hints when more exist', () => {
		const card = personalBestCard([
			{
				archerId: 'a',
				sessionId: 'old',
				total: 220,
				arrowCount: 30,
				achievedAt: '2026-09-01T00:00:00Z'
			},
			{
				archerId: 'a',
				sessionId: 'new',
				total: 498,
				arrowCount: 60,
				achievedAt: '2026-10-05T00:00:00Z'
			}
		]);
		expect(card.headline).toBe('PB · 60 arrows · 498');
		expect(card.hint).toBe('More personal bests are on History');
	});
});

describe('lastFive', () => {
	it('keeps the first five of a newest-first list', () => {
		expect(lastFive([1, 2, 3, 4, 5, 6])).toEqual([1, 2, 3, 4, 5]);
	});
});

describe('endRow', () => {
	it('sums scoreValue from the API rather than rescoring codes', () => {
		expect(
			endRow({
				endNumber: 1,
				arrows: [
					{ scoreCode: 'X', scoreValue: 10 },
					{ scoreCode: '9', scoreValue: 9 }
				]
			})
		).toEqual({ endNumber: 1, codes: ['X', '9'], total: 19 });
	});
});
