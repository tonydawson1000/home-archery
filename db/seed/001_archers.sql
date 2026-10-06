-- Stable ids so the PWA picker and localStorage can use known UUIDs in development.
INSERT INTO archers (id, display_name) VALUES
    ('8f0c1a2e-0b3d-4a7c-9e11-01a1a1a1a1a1', 'Tony'),
    ('8f0c1a2e-0b3d-4a7c-9e11-02b2b2b2b2b2', 'Becky')
ON CONFLICT (display_name) DO NOTHING;
