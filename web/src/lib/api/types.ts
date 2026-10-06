export type ScoreCode = 'X' | '10' | '9' | '8' | '7' | '6' | '5' | '4' | '3' | '2' | '1' | 'M';

export type SessionStatus = 'in_progress' | 'completed';

export type Archer = {
	id: string;
	displayName: string;
};

export type ScoreSummary = {
	total: number;
	hits: number;
	golds: number;
	xCount: number;
	arrowCount: number;
};

export type Arrow = {
	id: string;
	arrowNumber: number;
	scoreCode: ScoreCode;
	scoreValue: number;
	recordedAt: string;
};

export type End = {
	endNumber: number;
	arrows: Arrow[];
};

export type SessionListItem = {
	id: string;
	archerId: string;
	status: SessionStatus;
	startedAt: string;
	completedAt: string | null;
	summary: ScoreSummary;
};

export type Session = SessionListItem & {
	arrowsPerEnd: 6;
	ends: End[];
};

export type PersonalBest = {
	archerId: string;
	sessionId: string;
	total: number;
	arrowCount: number;
	achievedAt: string;
};

export type RecordArrowResponse = {
	arrow: Arrow;
	endNumber: number;
	summary: ScoreSummary;
};

export type ApiErrorBody = {
	code: string;
	message: string;
};

export const TONY_ID = '8f0c1a2e-0b3d-4a7c-9e11-01a1a1a1a1a1';
export const BECKY_ID = '8f0c1a2e-0b3d-4a7c-9e11-02b2b2b2b2b2';
