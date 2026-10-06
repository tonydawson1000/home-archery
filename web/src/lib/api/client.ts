import type { ApiErrorBody, Archer, PersonalBest, RecordArrowResponse, ScoreCode, Session } from './types';

export type FetchFn = typeof fetch;

export class ApiError extends Error {
	constructor(
		public status: number,
		public code: string,
		message: string
	) {
		super(message);
		this.name = 'ApiError';
	}
}

const defaultBase = '/api/v1';

async function parse<T>(res: Response): Promise<T> {
	const body = (await res.json()) as T | ApiErrorBody;
	if (!res.ok) {
		const err = body as ApiErrorBody;
		throw new ApiError(res.status, err.code, err.message);
	}
	return body as T;
}

export function createClient(fetchImpl: FetchFn = fetch, base = defaultBase) {
	return {
		async listArchers(): Promise<Archer[]> {
			const res = await fetchImpl(`${base}/archers`);
			return parse<Archer[]>(res);
		},
		async startSession(archerId: string): Promise<Session> {
			const res = await fetchImpl(`${base}/archers/${archerId}/sessions`, { method: 'POST' });
			return parse<Session>(res);
		},
		async listSessions(archerId: string): Promise<Session[]> {
			const res = await fetchImpl(`${base}/archers/${archerId}/sessions`);
			return parse<Session[]>(res);
		},
		async getPersonalBests(archerId: string): Promise<PersonalBest[]> {
			const res = await fetchImpl(`${base}/archers/${archerId}/personal-best`);
			return parse<PersonalBest[]>(res);
		},
		async getSession(sessionId: string): Promise<Session> {
			const res = await fetchImpl(`${base}/sessions/${sessionId}`);
			return parse<Session>(res);
		},
		async recordArrow(sessionId: string, scoreCode: ScoreCode): Promise<RecordArrowResponse> {
			const res = await fetchImpl(`${base}/sessions/${sessionId}/arrows`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ scoreCode })
			});
			return parse<RecordArrowResponse>(res);
		},
		async completeSession(sessionId: string): Promise<Session> {
			const res = await fetchImpl(`${base}/sessions/${sessionId}/complete`, { method: 'POST' });
			return parse<Session>(res);
		}
	};
}

export const api = createClient();
