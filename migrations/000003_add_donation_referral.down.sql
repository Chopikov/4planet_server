-- Remove referral_user_id field from donations table

ALTER TABLE donations DROP CONSTRAINT IF EXISTS fk_donations_referral_user;
DROP INDEX IF EXISTS idx_donations_referral_user;
ALTER TABLE donations DROP COLUMN IF EXISTS referral_user_id;
