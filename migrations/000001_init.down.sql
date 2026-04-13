DROP TRIGGER IF EXISTS trigger_users_updated_at ON splix.users;
DROP TRIGGER IF EXISTS trigger_groups_updated_at ON splix.groups;
DROP TRIGGER IF EXISTS trigger_expenses_updated_at ON splix.expenses;

DROP FUNCTION IF EXISTS splix.update_updated_at_column();

DROP TABLE IF EXISTS splix.expense_splits;
DROP TABLE IF EXISTS splix.group_members;
DROP TABLE IF EXISTS splix.expenses;
DROP TABLE IF EXISTS splix.groups;
DROP TABLE IF EXISTS splix.users;

DROP SCHEMA IF EXISTS splix CASCADE;
