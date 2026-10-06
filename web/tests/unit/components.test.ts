import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import ArcherPicker from '$lib/components/ArcherPicker.svelte';
import HomeScreen from '$lib/components/HomeScreen.svelte';
import LivePad from '$lib/components/LivePad.svelte';
import HistoryList from '$lib/components/HistoryList.svelte';
import Scorecard from '$lib/components/Scorecard.svelte';
import { BECKY_ID, TONY_ID, type Session } from '$lib/api/types';

afterEach(() => {
	cleanup();
});

describe('ArcherPicker', () => {
	it('offers Tony and Becky and reports the picked id', async () => {
		const user = userEvent.setup();
		const onPick = vi.fn();
		render(ArcherPicker, {
			props: {
				archers: [
					{ id: TONY_ID, displayName: 'Tony' },
					{ id: BECKY_ID, displayName: 'Becky' }
				],
				onPick
			}
		});
		expect(screen.getByText('Who is shooting?')).toBeTruthy();
		await user.click(screen.getByRole('button', { name: 'Becky' }));
		expect(onPick).toHaveBeenCalledWith(BECKY_ID);
	});
});

describe('HomeScreen', () => {
	it('shows the PB card and starts a session', async () => {
		const user = userEvent.setup();
		const onStart = vi.fn();
		render(HomeScreen, {
			props: {
				displayName: 'Becky',
				pb: { headline: 'PB · 60 arrows · 498', hint: null },
				sessions: [
					{
						id: 's1',
						dateLabel: 'Yesterday',
						arrowCount: 60,
						total: 498,
						isPersonalBest: true
					}
				],
				onStart,
				onSwitch: vi.fn(),
				onOpenSession: vi.fn(),
				onHistory: vi.fn()
			}
		});
		expect(screen.getByText('PB · 60 arrows · 498')).toBeTruthy();
		await user.click(screen.getByRole('button', { name: 'Start session' }));
		expect(onStart).toHaveBeenCalled();
	});
});

describe('LivePad', () => {
	it('posts pad codes and disables finish until an arrow exists', async () => {
		const user = userEvent.setup();
		const onScore = vi.fn();
		const onFinish = vi.fn();
		render(LivePad, {
			props: {
				view: {
					total: 0,
					hits: 0,
					golds: 0,
					xCount: 0,
					progress: 'End 1 · 0/6',
					slots: [null, null, null, null, null, null]
				},
				canFinish: false,
				onScore,
				onFinish
			}
		});
		expect((screen.getByRole('button', { name: 'Finish session' }) as HTMLButtonElement).disabled).toBe(
			true
		);
		await user.click(screen.getByRole('button', { name: 'X' }));
		expect(onScore).toHaveBeenCalledWith('X');
	});
});

describe('HistoryList', () => {
	it('opens a session from a row', async () => {
		const user = userEvent.setup();
		const onOpen = vi.fn();
		render(HistoryList, {
			props: {
				sessions: [
					{
						id: 's1',
						dateLabel: '5 Oct 2026',
						arrowCount: 60,
						total: 498,
						isPersonalBest: true
					}
				],
				onOpen,
				onHome: vi.fn()
			}
		});
		await user.click(screen.getByRole('button', { name: /5 Oct 2026/ }));
		expect(onOpen).toHaveBeenCalledWith('s1');
	});
});

describe('Scorecard', () => {
	it('offers Resume when the session is in progress', async () => {
		const user = userEvent.setup();
		const onResume = vi.fn();
		const session: Session = {
			id: 's1',
			archerId: BECKY_ID,
			status: 'in_progress',
			startedAt: '2026-10-05T10:00:00Z',
			completedAt: null,
			summary: { total: 19, hits: 2, golds: 2, xCount: 1, arrowCount: 2 },
			arrowsPerEnd: 6,
			ends: [
				{
					endNumber: 1,
					arrows: [
						{
							id: '1',
							arrowNumber: 1,
							scoreCode: 'X',
							scoreValue: 10,
							recordedAt: '2026-10-05T10:00:00Z'
						},
						{
							id: '2',
							arrowNumber: 2,
							scoreCode: '9',
							scoreValue: 9,
							recordedAt: '2026-10-05T10:01:00Z'
						}
					]
				}
			]
		};
		render(Scorecard, { props: { session, onResume } });
		expect(screen.getByText(/Total 19/)).toBeTruthy();
		await user.click(screen.getByRole('button', { name: 'Resume' }));
		expect(onResume).toHaveBeenCalled();
	});
});
