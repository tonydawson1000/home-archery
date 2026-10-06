import type { PersonalBest, ScoreSummary, Session, SessionListItem } from '$lib/api/types';

export type LivePadView = {
	total: number;
	hits: number;
	golds: number;
	xCount: number;
	progress: string;
	slots: Array<string | null>;
};

export function livePadView(session: Session): LivePadView {
	const { summary } = session;
	const current = currentEnd(session);
	const filled = current?.arrows.length ?? 0;
	const endNumber = current?.endNumber ?? 1;
	const slots: Array<string | null> = [null, null, null, null, null, null];
	if (current) {
		for (const arrow of current.arrows) {
			slots[arrow.arrowNumber - 1] = arrow.scoreCode;
		}
	}
	return {
		total: summary.total,
		hits: summary.hits,
		golds: summary.golds,
		xCount: summary.xCount,
		progress: `End ${endNumber} · ${filled}/6`,
		slots
	};
}

function currentEnd(session: Session) {
	if (session.ends.length === 0) {
		return undefined;
	}
	const last = session.ends[session.ends.length - 1];
	if (last.arrows.length >= session.arrowsPerEnd) {
		return { endNumber: last.endNumber + 1, arrows: [] as Session['ends'][0]['arrows'] };
	}
	return last;
}

export type SessionListRow = {
	id: string;
	dateLabel: string;
	arrowCount: number;
	total: number;
	isPersonalBest: boolean;
};

export function sessionListRows(
	sessions: SessionListItem[],
	bests: PersonalBest[],
	now = new Date()
): SessionListRow[] {
	const pbBySession = new Set(bests.map((b) => b.sessionId));
	return sessions.map((session) => ({
		id: session.id,
		dateLabel: formatSessionDate(session.startedAt, now),
		arrowCount: session.summary.arrowCount,
		total: session.summary.total,
		isPersonalBest: pbBySession.has(session.id)
	}));
}

export function formatSessionDate(iso: string, now = new Date()): string {
	const started = new Date(iso);
	if (sameCalendarDay(started, now)) {
		return 'Today';
	}
	const yesterday = new Date(now);
	yesterday.setDate(now.getDate() - 1);
	if (sameCalendarDay(started, yesterday)) {
		return 'Yesterday';
	}
	return new Intl.DateTimeFormat('en-GB', { day: 'numeric', month: 'short', year: 'numeric' }).format(
		started
	);
}

function sameCalendarDay(a: Date, b: Date): boolean {
	return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate();
}

export type PersonalBestCard = {
	headline: string;
	hint: string | null;
};

export function personalBestCard(bests: PersonalBest[]): PersonalBestCard {
	if (bests.length === 0) {
		return { headline: 'No personal best yet', hint: null };
	}
	const latest = [...bests].sort(
		(a, b) => new Date(b.achievedAt).getTime() - new Date(a.achievedAt).getTime()
	)[0];
	return {
		headline: `PB · ${latest.arrowCount} arrows · ${latest.total}`,
		hint: bests.length > 1 ? 'More personal bests are on History' : null
	};
}

export function lastFive<T>(items: T[]): T[] {
	return items.slice(0, 5);
}

export function summaryLine(summary: ScoreSummary): string {
	return `H ${summary.hits}   G ${summary.golds}   X ${summary.xCount}`;
}

export function endRow(end: { endNumber: number; arrows: { scoreCode: string; scoreValue: number }[] }): {
	endNumber: number;
	codes: string[];
	total: number;
} {
	return {
		endNumber: end.endNumber,
		codes: end.arrows.map((a) => a.scoreCode),
		total: end.arrows.reduce((sum, a) => sum + a.scoreValue, 0)
	};
}
