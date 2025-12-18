// TerminalTestHub Database Schema

-- Initialize exxtensions for UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table terminals (virtual POS devices)
CREATE TABLE IF NOT EXISTS terminals (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  serial_number TEXT UNIQUE NOT NULL,
  status TEXT NOT NULL DEFAULT 'idle',
  last_seen TIMESTAMPZ DEFAULT NOW(),
  metadata JSONB DEFAULT '{}'::jsonb,
  cerated_at TIMESTAMPZ DEFAULT NOW(),
)

-- Tabel of tasks (tests for terminals)
CREATE TABLE IF NOT EXISTS jobs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  terminal_id UUID NOT NULL REFERENCES terminals(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  paload JSONB DEFAULT '{}'::jsonb,
  status TEXT NOT NULL DEFAULT 'pending',
  result JSONB,
  cerated_at TIMESTAMPZ DEFAULT NOW(),
  updated_at TIMESTAMPZ DEFAULT NOW(),
)


-- Indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_terminals_status ON terminals(status);
CREATE INDEX IF NOT EXISTS idx_terminals_serilal ON terminals(serial_number);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_terminal_id ON jobs(terminal_id);
CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(cerated_at DESC);

-- Comments
COMMENT ON TABLE terminals IS 'Virtual POS-terminals for testing purposes';
COMMENT ON TABLE jobs IS 'Tasks to run tests on terminals (e.g., pending, in_progress, completed, failed)';

COMMENT ON COLUMN terminals.status IS "Status: idle(pending), running(running test), done(completed test), failed(error during test)";
COMMENT ON COLUMN jobs.staus IS "Statuses: pending(in queue) running(test in progress), done(test completed), failed(error during test)";




